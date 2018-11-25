FROM java:openjdk-8-alpine
MAINTAINER eaglering <442958506@qq.com>

WORKDIR /data
ADD script.sh /usr/local/bin/script.sh
ADD java_runner.sh /usr/local/bin/java_runner.sh
RUN chmod 777 /usr/local/bin/script.sh \
	&& chmod 777 /usr/local/bin/java_runner.sh \
	&& apk add --no-cache bc && rm -rf /var/cache/apk/*

ENTRYPOINT [ "sh", "-c" ]