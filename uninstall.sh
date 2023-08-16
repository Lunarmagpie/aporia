#!/bin/bash
set -e

if [ $(/usr/bin/id -u $USER) -ne 0 ]; then
  echo "This script must be run as root."
  exit 1
fi

systemctl disable aporia.service
rm /etc/pam.d/aporia
rm /etc/systemd/system/aporia.service
rm -rf /etc/aporia

if [ -f /bin/aporia ]; then
  rm /bin/aporia
fi

echo "Aporia has been uninstalled."
