"""Check the progress overlay against a local Server Go checkout. Does not send requests."""
import argparse
import json
import re
import subprocess
from pathlib import Path

def key_parameter(operation):
    return next((p for p in operation["parameters"]
                 if p.get("in") == "header" and p.get("name") == "Idempotency-Key"), None)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("server_checkout", type=Path)
    parser.add_argument("--run-fixtures", action="store_true",
                        help="execute all mapped HTTP test functions and verify each actually passed")
    args = parser.parse_args()
    docs = Path(__file__).resolve().parent
    baseline = json.loads((docs / "parity-contract-matrix.compact.json").read_text(encoding="utf-8"))
    overlay = json.loads((docs / "parity-progress.json").read_text(encoding="utf-8"))
    spec = json.loads((args.server_checkout / "docs/openapi.json").read_text(encoding="utf-8"))
    operations = {op["id"]: op for op in baseline["operations"]}
    if len(operations) != overlay["baselineCount"]: raise ValueError("baseline count changed")
    sources = [p.read_text(encoding="utf-8") for p in (args.server_checkout / "internal/httpapi").glob("*_test.go")]
    seen = set()
    fixtures = set()
    for item in overlay["entries"]:
        name = item["operation"]
        try:
            if name in seen: raise ValueError("duplicate entry")
            seen.add(name)
            journeys = {operations[z]["journey"] for z in item["refs"]}
            if not journeys or journeys != set(item["journeys"]): raise ValueError("baseline journey mismatch")
            op = spec["paths"]["/api/{session}/" + name][item["method"]]
            if not op.get("security") or op.get("x-implemented") is not True: raise ValueError("unprotected or unsupported")
            if str(item["success"]) not in op["responses"]: raise ValueError("success status missing")
            body = op["requestBody"]["content"]["application/json"]["schema"]
            if set(body.get("required", [])) != set(item["required"]): raise ValueError("required fields changed")
            key = key_parameter(op)
            if bool(key) != item["idempotency"]: raise ValueError("idempotency header changed")
            if item.get("requireKey"):
                if not key or not key.get("required"): raise ValueError("required idempotency header missing")
                lookup = spec["paths"]["/api/{session}/" + item["resultOperation"]]["get"]
                if (not lookup.get("security") or lookup.get("x-implemented") is not True
                        or not (key_parameter(lookup) or {}).get("required")
                        or "200" not in lookup["responses"]):
                    raise ValueError("protected result lookup missing")
            if not any(("func " + item["fixture"] + "(") in s for s in sources): raise ValueError("contract test missing")
            fixtures.add(item["fixture"])
            if item["liveEquivalence"] or item["evidence"] != "automated-contract" or not item["limit"]:
                raise ValueError("automated evidence must retain a parity limitation")
        except (KeyError, ValueError) as err:
            raise SystemExit(f"{name}: {err}") from err
    if args.run_fixtures:
        pattern = "^(" + "|".join(re.escape(name) for name in sorted(fixtures)) + ")$"
        run = subprocess.run(["go", "test", "-json", "-count=1", "-run", pattern, "./internal/httpapi"],
                             cwd=args.server_checkout, stdout=subprocess.PIPE, text=True, encoding="utf-8")
        events = [json.loads(line) for line in run.stdout.splitlines() if line.strip()]
        passed = {event.get("Test") for event in events if event.get("Action") == "pass"}
        if run.returncode or not fixtures.issubset(passed):
            for event in events:
                if event.get("Action") in ("output", "build-output"):
                    print(event.get("Output", ""), end="")
            raise SystemExit("Mapped fixtures without a pass: " + ", ".join(sorted(fixtures - passed)))
        print(f"{len(fixtures)} distinct mapped HTTP fixtures passed.")
    print(f"{len(seen)} native contracts mapped against {len(operations)} baseline operations; no live equivalence inferred.")
if __name__ == "__main__":
    main()
