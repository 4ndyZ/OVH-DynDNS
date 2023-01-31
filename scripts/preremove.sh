# Package uninstall
uninstall() {
  systemctl stop ovh-dyndns.service || :
  if [ -x "/usr/lib/systemd/systemd-update-helper" ]; then
    /usr/lib/systemd/systemd-update-helper remove-system-units ovh-dyndns.service || :
  fi
}

# Package uninstall and purge
purge() {
  uninstall
}

case "$1" in
  "0" | "remove")
    uninstall
    ;;
  "1" | "upgrade")
    ;;
  "purge")
    purge
    ;;
esac