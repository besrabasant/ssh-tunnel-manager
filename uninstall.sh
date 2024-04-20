#!/bin/bash

# Ensure the script is run as root
if [ "$(id -u)" -ne 0 ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# stop and disable the service
systemctl stop sshtmd
systemctl disable sshtmd

rm -f /usr/local/bin/sshtmd
rm -f /usr/local/bin/sshtm
rm -f /etc/systemd/system/sshtmd.service

systemctl daemon-reload


