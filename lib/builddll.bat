::for build windows dll
echo "building windows dll"
::brew install mingw-w64
::sudo apt-get install binutils-mingw-w64
::https://winlibs.com/#download-release
SET PATH=C:\mingw64\bin;%PATH%
SET CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64
go build -o ./build/amd64/windows/gateway.dll -buildmode=c-shared