#for build MacOS/iOS dylib
echo "building MacOS/iOS dylib"
#iOS
# shellcheck disable=SC2155
export CFLAGS="-arch arm64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphoneos --show-sdk-path)
# shellcheck disable=SC2155
export CGO_LDFLAGS="-arch arm64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphoneos --show-sdk-path)
CGO_ENABLED=1 GOARCH=arm64 GOOS=ios CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags ios -ldflags=-w -trimpath -o "./build/ios/gateway_arm64.a"  -buildmode c-archive main.go

# shellcheck disable=SC2155
export CFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
# shellcheck disable=SC2155
export CGO_LDFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
CGO_ENABLED=1 GOARCH=amd64 GOOS=ios CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags ios -ldflags=-w -trimpath -o "./build/ios/gateway_simulator_amd64.a"  -buildmode c-archive main.go

lipo -create ./build/ios/gateway_arm64.a ./build/ios/gateway_simulator_amd64.a -output ./build/ios/libgateway_amd64_arm64.a

#Mac
# shellcheck disable=SC2155
export CFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
# shellcheck disable=SC2155
export CGO_LDFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
#CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -o "test.a" -buildmode c-archive
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -o "./build/macos/gateway_amd64.a" -buildmode c-archive main.go
CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -o "./build/macos/gateway_arm64.a" -buildmode c-archive main.go

lipo -create ./build/macos/gateway_amd64.a ./build/macos/gateway_arm64.a -output ./build/macos/libgateway_amd64_arm64.a
