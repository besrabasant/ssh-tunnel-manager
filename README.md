# SSH Tunnel Manager (`sshtm`)

SSH Tunnel Manager (**sshtm**) is a CLI and background daemon for saving, starting, listing, and stopping SSH local port-forwarding tunnels. It is useful when you frequently connect to the same private services through SSH and want reusable tunnel profiles instead of typing long `ssh -L` commands.

Current version: **v1.1.5**

## What it does

- Manages reusable SSH tunnel configurations with names and descriptions.
- Starts tunnels from saved configurations, equivalent to:
  ```sh
  ssh -L [LOCAL_IP:]LOCAL_PORT:DESTINATION:DESTINATION_PORT [USER@]SSH_SERVER
  ```
- Runs through a background daemon (`sshtmd`) that the CLI (`sshtm`) communicates with over gRPC.

## Requirements

- **Go 1.25 or later** for building from source.
- **Git** for cloning the repository.

- **systemd user services** on Linux or **LaunchAgents** on macOS for the installed daemon service.
- **Air** *(optional)* for live reload during development.

## Installation

### Build from source

```sh
git clone https://github.com/besrabasant/ssh-tunnel-manager.git
cd ssh-tunnel-manager
./install.sh
```

The install script builds the daemon and CLI for your OS/architecture, then installs them under `~/.local/bin`.

Installed paths:

| Purpose | Path |
| --- | --- |
| CLI binary | `~/.local/bin/sshtm` |
| Daemon binary | `~/.local/bin/sshtmd` |
| Data/config directory | `~/.ssh-tunnel-manager` |
| Helper scripts | `~/.local/share/sshtm/scripts` |

Service behavior:

- **Linux**: installs and starts a systemd user unit for `sshtmd`.
- **macOS**: installs and starts a LaunchAgent for `sshtmd`.

Make sure `~/.local/bin` is on your `PATH` after installation.

### Arch Linux (AUR)

If you use `yay`, install directly from the AUR:

```sh
yay -S sshtm
```

You can also build the package manually with `makepkg` inside `packaging/arch`.

### Debian/Ubuntu package build

Build a Debian package using the files in `packaging/debian`:

```sh
cd packaging/debian
dpkg-buildpackage -us -uc
sudo dpkg -i ../sshtm_1.1.5-1_amd64.deb
```

## Uninstallation

If installed from source, run:

```sh
~/.local/share/sshtm/scripts/uninstall.sh
```

## Quick start

1. Add a tunnel configuration interactively:
   ```sh
   sshtm add
   ```

2. List saved configurations:
   ```sh
   sshtm list
   ```

3. Start a saved tunnel:
   ```sh
   sshtm tunnel my_configuration
   ```

4. Start a saved tunnel with an explicit local port:
   ```sh
   sshtm tunnel my_configuration 8080
   ```

5. View active tunnels:
   ```sh
   sshtm active
   ```

6. Stop a tunnel by name or local port:
   ```sh
   sshtm kill my_configuration
   sshtm kill 8080
   ```

## Configuration fields

When adding or editing a tunnel, `sshtm` stores these fields:

| Field | Description |
| --- | --- |
| `name` | Unique name for the tunnel configuration. |
| `description` | Human-readable notes about the tunnel. |
| `server` | SSH server or jump host. |
| `user` | SSH username. |
| `key_file` | Private key file used for authentication. |
| `remote_host` | Destination host reached from the SSH server. |
| `remote_port` | Destination port on the remote host. |
| `local_port` | Local port to bind when starting the tunnel. |

## Commands

```sh
sshtm [command]
```

| Command | Aliases | Description |
| --- | --- | --- |
| `list [search pattern]` | `ls`, `l` | List saved SSH tunnel configurations, optionally filtered with fuzzy search. |
| `add` | `a` | Add a new SSH tunnel configuration using an interactive form. |
| `edit` | `e` | Edit an existing SSH tunnel configuration. |
| `delete` | `del`, `d` | Delete an existing SSH tunnel configuration. |
| `tunnel <configuration name> [local port]` | `t` | Start an SSH tunnel from a saved configuration, optionally overriding the local port. |
| `active` | | List active SSH tunnels. |
| `kill <configuration name or local port>` | `k`, `terminate` | Terminate an active SSH tunnel. |
| `completion` | | Generate shell completion scripts. |
| `version` | | Print the version number. |
| `help [command]` | | Show command help. |

Use `sshtm help <command>` for command-specific details.

## Development

### Generate protobuf code

```sh
make gen_proto
```

### Start the daemon with live reload

```sh
make gen_proto
air
```

This assumes [Air](https://github.com/air-verse/air) is installed and the command is run from the project root.

### Run the CLI from source

```sh
cd client
go run main.go [command] [arguments]
```

For example:

```sh
cd client
go run main.go list
```

### Build locally

```sh
make build
```

This creates:

- `./sshtmd` for the daemon.
- `./sshtm` for the CLI.

### Cross-platform builds

The Makefile includes targets for Linux and macOS:

```sh
make build_linux   # linux/amd64
make build_macos   # darwin/amd64 + darwin/arm64
make build_all     # linux + macos
```

You can pass an explicit version into builds:

```sh
make VERSION=1.1.5 build_all
```

### Run tests

```sh
go test ./...
```

## Project layout

```text
client/      CLI commands, terminal forms, and daemon client helpers
daemon/      Background daemon and gRPC server task handlers
pkg/         Reusable configuration and tunnel-management packages
rpc/         Protobuf definitions and generated gRPC code
config/      Application constants and version metadata
packaging/   Arch Linux and Debian packaging files
scripts/     Install/uninstall helper scripts
utils/       Shared utility helpers
```

## Versioning policy

This project uses [Semantic Versioning](https://semver.org/) for version numbers.

## License

This project is licensed under the [MIT License](LICENSE).
