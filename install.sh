#!/bin/bash
set -euo pipefail

# Paths (all user-scoped)
DATA_DIR="$HOME/.ssh-tunnel-manager" # application data/config root

# Detect OS/arch for service management and binary selection
OS_NAME="$(uname -s)"
ARCH_NAME="$(uname -m)"
case "$ARCH_NAME" in
	x86_64|amd64)
		ARCH_NAME="amd64"
		;;
	arm64|aarch64)
		ARCH_NAME="arm64"
		;;
esac
if [ "$OS_NAME" = "Darwin" ]; then
	LAUNCH_AGENTS_DIR="$HOME/Library/LaunchAgents"   # launchd user agents directory
	SHARE_DIR="$HOME/Library/Application Support/sshtm" # application shared assets root
	BIN_DIR="$HOME/.local/bin"                       # user binaries directory
	DAEMON_BIN="$BIN_DIR/sshtmd"                     # sshtm daemon binary
	CLIENT_BIN="$BIN_DIR/sshtm"                      # sshtm client binary
	SCRIPTS_DIR="$SHARE_DIR/scripts"                 # bundled scripts directory
	PLIST_FILE="$LAUNCH_AGENTS_DIR/com.sshtm.daemon.plist" # launchd plist path
else
	SYSTEMD_USER_DIR="$HOME/.config/systemd/user"  # systemd user unit directory
	SHARE_DIR="$HOME/.local/share/sshtm"           # application shared assets root
	BIN_DIR="$HOME/.local/bin"                     # user binaries directory
	DAEMON_BIN="$BIN_DIR/sshtmd"                   # sshtm daemon binary
	CLIENT_BIN="$BIN_DIR/sshtm"                    # sshtm client binary
	SCRIPTS_DIR="$SHARE_DIR/scripts"               # bundled scripts directory
	UNIT_FILE="$SYSTEMD_USER_DIR/sshtmd.service"   # systemd unit file path
fi

# Ensure protoc gRPC plugin is available
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Build binaries for the detected OS/arch
TARGET_OS=""
if [ "$OS_NAME" = "Darwin" ]; then
	TARGET_OS="darwin"
else
	TARGET_OS="linux"
fi
SRC_DAEMON="sshtmd-${TARGET_OS}-${ARCH_NAME}"
SRC_CLIENT="sshtm-${TARGET_OS}-${ARCH_NAME}"
GOOS="$TARGET_OS" GOARCH="$ARCH_NAME" go build -o "$SRC_DAEMON" ./daemon
GOOS="$TARGET_OS" GOARCH="$ARCH_NAME" go build -o "$SRC_CLIENT" ./client

# Create required directories
if [ "$OS_NAME" = "Darwin" ]; then
	mkdir -p "$LAUNCH_AGENTS_DIR"
else
	mkdir -p "$SYSTEMD_USER_DIR"
fi
mkdir -p "$DATA_DIR"
mkdir -p "$SHARE_DIR"
mkdir -p "$BIN_DIR"

# Stop existing service before replacing binaries
if [ "$OS_NAME" = "Darwin" ]; then
	if launchctl list | grep -q "com.sshtm.daemon"; then
		launchctl unload "$PLIST_FILE" || true
	fi
else
	if systemctl --user is-active --quiet sshtmd; then
		systemctl --user stop sshtmd
	fi
fi

# Install/overwrite binaries in-place (updates existing installations)
install -m 755 -T "$SRC_DAEMON" "$DAEMON_BIN"
install -m 755 -T "$SRC_CLIENT" "$CLIENT_BIN"

# Refresh bundled scripts directory
rm -rf "$SCRIPTS_DIR"
cp -r scripts "$SHARE_DIR/"

# Register and start the daemon as a user service
if [ "$OS_NAME" = "Darwin" ]; then
	cat <<EOF > "$PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.sshtm.daemon</string>
	<key>ProgramArguments</key>
	<array>
		<string>$DAEMON_BIN</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
</dict>
</plist>
EOF

	# Load/enable launchd agent
	launchctl load "$PLIST_FILE"
else
	cat <<EOF > "$UNIT_FILE"
[Unit]
Description=SSH Tunnel Manager Daemon

[Service]
Type=simple
ExecStart=$DAEMON_BIN

[Install]
WantedBy=default.target
EOF

	# Reload, enable, and start systemd user service
	systemctl --user daemon-reload

	# Enable and start the service
	systemctl --user enable sshtmd

	# If the service is already running, restart it. Otherwise start it normally.
	if systemctl --user is-active --quiet sshtmd; then
		systemctl --user restart sshtmd
	else
		systemctl --user start sshtmd
	fi
fi
