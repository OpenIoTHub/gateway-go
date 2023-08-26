#for build linux/android so file
echo "building linux/android so file"
#linux和Android共用动态链接库
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -tags linux -ldflags=-w -trimpath -v -o "build/amd64/linux/gateway.so" -buildmode c-shared
#sudo apt-get install binutils-aarch64-linux-gnu
CGO_ENABLED=1 GOARCH=arm64 GOOS=linux CC=aarch64-linux-gnu-gcc && \
CXX=aarch64-linux-gnu-g++ AR=aarch64-linux-gnu-ar go build -tags linux -ldflags=-w -trimpath -v -o "build/arm64/linux/gateway.so" -buildmode c-shared