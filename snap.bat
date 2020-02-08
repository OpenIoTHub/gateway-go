rm -rf %GODIST%/natcloud/snap/bin/client/*
set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOOS=linux

set GOARCH=386
go build -ldflags -w main.go
ren main clientLinux386
upx clientLinux386
mv clientLinux386 %GODIST%/natcloud/snap/bin/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
go build -ldflags -w main.go
ren main clientLinuxArm
upx clientLinuxArm
mv clientLinuxArm %GODIST%/natcloud/snap/bin/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause