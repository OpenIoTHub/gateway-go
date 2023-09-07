gomobile bind -target=android
gomobile bind -ldflags '-w -s -extldflags "-lresolve"' --target=ios,macos,iossimulator
