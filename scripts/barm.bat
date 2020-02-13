set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
set GO111MODULE=on
::x86块
set GOARCH=arm
set GOOS=linux
go build -ldflags -w ../
ren gateway-go arm
upx arm
::upx arm
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause