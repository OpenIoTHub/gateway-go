export CGO_ENABLED=1
export GOROOT_BOOTSTRAP=C:/Go
export GO111MODULE=on
export GOARCH=mipsle
export GOOS=linux
export CGO_ENABLED=0
export GOMIPS=softfloat
go build -ldflags -w

export GOARCH=amd64
export GOOS=darwin