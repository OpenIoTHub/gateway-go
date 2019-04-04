set CGO_ENABLED=0
set GOROOT_BOOTSTRAP=C:/Go
::x86块
set GOARCH=386
set GOOS=windows
go build -ldflags -w main.go
ren main.exe windows386.exe
upx windows386.exe
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
set GOARCH=amd64
set GOOS=windows
pause