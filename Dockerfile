FROM alpine:latest

#默认的http api端口
EXPOSE 1082
#mdns端口
EXPOSE 5353/udp

RUN apk add --no-cache bash

ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY gateway-go /bin/gateway-go
