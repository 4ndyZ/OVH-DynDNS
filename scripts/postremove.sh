#!/usr/bin/env sh
systemctl stop ovh-dyndns
userdel -f ovh-dyndns >/dev/null
exit 0
