export CGO_ENABLED=0
export GO111MODULE=on
export GOARCH=386
export GOOS=windows
go build -ldflags -w ../
mv gateway-go.exe windows386.exe
upx windows386.exe
export GOOS=linux
go build -ldflags -w ../
mv gateway-go linux386
upx linux386
export GOOS=freebsd
go build -ldflags -w ../
mv gateway-go freebsd386
upx freebsd386
export GOOS=darwin
go build -ldflags -w ../
mv gateway-go darwin386
upx darwin386

export GOARCH=amd64
export GOOS=windows
go build -ldflags -w ../
mv gateway-go.exe windowsAmd64.exe
upx windowsAmd64.exe
export GOOS=linux
go build -ldflags -w ../
mv gateway-go linuxAMD64
upx linuxAMD64
export GOOS=freebsd
go build -ldflags -w ../
mv gateway-go freebsdAMD64
upx freebsdAMD64
export GOOS=darwin
go build -ldflags -w ../
mv gateway-go darwinAMD64
upx darwinAMD64

export GOARCH=arm
export GOOS=linux
go build -ldflags -w ../
mv gateway-go LinuxArm
upx LinuxArm

export GOARCH=mips64le
export GOOS=linux
go build -ldflags -w ../
mv gateway-go LinuxMips64le
upx LinuxMips64le

export GOARCH=mips64
export GOOS=linux
go build -ldflags -w ../
mv gateway-go LinuxMips64
upx LinuxMips64

export GOARCH=mipsle
export GOOS=linux
export CGO_ENABLED=0
export GOMIPS=softfloat
go build -ldflags -w ../
mv gateway-go LinuxMipsle
upx LinuxMipsle

export GOARCH=mips
export GOOS=linux
export CGO_ENABLED=0
export GOMIPS=softfloat
go build -ldflags -w ../
mv gateway-go LinuxMips
upx LinuxMips

export GOARCH=amd64
export GOOS=windows