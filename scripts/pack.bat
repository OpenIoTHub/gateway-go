rm -rf %GODIST%/natcloud/gateway/*
set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
set GO111MODULE=on
::x86块
set GOARCH=386
set GOOS=windows
go build -ldflags -w ../
ren gateway-go.exe gatewayWindows386.exe
::upx windows386.exe
mv gatewayWindows386.exe %GODIST%/natcloud/gateway/
set GOOS=linux
go build -ldflags -w ../
ren gateway-go gatewayLinux386
upx gatewayLinux386
mv gatewayLinux386 %GODIST%/natcloud/gateway/
set GOOS=freebsd
go build -ldflags -w ../
ren gateway-go gatewayFreebsd386
upx gatewayFreebsd386
mv gatewayFreebsd386 %GODIST%/natcloud/gateway/
set GOOS=darwin
go build -ldflags -w ../
ren gateway-go gatewayDarwin386
upx gatewayDarwin386
mv gatewayDarwin386 %GODIST%/natcloud/gateway/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::x64块
set GOARCH=amd64
set GOOS=windows
go build -ldflags -w ../
ren gateway-go.exe gatewayWindowsAmd64.exe
::upx windowsAmd64.exe
mv gatewayWindowsAmd64.exe %GODIST%/natcloud/gateway/

set GOOS=linux
go build -ldflags -w ../
ren gateway-go gatewayLinuxAMD64
upx gatewayLinuxAMD64
mv gatewayLinuxAMD64 %GODIST%/natcloud/gateway/

set GOOS=freebsd
go build -ldflags -w ../
ren gateway-go gatewayFreebsdAMD64
upx gatewayFreebsdAMD64
mv gatewayFreebsdAMD64 %GODIST%/natcloud/gateway/

set GOOS=darwin
go build -ldflags -w ../
ren gateway-go gatewayDarwinAMD64
upx gatewayDarwinAMD64
mv gatewayDarwinAMD64 %GODIST%/natcloud/gateway/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
set GOOS=linux
go build -ldflags -w ../
ren gateway-go gatewayLinuxArm
upx gatewayLinuxArm
mv gatewayLinuxArm %GODIST%/natcloud/gateway/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::mips块
set GOARCH=mips64le
set GOOS=linux
go build -ldflags -w ../
ren gateway-go gatewayLinuxMips64le
upx gatewayLinuxMips64le
mv gatewayLinuxMips64le %GODIST%/natcloud/gateway/

set GOARCH=mips64
set GOOS=linux
go build -ldflags -w ../
ren gateway-go gatewayLinuxMips64
upx gatewayLinuxMips64
mv gatewayLinuxMips64 %GODIST%/natcloud/gateway/

set GOARCH=mipsle
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w ../
ren gateway-go gatewayLinuxMipsle
upx gatewayLinuxMipsle
mv gatewayLinuxMipsle %GODIST%/natcloud/gateway/

set GOARCH=mips
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w ../
ren gateway-go gatewayLinuxMips
upx gatewayLinuxMips
mv gatewayLinuxMips %GODIST%/natcloud/gateway/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause