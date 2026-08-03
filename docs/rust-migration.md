# Go-to-Rust Migration Plan

> **Status: implemented.** The repository now builds `sshtm` and `sshtmd` from a
> Rust package. The protobuf field numbers, binary names, installation paths,
> profile JSON, active-tunnel JSON, and service identifiers were retained. The
> implementation uses tonic/Tokio for the daemon, Clap for the CLI, Serde for
> persistence, and supervised OpenSSH processes for forwarding. The remainder of
> this document is retained as the compatibility and verification record that
> should guide future changes.

## Objective

Replace the Go implementation of `sshtm` and `sshtmd` with Rust while preserving
the command-line interface, the daemon's gRPC contract, saved configuration data,
active-tunnel restoration, installation paths, and Linux/macOS service behavior.
The migration should be incremental: a released Go client must be able to talk to
a Rust daemon, and a Rust client must be able to talk to a released Go daemon,
until the final Go removal release.

This is a delivery plan, not a commitment to a big-bang rewrite. Each phase has
an independently reversible release gate.

## Current-system baseline

Before implementation starts, record the behavior of the current application as
the compatibility baseline:

- Two executables are produced: the Cobra-based `sshtm` client and the `sshtmd`
  daemon.
- Client and daemon communicate on `localhost:50051` with the service defined in
  `rpc/daemon.proto`.
- Tunnel profiles and `active_tunnels.json` live under
  `~/.ssh-tunnel-manager` by default.
- The daemon owns configuration mutation, SSH connection lifecycle, port
  allocation, persistence, and restoration after restart.
- The CLI supports `list`, `add`, `edit`, `delete`, `tunnel`, `active`, `kill`,
  `completion`, `version`, and `help`, including the existing aliases and
  human-readable output.
- Linux systemd user services and macOS LaunchAgents invoke binaries installed in
  the existing user-local paths.

The first migration pull request should turn these statements into executable
black-box acceptance tests. Capture stdout, stderr, exit status, RPC status and
message fields, configuration JSON, signal handling, and tunnel restoration.
Normalize only nondeterministic values such as allocated ports and temporary
paths.

## Compatibility contracts

These contracts remain fixed throughout the dual-language period:

1. **Protocol:** Keep `rpc/daemon.proto` as the single source of truth. Do not
   renumber or reuse fields or enum values. Additive RPC changes must be optional,
   and generated Go and Rust stubs must be checked in or reproducibly generated.
2. **Storage:** Rust must read profiles and `active_tunnels.json` written by Go,
   and Go must continue reading files written by Rust. Preserve field names,
   integer semantics, permissions, and the zero/missing-field behavior. Writes
   should use a temporary file, flush it, and atomically rename it.
3. **CLI:** Preserve command names, aliases, positional arguments, completion,
   exit codes, and scripting-oriented output. Intentional UX changes require a
   separately documented compatibility decision.
4. **Networking:** Keep the client endpoint (`localhost:50051`) and port during
   migration. The Go daemon currently listens on the wildcard address, so the
   Rust binding must be an explicit pre-rollout decision: reproduce that behavior
   for strict compatibility or move both implementations to loopback in a
   separately reviewed hardening change. Changing transport or adding local
   authentication is likewise outside the rewrite's compatibility scope.
5. **Lifecycle:** SIGINT/SIGTERM must stop accepting work, persist active tunnels,
   close listeners and SSH connections, and exit within the service manager's
   timeout. Restart must restore the same set of tunnels or report each failure.
6. **Distribution:** Retain binary names, install locations, service identifiers,
   configuration locations, and uninstall behavior so upgrades do not require a
   manual migration.

## Proposed Rust architecture

Use a Cargo workspace so domain logic is shared without coupling it to either
executable:

```text
Cargo.toml
crates/
  sshtm-core/       domain types, validation, configuration repository
  sshtm-proto/      generated protobuf/gRPC types
  sshtm-tunnel/     SSH sessions, forwards, registry, persistence/restoration
  sshtm-daemon/     gRPC handlers, startup, shutdown, observability
  sshtm-cli/        commands, formatters, interactive forms, daemon client
tests/
  compatibility/    Go/Rust black-box and cross-version fixtures
```

Keep dependency choices behind project-owned interfaces. The initial technical
spike should evaluate `tokio` for async execution, `tonic`/`prost` for gRPC and
protobuf, `clap` for command parsing and completion, `serde`/`serde_json` for
storage, `tracing` for diagnostics, and a maintained Rust SSH implementation.
`ratatui` plus a terminal backend is a candidate for the interactive UI. Do not
select the SSH or TUI library until the spike proves keyboard behavior, private
key support, cancellation, reconnect behavior, Linux/macOS support, and license
compatibility.

Suggested boundaries:

- `ConfigRepository` provides list/get/add/update/delete operations and owns
  atomic persistence.
- `TunnelBackend` abstracts establishing an SSH session and opening a local
  forward so lifecycle tests can use a deterministic fake.
- `TunnelRegistry` owns active tunnel metadata and cancellation handles; it must
  not expose locks across `.await` points.
- `TunnelService` implements application use cases and is independent of gRPC.
- Thin gRPC adapters translate protobuf requests and domain errors into the
  existing response status/message fields.
- CLI formatters accept structured response types and write to an injected
  writer, allowing golden tests without a terminal.

## Delivery phases

### Phase 0 — Freeze and measure behavior

**Work**

- Add golden CLI fixtures for every command, alias, validation failure, and empty
  state.
- Add protocol descriptor checks and RPC integration tests covering all methods.
- Add storage fixtures for valid, missing, malformed, and legacy profile data.
- Add an in-process SSH test server or hermetic test container for forwarding,
  cancellation, port conflicts, authentication failures, and reconnects.
- Record baseline startup time, idle memory, tunnel setup latency, and binary
  sizes. These are regression signals rather than reasons to relax correctness.

**Exit gate:** Current Go tests and the new black-box suite pass reliably on
Linux and macOS CI. No Rust production path exists yet.

### Phase 1 — Scaffold and prove risky dependencies

**Work**

- Create the Cargo workspace, formatting/lint policy, dependency audit, locked
  dependency file, and CI jobs for supported targets.
- Generate Rust protobuf code from the unchanged schema and add a descriptor
  comparison check.
- Prototype SSH key loading, local forwarding, cancellation, connection loss,
  and listener cleanup against the hermetic SSH server.
- Prototype the interactive form's focus, escape, Ctrl-C, validation, and
  non-interactive terminal behavior.
- Write an architecture decision record for the selected SSH and TUI crates,
  including maintenance, platform, security, and licensing evidence.

**Exit gate:** Both spikes meet the compatibility suite on Linux and macOS. If
they do not, replace the candidate dependency without affecting later crates.

### Phase 2 — Port domain and storage logic

**Work**

- Implement domain types and validation in `sshtm-core` with explicit, typed
  errors.
- Implement `ConfigRepository`, home-directory expansion, duplicate handling,
  filtering, and atomic writes.
- Implement serializers for saved configurations and active-tunnel state using
  checked-in fixtures produced by Go.
- Run bidirectional tests: Go writes/Rust reads and Rust writes/Go reads.

**Exit gate:** Fixture parity is exact where formatting is contractual and
semantically equivalent otherwise. The Go executables remain the release
artifacts.

### Phase 3 — Port the daemon behind the stable RPC API

**Work**

- Implement `TunnelBackend`, registry, service layer, all gRPC handlers, and
  graceful shutdown.
- Preserve structured response fields and legacy `result` text while both are
  supported by the protobuf schema.
- Add concurrency tests for simultaneous start/kill/list operations, duplicate
  local ports, dropped SSH sessions, daemon shutdown, and restoration.
- Run the full matrix: Go client/Go daemon, Go client/Rust daemon, Rust test
  client/Go daemon, and Rust test client/Rust daemon.
- Package the Rust daemon as an opt-in artifact and canary it before making it the
  default. Keep the Go daemon artifact available for one stable release as the
  rollback path.

**Exit gate:** The Go CLI passes every acceptance test against the Rust daemon;
there are no known persistence, file-descriptor, task, or connection leaks; and
service-manager restart tests pass on both platforms.

### Phase 4 — Port the non-interactive CLI

**Work**

- Implement global options, every command and alias, argument validation,
  completion, version injection, and daemon connection errors.
- Port formatters with golden tests for TTY/non-TTY and color/no-color output.
- Execute the four-way client/daemon matrix and compare results with normalized
  Go baselines.
- Ship the Rust CLI as an opt-in artifact, then promote it while retaining the Go
  CLI for one stable rollback release.

**Exit gate:** Existing shell usage and documented commands work unchanged
against both daemons, and the compatibility suite finds no unexplained output or
exit-code differences.

### Phase 5 — Port the interactive UI

**Work**

- Reproduce add/edit/delete flows, focus order, selection state, filtering,
  active-state decoration, logging, resize behavior, escape, and Ctrl-C handling.
- Add state-machine tests below the terminal layer and pseudo-terminal smoke
  tests for critical keyboard flows.
- Perform manual accessibility and terminal checks on common Linux terminals,
  macOS Terminal, and a small/narrow viewport.

**Exit gate:** The keyboard-flow checklist and pseudo-terminal tests pass, with
any deliberate visual change approved and documented.

### Phase 6 — Switch builds, packaging, and releases

**Work**

- Make Cargo the default build while retaining temporary targets for the Go
  rollback binaries.
