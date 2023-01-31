# Package uninstall
uninstall() {
  userdel -f ovh-dyndns >/dev/null || :
  systemctl daemon-reload || :
}

# Package uninstall and purge
purge() {
  rm -drf /etc/ovh-dyndns || :
  rm -drf /var/log/ovh-dyndns || :
}

# Package upgrade
upgrade() {
  :
}

action="$1"
case "$action" in
  "0" | "remove")
    uninstall
    ;;
  "1" | "upgrade")
    upgrade
    ;;
  "purge")
    uninstall
    purge
    ;;
esac
