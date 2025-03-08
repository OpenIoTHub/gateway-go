FROM alpine
LABEL name=gateway-go
LABEL url=https://github.com/OpenIoTHub/OpenIoTHub
RUN apk add --no-cache bash

WORKDIR /app
COPY gateway-go /app/
ENV TZ=Asia/Shanghai
#mdns端口
EXPOSE 5353/udp
EXPOSE 34323
ENTRYPOINT ["/app/gateway-go"]
CMD ["-c", "/root/config.yaml"]