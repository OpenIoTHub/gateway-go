#for build MacOS/iOS dylib
echo "building MacOS/iOS dylib"
#iOS
export CFLAGS="-arch arm64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphoneos --show-sdk-path)
export CGO_LDFLAGS="-arch arm64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphoneos --show-sdk-path)
CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags ios -ldflags=-w -trimpath -v -o "./build/arm64/ios/gateway.a" -buildmode c-archive

export CFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
export CGO_LDFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags ios -ldflags=-w -trimpath -v -o "./build/amd64/ios/gateway.a" -buildmode c-archive

#Mac
export CFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
export CGO_LDFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
#CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -v -o "test.a" -buildmode c-archive
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -v -o "./build/amd64/macos/gateway.dylib" -buildmode c-shared
CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -v -o "./build/arm64/macos/gateway.dylib" -buildmode c-shared

