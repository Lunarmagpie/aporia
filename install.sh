#!/bin/bash
set -e

if
  go build
  cp extra/aporia.pam /etc/pam.d/aporia
  cp extra/aporia.service /etc/systemd/system/aporia.service
  cp aporia /bin/aporia
  mkdir /etc/aporia
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
