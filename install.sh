#!/bin/bash

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

make

mkdir -p $HOME/.config/systemd/user
mkdir -p $HOME/.local/share/sshtm

mv sshtmd $HOME/.local/bin/
mv sshtm $HOME/.local/bin/
cp -r scripts $HOME/.local/share/sshtm

# Optional: Set up the application as a system service (for systems using systemd)
cat <<EOF > $HOME/.config/systemd/user/sshtmd.service
[Unit]
Description=SSH Tunnel Manager Daemon

[Service]
Type=simple
ExecStart=$HOME/.local/bin/sshtmd

[Install]
WantedBy=default.target
EOF

systemctl --user daemon-reload

# Enable and start the service
systemctl --user enable sshtmd

# If the service is already running, restart it. Otherwise start it normally.
if systemctl --user is-active --quiet sshtmd; then
    systemctl --user restart sshtmd
else
    systemctl --user start sshtmd
fi
