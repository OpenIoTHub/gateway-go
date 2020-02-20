# gateway-go
[![Build Status](https://travis-ci.org/OpenIoTHub/gateway-go.svg?branch=master)](https://travis-ci.org/OpenIoTHub/gateway-go)

You can install the pre-compiled binary (in several different ways),
use Docker.

Here are the steps for each of them:

## Install the pre-compiled binary

**homebrew tap** (only on macOS for now):

```sh
$ brew install OpenIoTHub/tap/gateway-go
```

**homebrew** (may not be the latest version):

```sh
$ brew install gateway-go（not support yet）
```

**snapcraft**:

```sh
$ sudo snap install gateway-go
```

**scoop**:

```sh
$ scoop bucket add OpenIoTHub https://github.com/OpenIoTHub/scoop-bucket.git
$ scoop install gateway-go
```

**deb/rpm**:

Download the `.deb` or `.rpm` from the [releases page][releases] and
install with `dpkg -i` and `rpm -i` respectively.

**Shell script**:

```sh
$ curl -sfL https://install.goreleaser.com/github.com/OpenIoTHub/gateway-go.sh | sh
```

**manually**:

Download the pre-compiled binaries from the [releases page][releases] and
copy to the desired location.

## Running with Docker

You can also use it within a Docker container. To do that, you'll need to
execute something more-or-less like the following:

```sh
$ docker run openiothub/gateway:latest
```

Note that the image will almost always have the last stable Go version.

[releases]: https://github.com/goreleaser/goreleaser/releases
