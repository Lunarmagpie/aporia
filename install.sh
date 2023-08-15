#!/bin/bash
set -e

if [ $(/usr/bin/id -u $USER) -ne 0 ]; then
  echo "This script must be run as root."
  exit 1
fi

if
  (sudo -u $SUDO_USER -s go build) &&  # Allow user to have go installed in userspace, ex asdf
  cp extra/aporia.pam /etc/pam.d/aporia &&
  cp extra/aporia.service /etc/systemd/system/aporia.service &&
  (mkdir /etc/aporia || true) &&
  (mkdir /etc/aporia/.scripts || true) &&
  cp extra/startx.sh /etc/aporia/.scripts/startx.sh &&
  chmod +x /etc/aporia/.scripts/startx.sh &&
  cp extra/xsetup.sh /etc/aporia/.scripts/xsetup.sh &&
  chmod +x /etc/aporia/.scripts/xsetup.sh &&
  cp aporia /bin/aporia -f  # this one is pretty bad
then
  set +e
  echo -e "Run the command \033[94msystemctl enable aporia.service\033[39;49m to enable aporia"
  echo " | See more information at https://github.com/lunarmagpie/aporia"
else
  set +e
  echo -e "\033[31mAporia failed to install\033[39;49m"
  echo " | See more information at https://github.com/lunarmagpie/aporia"
  exit 1
fi
