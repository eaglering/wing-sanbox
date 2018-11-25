FROM python:3.4-alpine
MAINTAINER eaglering <442958506@qq.com>

ADD script.sh /usr/local/bin/script.sh
RUN chmod 777 /usr/local/bin/script.sh \
	&& apk add --no-cache bc && rm -rf /var/cache/apk/*

ENTRYPOINT [ "sh", "-c" ]