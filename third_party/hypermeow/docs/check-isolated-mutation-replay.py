"""Exercise all five receipts across a process restart with a fresh, disconnected Go server."""
import argparse
from contextlib import contextmanager
import json
import os
from pathlib import Path
import secrets
import socket
import subprocess
import tempfile
import time
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen

COMMANDS = {
    "delete-newsletter": {"id": "123@newsletter"},
    "del-collection": {"collectionIds": ["fixture"]},
    "add-product": {"product": {"name": "fixture", "image_urls": ["https://mmg.whatsapp.net/fixture"]}},
    "edit-product": {"productId": "p", "product": {"name": "fixture", "is_hidden": True,
                                                  "image_urls": ["https://mmg.whatsapp.net/fixture"]}},
    "del-products": {"productIds": ["p"]},
}


def call(url, method, expected, body=None, headers=None):
    data = json.dumps(body).encode() if body is not None else None
    request = Request(url, data=data, method=method, headers=headers or {})
    try:
        response = urlopen(request, timeout=5)
    except HTTPError as err:
        response = err
    with response:
        if response.status != expected:
            raise RuntimeError(f"{method} expected {expected}, got {response.status}")
        return json.load(response), response.headers


@contextmanager
def server(executable, directory, env, base):
    proc = subprocess.Popen([executable], cwd=directory, env=env,
                            stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    try:
        for _ in range(100):
            if proc.poll() is not None:
                raise RuntimeError("isolated server exited before readiness")
            try:
                with urlopen(base + "/healthz", timeout=1) as result:
                    if result.status == 200:
                        break
            except (URLError, OSError):
                time.sleep(0.1)
        else:
            raise RuntimeError("isolated server did not become ready")
        yield
    finally:
        proc.kill()
        proc.wait(timeout=15)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("server_checkout", type=Path)
    args = parser.parse_args()
    with tempfile.TemporaryDirectory() as directory:
        executable = Path(directory) / ("server.exe" if os.name == "nt" else "server")
        subprocess.run(["go", "build", "-o", executable, "./cmd/server"], cwd=args.server_checkout, check=True)
        with socket.socket() as sock:
            sock.bind(("127.0.0.1", 0))
            port = sock.getsockname()[1]
        base = f"http://127.0.0.1:{port}"
        secret = secrets.token_hex(16)
        env = {k: os.environ[k] for k in ("SystemRoot", "TEMP", "TMP", "PATH") if k in os.environ}
        env.update(PORT=str(port), SECRET_KEY=secret, DATA_DIR=str(Path(directory) / "data"),
                   AUTO_RESTORE_SESSIONS="false", ENABLE_MANAGER="false")
        session = base + "/api/fixture/"
        with server(executable, directory, env, base):
            auth, _ = call(session + secret + "/generate-token", "POST", 201)
            headers = {"Authorization": "Bearer " + auth["token"], "Content-Type": "application/json"}
            for command, body in COMMANDS.items():
                call(session + command, "POST", 400, body, headers)
                result, _ = call(session + command, "POST", 404, body, {**headers, "Idempotency-Key": command})
                if result.get("status") != "disconnected":
                    raise RuntimeError(f"{command}: expected disconnected session")
        with server(executable, directory, env, base):
            for command, body in COMMANDS.items():
                keyed = {**headers, "Idempotency-Key": command}
                lookup = "newsletter-admin-result" if command == "delete-newsletter" else "catalog-write-result"
                stored, _ = call(session + lookup, "GET", 200, headers=keyed)
                result, replay = call(session + command, "POST", 404, body, keyed)
                if (replay.get("Idempotency-Replayed") != "true" or stored.get("httpStatus") != 404
                        or stored.get("response") != result):
                    raise RuntimeError(f"{command}: original receipt was not recovered")
                call(session + command, "POST", 409, {}, keyed)
    print("Five disconnected commands retained their scoped receipts and rejected changed input after restart.")


if __name__ == "__main__":
    main()
