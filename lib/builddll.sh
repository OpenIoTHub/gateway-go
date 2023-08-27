#for build windows dll
echo "building windows dll"
#brew install mingw-w64
#sudo apt-get install binutils-mingw-w64
# shellcheck disable=SC2034
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 && \
go build -o ./build/windows/gateway_amd64.dll -buildmode=c-shared