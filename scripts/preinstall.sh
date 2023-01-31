# Initial package installation
install() {
  # Create the user and group
  getent group ovh-dyndns >/dev/null || groupadd -r ovh-dyndns || :
  getent passwd ovh-dyndns >/dev/null || useradd -r -g  ovh-dyndns -s /sbin/nologin \
     -c "User for the OVH-DynDNS"  ovh-dyndns || :
}

# Package upgrade
upgrade() {
  :
}

action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  # deb passes $1=configure
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
  # deb passes $1=configure $2=<current version>
  action="upgrade"
fi

case "$action" in
  "1" | "install")
    install
    ;;
  "2" | "upgrade")
    upgrade
    ;;
esac