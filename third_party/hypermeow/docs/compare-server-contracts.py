"""Compare pinned public Server Go and Node contracts against the Zuck baseline."""

import argparse
import json
import subprocess
from collections import Counter
from pathlib import Path

HTTP_METHODS = {"get", "head", "options", "post", "put", "patch", "delete", "trace"}


def endpoints(repo, revision, path):
    raw = subprocess.check_output(["git", "-C", str(repo), "show", f"{revision}:{path}"])
    return {(method, route) for route, item in json.loads(raw)["paths"].items()
            for method in item if method in HTTP_METHODS}


def resolve(value, definitions):
    # The compact matrix may store a scalar, an operation or its status by reference.
    while isinstance(value, dict) and set(value) == {"$ref"}:
        value = definitions[value["$ref"].removeprefix("#/$defs/")]
    return value


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("go_checkout", type=Path)
    parser.add_argument("node_checkout", type=Path)
    args = parser.parse_args()
    docs = Path(__file__).resolve().parent
    matrix = json.loads((docs / "parity-contract-matrix.compact.json").read_bytes())
    progress = json.loads((docs / "parity-progress.json").read_bytes())
    node_revision = matrix["baseline"]["node"]["sha"]
    go_revision = progress["serverRevision"]
    go = endpoints(args.go_checkout, go_revision, "docs/openapi.json")
    node = endpoints(args.node_checkout, node_revision, "src/swagger.json")
    print(f"Go {go_revision}: {len(go)} operations")
    print(f"Node {node_revision}: {len(node)} operations")
    print(f"Same method and path: {len(go & node)}. Go only: {len(go - node)}. Node only: {len(node - go)}.")
    mapped = {ref for entry in progress["entries"] for ref in entry["refs"]}
    by_journey = Counter()
    mapped_by_journey = Counter()
    node_candidates = []
    for op in matrix["operations"]:
        journey = op["journey"]
        by_journey[journey] += 1
        if op["id"] in mapped:
            mapped_by_journey[journey] += 1
        baseline_go = resolve(op["go"], matrix["$defs"])
        baseline_node = resolve(op["node"], matrix["$defs"])
        if (resolve(baseline_go["status"], matrix["$defs"]) == "missing"
                and resolve(baseline_node["status"], matrix["$defs"]) == "core"):
            node_candidates.append(op)
    print("Journey | baseline operations | Zuck IDs with native contract evidence")
    for journey in sorted(by_journey):
        print(f"{journey} | {by_journey[journey]} | {mapped_by_journey[journey]}")
    print("Baseline Node core / Go missing candidates")
    for op in node_candidates:
        status = "mapped native contract" if op["id"] in mapped else "requires comparison"
        print(f"{op['id']} {op['priority']} {op['method']} {op['path']}: {status}")
    print("Matching routes and mapped contracts do not establish functional equivalence.")


if __name__ == "__main__":
    main()
