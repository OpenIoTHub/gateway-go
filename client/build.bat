gomobile bind -target=android  -javapkg=cloud.iothub
:: gpg --clearsign gateway-0.0.1.aar
:: mvn gpg:sign-and-deploy-file -Durl=https://s01.oss.sonatype.org/service/local/staging/deploy/maven2/ -DrepositoryId=ossrh -Dpackaging=aar -DpomFile=gateway-0.0.2.pom -Dfile=gateway-0.0.2.aar -Dsources=gateway-0.0.2.jar -Djavadoc=gateway-0.0.2.jar
:: mvn deploy:deploy-file -Dfile=gateway-0.0.2.aar -DgroupId=cloud.iothub -DartifactId=gateway -Dversion=0.0.1 -Dpackaging=aar -DrepositoryId=github -Durl=https://maven.pkg.github.com/OpenIoTHub/gateway-go
gomobile bind -ldflags '-w -s -extldflags "-lresolve"' --target=ios,macos,iossimulator
