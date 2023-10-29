#for build windows dll
echo "building windows dll"
#brew install mingw-w64
#sudo apt-get install binutils-mingw-w64
# shellcheck disable=SC2034
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
export GOOS=windows GOARCH=amd64
go build -tags windows -ldflags=-w -trimpath -o ./build/windows/gateway_amd64.dll -buildmode=c-shared lib/lib.go