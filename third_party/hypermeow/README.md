# WPPConnect HyperMeow

Public WPPConnect-maintained fork of [Polymorfa HyperMeow](https://github.com/polymorfa/hypermeow),
which derives from [tulir/whatsmeow](https://github.com/tulir/whatsmeow).
This repository is a Go WhatsApp protocol SDK, not the ZuckZapGo application server.

## Provenance and licensing

Initial public source: Polymorfa `f43219816d8600ec963371881a1467d22fa7b3e6`.
The complete reachable history is preserved. Read [LICENSING.md](LICENSING.md),
[LICENSE](LICENSE) and [LICENSE-MIT](LICENSE-MIT) for file-specific licensing.
Existing copyright notices and third-party attributions remain intact.

The inspected ZuckZapGo v2.4.2 executable uses a local `hypermeow-private` replacement.
Its container contains compiled application code and generated documentation, not
that private SDK's compilable source. This public fork does not claim to contain
those private changes or all 303 application API operations.

## Using this fork

The inherited Go module path remains `github.com/polymorfa/hypermeow` so that
internal imports, generated protobuf descriptors and upstream merges remain coherent.
Pin a verified WPPConnect commit using a Go replace directive:

```sh
go mod edit -replace=github.com/polymorfa/hypermeow=github.com/wppconnect-team/hypermeow@<verified-commit>
go mod tidy
```

Applications must consistently import the HyperMeow module. Replacing
`go.mau.fi/whatsmeow` alone does not rewrite its imports or make the APIs identical.
Before migrating WPPConnect Server Go, validate API changes, device stores,
SQL migrations, protocol compatibility, reconnect and backup/restore using a copy
of an authorized test session. Do not point two clients at the same session data.

## Development and checks

See [upstream documentation](README.upstream.md) for SDK features and usage.
The inherited CI supports Go 1.26 and 1.27, checks generated MEX bindings and builds
JS/WASM. The go.mod toolchain preference can be overridden with a supported,
verified Go toolchain on the developer machine.

```sh
go test ./...
go vet ./...
go build ./...
go run ./internal/cmd/genmex -input mex/spec.json -output mex/bindings.go -check
GOOS=js GOARCH=wasm go build ./...
```

Offline tests do not establish successful pairing, WhatsApp message delivery,
private-extension parity or safe migration of a production database.

## Upstream maintenance

Daily maintenance compares against `tulir/whatsmeow` and also records the public
Polymorfa source. [UPSTREAM.md](UPSTREAM.md) defines the merge and validation policy.
No force-push or replacement of WPPConnect/HyperMeow-specific code is permitted.
A conflicting or failing update remains unmerged and is reported with evidence.

The phased Cloud and parity work is tracked in [PLAN.md](PLAN.md).

See the [303-operation contract matrix](docs/parity-contract-matrix.md) for the fixed Zuck/Go/Node baseline, candidate adapters, limitations and acceptance criteria. A matching endpoint is not proof of functional parity.
