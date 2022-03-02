#!/usr/bin/env sh
getent group ovh-dyndns >/dev/null || \
	groupadd -r ovh-dyndns
getent passwd ovh-dyndns >/dev/null || \
	useradd -r -g ovh-dyndns -s /sbin/nologin \
    -c "User for the OVH-DynDNS" ovh-dyndns
exit 0
