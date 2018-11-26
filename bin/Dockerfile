FROM alpine:latest
MAINTAINER ealgering <442958506@qq.com>

ENV SERVER_PEM /data/server.pem
ENV SERVER_KEY /data/server.key

RUN apk add --no-cache bash docker \
    && rm -rf /var/cache/apk/*

ADD wing.sh /usr/local/bin/wing.sh
ADD wing_server_release_linux /usr/local/bin/wing
ADD server.pem ${SERVER_PEM}
ADD server.key ${SERVER_KEY}

RUN chmod 755 /usr/local/bin/wing.sh \
        && chmod 755 /usr/local/bin/wing

CMD ["/usr/local/bin/wing"]
