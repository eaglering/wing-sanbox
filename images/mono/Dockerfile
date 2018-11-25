FROM mono:5.12
MAINTAINER eaglering <442958506@qq.com>

WORKDIR /data
ADD script.sh /usr/local/bin/script.sh
RUN chmod 777 /usr/local/bin/script.sh \
	&& echo "deb http://ftp.de.debian.org/debian jessie main" > /etc/apt/sources.list.d/bc.list \
	&& apt-get update \
	&& apt-get install -y bc \
	&& rm -rf /var/lib/apt/lists/* /tmp/*

ENTRYPOINT [ "bash", "-c" ]