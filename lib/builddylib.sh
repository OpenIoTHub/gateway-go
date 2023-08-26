#for build MacOS/iOS dylib
echo "building MacOS/iOS dylib"
#iOS
# shellcheck disable=SC2155
export CFLAGS="-arch arm64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphoneos --show-sdk-path)
# shellcheck disable=SC2155
export CGO_LDFLAGS="-arch arm64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphoneos --show-sdk-path)
CGO_ENABLED=1 GOARCH=arm64 GOOS=ios CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags ios -ldflags=-w -trimpath -v -o "./build/arm64/ios/gateway.a"  -buildmode c-archive

# shellcheck disable=SC2155
export CFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
# shellcheck disable=SC2155
export CGO_LDFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
CGO_ENABLED=1 GOARCH=amd64 GOOS=ios CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags ios -ldflags=-w -trimpath -v -o "./build/amd64/ios/gateway.a"  -buildmode c-archive

mkdir ./build/ios
lipo -create ./build/arm64/ios/gateway.a ./build/amd64/ios/gateway.a -output ./build/ios/gateway.a

#Mac
# shellcheck disable=SC2155
export CFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
# shellcheck disable=SC2155
export CGO_LDFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
#CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -v -o "test.a" -buildmode c-archive
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -v -o "./build/amd64/macos/gateway.dylib" -buildmode c-shared
CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -v -o "./build/arm64/macos/gateway.dylib" -buildmode c-shared

