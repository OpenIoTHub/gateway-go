::for build windows dll
echo "building windows dll"
::brew install mingw-w64
::sudo apt-get install binutils-mingw-w64
::https://winlibs.com/#download-release
SET PATH=C:\mingw64\bin;%PATH%
SET CGO_ENABLED=1
SET CC=x86_64-w64-mingw32-gcc
SET CXX=x86_64-w64-mingw32-g++
SET GOOS=windows
SET GOARCH=amd64
go build -tags windows -ldflags=-w -trimpath -o ./build/windows/gateway_amd64.dll -buildmode=c-shared lib/lib.go