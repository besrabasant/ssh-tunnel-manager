#!/bin/bash

make

mkdir -p $HOME/.config/systemd/user
mkdir -p $HOME/.local/share/sshtm

mv sshtmd $HOME/.local/bin/
mv sshtm $HOME/.local/bin/
cp -r scripts $HOME/.local/share/sshtm

# Optional: Set up the application as a system service (for systems using systemd)
cat <<EOF > $HOME/.config/systemd/user/sshtmd.service
[Unit]
Description=SSH Tunnnel manager Daemon

[Service]
Type=simple
ExecStart=$HOME/.local/bin/sshtmd

[Install]
WantedBy=default.target
EOF

systemctl --user daemon-reload

# Enable and start the service
systemctl --user enable sshtmd
systemctl --user start sshtmd
