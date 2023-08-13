#!/bin/bash

go build

cp extra/aporia.pam /etc/pam/aporia
cp extra/aporia.service /etc/systemd/system/aporia.service
cp aporia /bin/aporia

echo "It is ok this command fails"
systemctl disable display-manager.service
systemctl enable aporia.service
