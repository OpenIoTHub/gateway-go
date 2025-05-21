# avahi in docker
docker run --privileged -v /var/run/dbus:/var/run/dbus -v /var/run/avahi-daemon/socket:/var/run/avahi-daemon/socket -it openiothub/gateway-go:latest bash
apt-get update && apt-get install avahi-utils iputils-ping -y
ping whatever.local
avahi-browse -a