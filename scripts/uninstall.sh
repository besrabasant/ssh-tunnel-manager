#!/bin/bash

# stop and disable the service
systemctl --user stop sshtmd
systemctl --user disable sshtmd

rm -f $HOME/.local/bin/sshtmd
rm -f $HOME/.local/bin/sshtm
rm -f $HOME/.config/systemd/user/sshtmd.service
rm -rf $HOME/.local/share/sshtm

systemctl --user daemon-reload


