rm -rf %GODIST%/natcloud/client/*
set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOARCH=386
set GOOS=windows
go build -ldflags -w main.go
ren main.exe clientWindows386.exe
::upx windows386.exe
mv clientWindows386.exe %GODIST%/natcloud/client/
set GOOS=linux
go build -ldflags -w main.go
ren main clientLinux386
upx clientLinux386
mv clientLinux386 %GODIST%/natcloud/client/
set GOOS=freebsd
go build -ldflags -w main.go
ren main clientFreebsd386
upx clientFreebsd386
mv clientFreebsd386 %GODIST%/natcloud/client/
set GOOS=darwin
go build -ldflags -w main.go
ren main clientDarwin386
upx clientDarwin386
mv clientDarwin386 %GODIST%/natcloud/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::x64块
set GOARCH=amd64
set GOOS=windows
go build -ldflags -w main.go
ren main.exe clientWindowsAmd64.exe
::upx windowsAmd64.exe
mv clientWindowsAmd64.exe %GODIST%/natcloud/client/

set GOOS=linux
go build -ldflags -w main.go
ren main clientLinuxAMD64
upx clientLinuxAMD64
mv clientLinuxAMD64 %GODIST%/natcloud/client/

set GOOS=freebsd
go build -ldflags -w main.go
ren main clientFreebsdAMD64
upx clientFreebsdAMD64
mv clientFreebsdAMD64 %GODIST%/natcloud/client/

set GOOS=darwin
go build -ldflags -w main.go
ren main clientDarwinAMD64
upx clientDarwinAMD64
mv clientDarwinAMD64 %GODIST%/natcloud/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
set GOOS=linux
go build -ldflags -w main.go
ren main clientLinuxArm
upx clientLinuxArm
mv clientLinuxArm %GODIST%/natcloud/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::mips块
set GOARCH=mips64le
set GOOS=linux
go build -ldflags -w main.go
ren main clientLinuxMips64le
upx clientLinuxMips64le
mv clientLinuxMips64le %GODIST%/natcloud/client/

set GOARCH=mips64
set GOOS=linux
go build -ldflags -w main.go
ren main clientLinuxMips64
upx clientLinuxMips64
mv clientLinuxMips64 %GODIST%/natcloud/client/

set GOARCH=mipsle
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w main.go
ren main clientLinuxMipsle
upx clientLinuxMipsle
mv clientLinuxMipsle %GODIST%/natcloud/client/

set GOARCH=mips
set GOOS=linux
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w main.go
ren main clientLinuxMips
upx clientLinuxMips
mv clientLinuxMips %GODIST%/natcloud/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause