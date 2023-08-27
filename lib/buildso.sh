#for build linux/android so file
echo "building linux/android so file"
#linux和Android共用动态链接库
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -tags linux -ldflags=-w -trimpath -v -o "build/linux/gateway_amd64.so" -buildmode c-shared
#sudo apt-get install binutils-aarch64-linux-gnu
# shellcheck disable=SC2034
CGO_ENABLED=1 GOARCH=arm64 GOOS=linux CC=aarch64-linux-gnu-gcc && \
CXX=aarch64-linux-gnu-g++ AR=aarch64-linux-gnu-ar go build -tags linux -ldflags=-w -trimpath -v -o "build/linux/gateway_arm64.so" -buildmode c-shared

#android/app/src/main/jniLibs/armeabi-v7a
# shellcheck disable=SC2034
CGO_ENABLED=1 GOARCH=arm GOOS=android CC=aarch64-linux-gnu-gcc && \
CXX=aarch64-linux-gnu-g++ AR=aarch64-linux-gnu-ar go build -tags android -ldflags=-w -trimpath -v -o "build/android/gateway_arm.so" -buildmode c-shared
#android/app/src/main/jniLibs/arm64-v8a
# shellcheck disable=SC2034
CGO_ENABLED=1 GOARCH=arm64 GOOS=android CC=aarch64-linux-gnu-gcc && \
CXX=aarch64-linux-gnu-g++ AR=aarch64-linux-gnu-ar go build -tags android -ldflags=-w -trimpath -v -o "build/android/gateway_arm64.so" -buildmode c-shared
#android/app/src/main/jniLibs/x86
CGO_ENABLED=1 GOARCH=386 GOOS=android go build -tags android -ldflags=-w -trimpath -v -o "build/android/gateway_i386.so" -buildmode c-shared
#android/app/src/main/jniLibs/x86_64
CGO_ENABLED=1 GOARCH=amd64 GOOS=android go build -tags android -ldflags=-w -trimpath -v -o "build/android/gateway_amd64.so" -buildmode c-shared