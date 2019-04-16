rm -rf %GODIST%/natcloud/snap/client/*
set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOOS=linux

set GOARCH=386
go build -ldflags -w main.go
ren main clientLinux386
upx clientLinux386
mv clientLinux386 %GODIST%/natcloud/snap/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::x64块
set GOARCH=amd64
go build -ldflags -w main.go
ren main clientLinuxAMD64
upx clientLinuxAMD64
mv clientLinuxAMD64 %GODIST%/natcloud/snap/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::arm块
set GOARCH=arm
go build -ldflags -w main.go
ren main clientLinuxArm
upx clientLinuxArm
mv clientLinuxArm %GODIST%/natcloud/snap/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
  
::mips块
set GOARCH=mips64le
go build -ldflags -w main.go
ren main clientLinuxMips64le
upx clientLinuxMips64le
mv clientLinuxMips64le %GODIST%/natcloud/snap/client/

set GOARCH=mips64
go build -ldflags -w main.go
ren main clientLinuxMips64
upx clientLinuxMips64
mv clientLinuxMips64 %GODIST%/natcloud/snap/client/

set GOARCH=mipsle
set CGO_ENABLED=0
set GOMIPS=softfloat
go build -ldflags -w main.go
ren main clientLinuxMipsle
upx clientLinuxMipsle
mv clientLinuxMipsle %GODIST%/natcloud/snap/client/
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause