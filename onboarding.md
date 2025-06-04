# Onboarding Guide

Welcome to **SSH Tunnel Manager (sshtm)**! This guide provides a quick overview of the repository structure and pointers on where to start.

## Repository Overview
The project implements a CLI client and background daemon for managing SSH tunnel configurations. You can save, list, edit, delete, and start tunnels.

```
.
├── client/        # CLI application code
│   ├── cmd/       # Cobra commands (add, list, tunnel, etc.)
│   └── lib/       # Helper library (forms, gRPC client)
├── daemon/        # Background service
│   ├── tasks/     # Logic for configuration and tunnel management
│   ├── main.go    # gRPC server startup
│   └── rpcserver.go
├── pkg/           # Reusable libraries
│   ├── configmanager/  # Read/write configuration files
│   └── tunnelmanager/  # Maintain SSH connections
├── rpc/           # Protocol buffer definitions & generated code
├── utils/         # Small helper functions
└── scripts/       # Install/uninstall helpers
```

### Key Files and Features
- `config/config.go` defines defaults like the gRPC port (`50051`) and configuration directory.
- The daemon starts a gRPC server (`daemon/main.go`) and registers handlers implemented in the `tasks` package.
- `pkg/tunnelmanager` orchestrates SSH connections and tracks active tunnels.
- CLI commands under `client/cmd/` communicate with the daemon via the generated gRPC client.
- Interactive forms in `client/lib` use `tview` for adding or editing tunnels.

### Usage and Development
- Run `./install.sh` to build binaries and install a user-level systemd service for the daemon.
- During development, use `make gen_proto` followed by `air` for live reloading.
- `scripts/uninstall.sh` stops the service and removes installed files.

### Dependencies
The project relies on:
- [Cobra](https://github.com/spf13/cobra) for CLI command handling.
- [tview](https://github.com/rivo/tview) and [tcell](https://github.com/gdamore/tcell) for terminal UIs.
- gRPC/protobuf for client–daemon communication.
- `golang.org/x/crypto/ssh` for SSH tunneling.

## Learning Pointers
To dive deeper:
1. Review `rpc/daemon.proto` and its generated files to understand the gRPC API.
2. Explore `pkg/configmanager` for how configurations are stored and loaded.
3. Study `pkg/tunnelmanager` to see how SSH connections are started and managed.
4. Look at the Cobra commands in `client/cmd/` to learn how the CLI interfaces with the daemon.
5. Check out the systemd setup in `install.sh` and `scripts/uninstall.sh`.

This should give you a solid starting point for working with **sshtm**. Happy hacking!
