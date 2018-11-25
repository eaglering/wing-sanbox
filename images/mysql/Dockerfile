FROM mysql:5.6
MAINTAINER eaglering <442958506@qq.com>

WORKDIR /data
ADD script.sh /usr/local/bin/script.sh
ADD sql_runner.sh /usr/local/bin/sql_runner.sh
ADD dist /usr/local/dist
RUN chmod 777 /usr/local/bin/script.sh \
	&& chmod 777 /usr/local/bin/sql_runner.sh \
	&& echo "deb http://ftp.de.debian.org/debian jessie main" > /etc/apt/sources.list.d/bc.list \
	&& apt-get update \
	&& apt-get install -y bc \
	&& rm -rf /var/lib/apt/lists/* /tmp/*

ENTRYPOINT [ "bash", "-c" ]