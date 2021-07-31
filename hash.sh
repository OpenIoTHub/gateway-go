#提交openwrt版获取hash
PKG_VERSION=0.1.92
URL=https://codeload.github.com/OpenIoTHub/gateway-go/tar.gz/v${PKG_VERSION}
wget ${URL}
openssl sha256 v${PKG_VERSION}
rm v${PKG_VERSION}
#SHA256(v0.1.92.tar.gz)= dd8074d9312e00ff957ffd1f3be7ba118a9b8cc31f07aa1ed594ef07931dab16