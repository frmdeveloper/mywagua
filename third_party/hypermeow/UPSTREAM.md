# Upstream synchronization policy

- Primary protocol upstream: https://github.com/tulir/whatsmeow, resolve default branch via GitHub.
- Public feature upstream: https://github.com/polymorfa/hypermeow, resolve its default branch separately.
- Destination: https://github.com/wppconnect-team/hypermeow, main.
- Daily schedule: 09:30 America/Sao_Paulo through the requested Codex task. Review the public Web transport, the active HyperMeow implementation and the version pinned by server-go. Record versioned evidence and run socket regression tests when a behavior changes. [Transport study and reproduction](docs/web-socket.md).

Use a fresh clone or clean worktree. Fetch both upstreams and origin; record exact
SHAs, merge base, ahead/behind and patch-equivalent commits. Preserve local work.
Create a candidate branch for integration. Merge normally, never reset main or
force-push. Resolve conflicts by preserving intended behavior on both sides;
renaming upstream imports must remain consistent with the inherited module path.
Generated protobuf/MEX files require their corresponding sources/generator checks.
Do not lower SQL schema versions or discard fork migrations, namespaces, locks,
performance changes, public API, attribution or protocol corrections.

Run the CI-equivalent tests, vet, build, MEX generation check, JS/WASM build and
pre-commit checks. Test SQL migrations when the store changes, and dependent
server-go builds when public APIs change. Offline passing tests do not waive
account-dependent validation of protocol changes. Never reuse production data
for tests or send messages to unrelated recipients.

Push a candidate and open or update a single upstream-sync PR when appropriate.
Only merge/push main after all applicable checks pass and required branch rules
allow it. Record CI failures or unavailable runners; absence of CI is not success.
Re-fetch main immediately before merge, incorporate intervening changes and rerun
affected checks. Verify the destination SHA after update. No npm publication or
semver tag is implied by this task; Go consumers pin the verified commit.

On conflicts, keep the candidate and a sanitized report. Do not repeatedly open
identical PRs. Reuse the pending candidate where safe. Report source/destination
SHAs, relevant changes, validation, PR/merge link and remaining blockers. If neither
upstream changed and no pending integration needs work, produce no empty commit.

Initial comparison (2026-11-10): Polymorfa f432198 vs tulir c88388f has 37 tulir-only
and 293 fork-only commits; merge-tree reports conflicts in 15 files. This is not
an already synchronized base. Resolve and validate those conflicts before claiming
that the fork includes the current tulir main.


## Integration candidate 2026-11-10

Source: tulir/whatsmeow c88388f983b4bc43dfd0bcc2b46969fcde82d61f.
The 37 upstream commits were merged into the public HyperMeow base. Resolution:

- Keep bounded caches, shadow/remote signal support, synchronous name updates and
  deferred app-state error dispatch; add upstream backoff and unified sessions.
- Keep HTTP rejection precedence when racing sockets, add handshake/closed socket
  retry classifications. Keep HyperMeow's lock-protected SetOnDisconnect API:
  it already serializes setter/close access and avoids breaking existing callers.
  Direct writes to OnDisconnect while in use remain unsupported.
- Retain context-aware frame completion and the fork's boolean return contract.
- Combine both XML credential-redaction sets without mutating protocol payloads.
- Preserve ResolveLID cache/USync behavior, wrapping missing-LID failures with the
  upstream ErrRecipientNotOnWhatsApp sentinel without a second network lookup.
- Keep fork module/import namespace, libsignal fork and schema version 19.
  Include upstream group-admin metadata in the bounded cache.
- Regenerate 13 changed protobuf definitions from merged schemas, including the
  new Hatch schema in the inherited namespace. Generator: protoc 33.6 (Windows
  archive SHA256 96a83b8ee942f1046f079c43657cc09fd41f37a4e349ddb0e8ba5064c48bb5ff),
  protoc-gen-go v1.36.12. MEX generation remains unchanged.

Local Go 1.26.10 test, vet, build, MEX check and JS/WASM build passed. New regression
coverage checks backoff before store access, release of its lock, handshake retry
classification, HTTP rejection precedence and combined credential redaction.
CI, race and dependent server validation are recorded in the integration PR.


Dependent server validation exposed overlapping schema version 16: tulir has
wasa_root_secret_id while HyperMeow uses that version for a sender-key index.
The v18->v19 upgrade now checks column presence and ensures the missing index
without overwriting root secrets. The server regression uses the exact previous
v16 schema and verifies credential/root-secret preservation and repeat upgrade.
Do not downgrade a live database in place; rollback restores a tested backup.
