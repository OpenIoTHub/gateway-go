set CGO_ENABLED=1
set GOROOT_BOOTSTRAP=C:/Go
set GO111MODULE=on
set GOARCH=mipsle
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w ../
ren gateway-go LinuxMipsle
upx -9 LinuxMipsle

set GOARCH=amd64
set GOOS=windows
pause