# Onboarding Guide

SSH Tunnel Manager is a Rust CLI and background daemon for reusable SSH local
port forwards. Both binaries share the `sshtm` library crate.

## Repository layout

```text
Cargo.toml          Package, binaries, and dependencies
build.rs            Protobuf binding generation
src/bin/sshtm.rs    Clap CLI and interactive terminal workflow
src/bin/sshtmd.rs   Tokio/tonic daemon entry point and graceful shutdown
src/daemon.rs       gRPC handlers and supervised OpenSSH tunnel registry
src/store.rs        Go-compatible JSON configuration persistence
rpc/daemon.proto    Stable client/daemon API
packaging/          Arch, Debian, systemd, and shell integration
```

The daemon listens on `127.0.0.1:50051`. It validates and persists profiles,
starts one supervised OpenSSH process per active local forward, and writes
`active_tunnels.json` so tunnels can be restored after restart. The CLI uses the
generated tonic client for every operation.

## Development

Install Rust, Cargo, `protoc`, and OpenSSH, then run:

```sh
cargo fmt --check
cargo clippy --workspace --all-targets -- -D warnings
cargo test
cargo run --bin sshtmd
cargo run --bin sshtm -- list
```

`cargo build` automatically regenerates protobuf bindings. Add domain and storage
tests close to their modules; exercise RPC behavior through the daemon service.
Preserve the protobuf field numbers, profile JSON's PascalCase field names, and
the active-tunnel JSON schema because released Go versions used those formats.