- Update `install.sh`, Arch, Debian, systemd, launchd, checksum, and release jobs
  without changing installed paths or service identities.
- Verify clean install, in-place upgrade from the last Go release, downgrade,
  uninstall, and reinstall on Linux amd64/arm64 and macOS amd64/arm64 as supported
  by release policy.
- Document Rust toolchain prerequisites, protobuf generation, local development,
  and release/version injection.

**Exit gate:** Signed release artifacts pass installation and upgrade smoke tests,
and rollback restores the Go binaries without losing configurations.

### Phase 7 — Remove Go after the rollback window

**Work**

- Observe at least one stable Rust release and resolve migration regressions.
- Delete Go sources, generated Go protobuf files, Go-specific tooling, and dual
  implementation tests only after the support window closes.
- Retain immutable compatibility fixtures and end-to-end tests so future Rust
  refactors cannot silently break old installations.
- Update contributor, packaging, security, and release documentation.

**Exit gate:** Rust is the only shipped implementation, no packaging path refers
to Go, upgrade from the final Go release remains tested, and maintainers approve
the removal.

## Verification matrix

Every phase should run the applicable subset of this matrix:

| Area | Required checks |
| --- | --- |
| Rust quality | `cargo fmt --check`, `cargo clippy --workspace --all-targets --all-features -- -D warnings`, `cargo test --workspace --all-features` |
| Existing implementation | `go test ./...` until Go removal |
| Protocol | descriptor compatibility plus all four Go/Rust client-daemon combinations |
| Storage | checked-in cross-language fixtures, malformed input, permissions, atomic-write interruption |
| SSH | auth success/failure, forwarding, concurrent connections, port collision, remote drop, cancellation, reconnect |
| Lifecycle | SIGINT, SIGTERM, service restart, active-state persist/restore, repeated start/stop leak test |
| CLI | golden stdout/stderr/exit codes, aliases, completion, TTY/color variations |
| Distribution | clean install, Go-to-Rust upgrade, rollback, uninstall, and package validation per target |
| Security | dependency/license audit, approved listener-address assertion, secret-redaction and file-permission tests |

CI should fail on an unreviewed protobuf descriptor change, generated-code drift,
or compatibility fixture change. Performance comparisons should report trends;
maintainers should define numerical release thresholds after Phase 0 measurements
rather than inventing them in advance.

## Major risks and mitigations

- **SSH semantic differences:** Authentication algorithms, key formats, keepalive,
  reconnect, and cancellation may differ. Mitigate with the Phase 1 spike and a
  hermetic SSH compatibility suite before choosing a library.
- **Async lifecycle leaks:** Detached tasks can retain listeners or connections.
  Use structured ownership and cancellation, bounded shutdown, repeated
  start/stop tests, and file-descriptor monitoring.
- **Data loss during mixed-version use:** Preserve schemas, use atomic writes,
  maintain golden fixtures, and test both read/write directions before the Rust
  daemon writes user data.
- **CLI/TUI drift:** Separate formatters and state transitions from terminal I/O;
  use golden, state-machine, and pseudo-terminal tests.
- **Packaging regression:** Binary names may stay unchanged while build inputs
  change. Test actual packages and in-place upgrades in clean Linux/macOS runners,
  not only raw binaries.
- **Unbounded scope:** Transport hardening, configuration redesign, and UX cleanup
  are worthwhile but should not share the migration critical path unless required
  to fix a demonstrated security issue.

## Pull-request sequence and ownership

Prefer small pull requests that map to reversible milestones:

1. Compatibility test harness and baseline fixtures.
2. Cargo workspace, CI, protobuf generation, and dependency ADRs.
3. Core domain/config repository and cross-language storage tests.
4. SSH backend and tunnel registry.
5. Daemon service and gRPC adapters.
6. Canary daemon packaging and cross-version matrix.
7. CLI commands and output formatters.
8. Interactive UI.
9. Default Rust packaging, install/upgrade/rollback validation, and documentation.
10. Go removal after the stable-release observation window.

Assign one maintainer as compatibility owner (protocol, storage, and CLI contract)
and another as release owner (CI, packages, upgrade, and rollback). Security review
is required for the SSH dependency decision and before the Rust daemon becomes
the default.

## Definition of done

The migration is complete only when:

- all shipped commands and daemon RPCs meet the recorded compatibility baseline;
- configurations and active-tunnel data survive Go-to-Rust upgrade and tested
  rollback without manual conversion;
- forwarding, cancellation, shutdown, persistence, and restoration are covered by
  hermetic automated tests;
- supported Linux/macOS artifacts install, upgrade, run, and uninstall through
  their real service managers;
- dependency, license, secret-handling, and listener-binding reviews pass;
- operator and contributor documentation describes only the supported Rust path;
  and
- Go code is removed only after the announced rollback window has elapsed.
