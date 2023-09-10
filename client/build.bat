gomobile bind -target=android  -javapkg=cloud.iothub
:: gpg --clearsign gateway-0.0.1.aar
:: mvn gpg:sign-and-deploy-file -Durl=https://s01.oss.sonatype.org/service/local/staging/deploy/maven2/ -DrepositoryId=ossrh -Dpackaging=aar -DpomFile=gateway-0.0.2.pom -Dfile=gateway-0.0.2.aar -Dsources=gateway-0.0.2.jar -Djavadoc=gateway-0.0.2.jar
:: mvn deploy:deploy-file -Dfile=client.aar -DgroupId=cloud.iothub -DartifactId=gateway -Dversion=0.0.1 -Dpackaging=aar -DrepositoryId=github -Durl=https://maven.pkg.github.com/OpenIoTHub/gateway-go
gomobile bind -ldflags '-w -s -extldflags "-lresolve"' --target=ios,macos,iossimulator
::gomobile bind -ldflags '-w -s -extldflags "-lresolve"' --target=ios,macos,iossimulator -o OpenIoTHubGateway.xcframework ./client
::https://gitee.com/OpenIoThub/mobile-lib-podspec
::git tag -a 0.0.1 -m '0.0.1'
::git pus --tags
::pod trunk push ./OpenIoTHubGateway.podspec --skip-import-validation --allow-warnings
