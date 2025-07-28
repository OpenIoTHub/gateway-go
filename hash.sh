#提交openwrt版获取hash
PKG_VERSION=2.0.10
URL=https://codeload.github.com/OpenIoTHub/gateway-go/tar.gz/v${PKG_VERSION}
wget ${URL}
openssl sha256 v${PKG_VERSION}
rm v${PKG_VERSION}
#SHA2-256(v2.0.10)= 95de453e76a22ae69a1bc9c32d6930278980c98decbe7a511fd9870bcdf90d2d

git commit --amend -m "gateway-go: update to 2.0.10" -s