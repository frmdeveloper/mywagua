# Implementation and validation checkpoint: 2026-11-10

This checkpoint records completed platform work and the authorized real-phone
checks described below. It does not claim complete ZuckZapGo parity. The [303-operation baseline](parity-contract-matrix.md)
and [phased plan](../PLAN.md) retain the outstanding work explicitly.

| Component | Validated source | Result |
| --- | --- | --- |
| HyperMeow | d326350a86e98c405c9f6adf14cb6bd03ea8e66a | Tulir updates, migration correction, inherited tests and registered-recipient LID fallback incorporated |
| Server Go | e01698bf3e8b0b9d9dc749e7ffbb49e65deab22f | HyperMeow adapter, signed webhooks, send ledger, selective event recovery and recipient-resolution fix incorporated |
| Cloud Manager | 3ade44f97868c008c30a1e0894f9636cfa851c51 | Go-only allocation, authenticated workspace ownership, quota/concurrency checks |
| Cloud worker | 95aa5d279778e84707b5bfde886cdd455bee48fe | Go-only reserves and isolated per-session signed receiver provisioning |
| Cloud | 35ac0c4080dd375a3efa1625567e55d93c856958 | No runtime selection required; PT/EN docs, idempotency forwarding/recovery and rejection on old runtimes before sending |

Validated Go runtime image:
`sha256:b0d96f63e775038ff4c71cfdf4630763712fe8e978ff555dedb685cae8ed9b5b`.
The build used a clean checkout of the source above; its git source archive SHA256
was `17c5ed3a0018f94f2b3527b547ccc751b90092060067d62e402d43cb454e735e`.
All three Cloud reserves use this image. Existing customer allocations were
preserved; the previously connected session remained connected.

## Automated validation

- [Recipient SDK correction CI](https://github.com/wppconnect-team/hypermeow/actions/runs/42886170997): regression reproduced through USync serialization/parsing, full tests, race, supported Go versions, WASM, MEX and pre-commit.
- [Updated server CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42887160468): full Go/race/vet, broker/S3 integrations and AMD64/ARM64 image smokes passed before rollout.
- [HyperMeow integration CI](https://github.com/wppconnect-team/hypermeow/actions/runs/42866085986): supported Go versions, race, WASM, MEX and pre-commit.
- [Signed public receiver CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42868894279): Go/race, broker/S3 integrations and AMD64/ARM64 image smokes.
- [Send ledger CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42873726008): restart, concurrency, storage failures, cancellation, tenant isolation and simulated crash regressions, plus the full suite.
- [Delivery operations CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42875435897): old schema migration, atomic audit, selective replay, leases, target/session filters and receiver recovery, plus the full suite.
- [Cloud CI](https://github.com/joaosouz4dev/wppconnect-cloud/actions/runs/42876439548): 48 local unit tests and production build also passed; checks include preserving POST bytes across the GET ledger probe and rejecting old runtimes without a POST.

## Runtime checks

Two temporary HTTPS receivers with different keys verified HMAC over raw bytes,
stable event IDs across retries, no cross-session events and restart persistence.
A Manager allocation completed in 5.91 seconds, returned a valid QR PNG, delivered
its signed event, enforced cross-workspace 404 and quota 409, and was deleted.

The final Go image was separately exercised against a receiver held at HTTP 503.
The delivery became failed, selective replay recovered it after receiver repair,
audit and capacity metrics were readable, and another session could not replay
it. An invalid send on an unpaired fixture was recorded once, recovered with the
same key, rejected after a payload change, and remained recoverable after restart.
The production Manager also preserved the key and exposed the stored result.

Temporary receivers, containers, network and fixture secrets were removed.
The final runtime audit found three healthy Go reserves, no OOM, and no
panic/fatal/persistence-failure matches in the inspected ten-minute log window.

## Authorized real-phone checks

The dedicated test account paired successfully through a temporary page that
refreshed its QR automatically. The page, route and private helper configuration
were removed after connection.

The first text attempt returned HTTP 502. A probe that sent no messages showed
that profile USync omitted the LID while contact discovery recognized the same
registered recipient. [SDK PR #4](https://github.com/wppconnect-team/hypermeow/pull/4)
added a strictly matched contact lookup fallback, and
[server PR #20](https://github.com/wppconnect-team/wppconnect-server-go/pull/20)
pinned the correction. The failed command remains in the ledger and was not
automatically resent with a different key.

After the update, a separately identified text test and one small text file
returned HTTP 201, received delivery and read receipts, and elicited the exact
test reply observed by the server. Repeating each command with its original key
returned the identical response; changing the payload returned HTTP 409. Another
workspace could not read the paired session's send result.

The paired runtime reconnected after restart. Both successful results and the
original failed result survived. The file downloaded after restart matched the
original bytes. An offline restoration of the consistent pre-update backup
passed integrity checks for all five SQLite databases; no second WhatsApp client
was connected to the restored copy.

The first container recreation attempt rejected a network alias on Docker's
default bridge. The helper script was corrected, the previous runtime was
restored and reconnected, and the subsequent update preserved the session,
environment and published port. Customer allocations were not changed.

See the [sanitized acceptance receipt](phone-acceptance-2026-11-10.json). It contains
no phone number, session identifier, message content, pairing URL or credential.

## Limits of this checkpoint

The earlier isolated runtime checks used unpaired fixtures. The real-phone
checks above cover one sending account and one authorized recipient, not the
complete contract matrix. Logout, two paired tenants, phone-code pairing,
Android/iOS/Web interactive behavior, calls and the remaining application
integrations still require acceptance work. The test account remains paired for
those remaining journeys.
Prometheus alert rules are examples; no new operational alerting service or
notification destination was installed.

The supplied Docker image contains the compiled private application, UI and
contract documentation, not its complete buildable private source. The public
repository preserves the available HyperMeow source and history; it is not a
recovered copy of every private ZuckZapGo component.
