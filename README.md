# gateway-go
[![Build Status](https://travis-ci.com/OpenIoTHub/gateway-go.svg?branch=master)](https://travis-ci.com/OpenIoTHub/gateway-go)

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/gateway-go)

You can install the pre-compiled binary (in several different ways),
use Docker.

Here are the steps for each of them:

## Install the pre-compiled binary

**openwrt/entware/optware (Usually on the router)**:
#### use snapshot branch：https://downloads.openwrt.org/snapshots/
```sh
opkg update
opkg install gateway-go
```

**homebrew tap** :

```sh
$ brew install OpenIoTHub/tap/gateway-go
```

**homebrew** (may not be the latest version):

```sh
$ brew install gateway-go
```
homebrew pr [gateway-go](https://github.com/Homebrew/homebrew-core/blob/master/Formula/gateway-go.rb)
```text
*** config file : 
/usr/local/etc/gateway-go/gateway-go.yaml
```


**snapcraft**:

```sh
$ sudo snap install gateway-go
```
```text
*** config file :
 /root/snap/gateway-go/current/gateway-go.yaml
```


**scoop**:

```sh
$ scoop bucket add OpenIoTHub https://github.com/OpenIoTHub/scoop-bucket.git
$ scoop install gateway-go
```

**deb/rpm**:

Download the `.deb` or `.rpm` from the [releases page][releases] and
install with `dpkg -i` and `rpm -i` respectively.
```text
*** config file :
 /etc/gateway-go/gateway-go.yaml
```

**manually**:

Download the pre-compiled binaries from the [releases page][releases] and
copy to the desired location.

## Running with Docker

You can also use it within a Docker container. To do that, you'll need to
execute something more-or-less like the following:

```sh
$ docker run -it --net=host openiothub/gateway-go:latest -t <your Token>
```

Note that the image will almost always have the last stable Go version.

[releases]: https://github.com/OpenIoTHub/gateway-go/releases
