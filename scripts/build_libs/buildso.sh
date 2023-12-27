#for build linux/android so file
echo "building linux/android so file"
#linux和Android共用动态链接库
export CGO_ENABLED=1
export GOARCH=amd64
export GOOS=linux
go build -tags linux -ldflags=-w -trimpath -o build/linux/libgateway_amd64.so -buildmode c-shared main.go
# shellcheck disable=SC2034
export CGO_ENABLED=1
export GOARCH=arm64
export GOOS=linux
#sudo apt install gcc-aarch64-linux-gnu
export CC=aarch64-linux-gnu-gcc
##sudo apt install g++-aarch64-linux-gnu
#export CXX=aarch64-linux-gnu-g++
##sudo apt-get install binutils-aarch64-linux-gnu
#export AR=aarch64-linux-gnu-ar
go build -tags linux -ldflags=-w -trimpath -o build/linux/libgateway_arm64.so -buildmode c-shared main.go

#export PATH=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64/bin:~/Android/Sdk/ndk/25.2.9519653/toolchains/llvm/prebuilt/linux-x86_64/bin:$PATH
##android/app/src/main/jniLibs/armeabi-v7a
## shellcheck disable=SC2034
#export CGO_ENABLED=1
#export GOARCH=arm
#export GOOS=android
#export CC=armv7a-linux-androideabi33-clang
#go build -tags android -ldflags=-w -trimpath -o build/android/armeabi-v7a/libgateway.so -buildmode c-shared main.go
##android/app/src/main/jniLibs/arm64-v8a
## shellcheck disable=SC2034
#export CGO_ENABLED=1
#export GOARCH=arm64
#export GOOS=android
#export CC=aarch64-linux-android33-clang
#go build -tags android -ldflags=-w -trimpath -o build/android/arm64-v8a/libgateway.so -buildmode c-shared main.go
##android/app/src/main/jniLibs/x86
## shellcheck disable=SC2034
#export CGO_ENABLED=1
#export GOARCH=386
#export GOOS=android
#export CC=i686-linux-android33-clang
#go build -tags android -ldflags=-w -trimpath -o build/android/x86/libgateway.so -buildmode c-shared main.go
##android/app/src/main/jniLibs/x86_64
## shellcheck disable=SC2034
#export CGO_ENABLED=1
#export GOARCH=amd64
#export GOOS=android
#export CC=x86_64-linux-android33-clang
#go build -tags android -ldflags=-w -trimpath -o build/android/x86_64/libgateway.so -buildmode c-shared main.go