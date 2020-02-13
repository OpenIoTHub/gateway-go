set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
set GO111MODULE=on
::x86块
set GOARCH=386
set GOOS=windows
go build -ldflags -w ../
ren gateway-go.exe windows386.exe
::upx windows386.exe
set GOOS=linux
go build -ldflags -w ../
ren gateway-go linux386
upx linux386
set GOOS=freebsd
go build -ldflags -w ../
ren gateway-go freebsd386
upx freebsd386
set GOOS=darwin
go build -ldflags -w ../
ren gateway-go darwin386
upx darwin386
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::x64块
set GOARCH=amd64
set GOOS=windows
go build -ldflags -w ../
ren gateway-go.exe windowsAmd64.exe
::upx windowsAmd64.exe
set GOOS=linux
go build -ldflags -w ../
ren gateway-go linuxAMD64
upx linuxAMD64
set GOOS=freebsd
go build -ldflags -w ../
ren gateway-go freebsdAMD64
upx freebsdAMD64
set GOOS=darwin
go build -ldflags -w ../
ren gateway-go darwinAMD64
upx darwinAMD64
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
set GOOS=linux
go build -ldflags -w ../
ren gateway-go LinuxArm
upx LinuxArm
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::mips块
set GOARCH=mips64le
set GOOS=linux
go build -ldflags -w ../
ren gateway-go LinuxMips64le
upx LinuxMips64le

set GOARCH=mips64
set GOOS=linux
go build -ldflags -w ../
ren gateway-go LinuxMips64
upx LinuxMips64

set GOARCH=mipsle
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w ../
ren gateway-go LinuxMipsle
upx LinuxMipsle

set GOARCH=mips
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w ../
ren gateway-go LinuxMips
upx LinuxMips
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause