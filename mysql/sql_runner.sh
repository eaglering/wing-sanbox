#!/bin/sh

service mysql start >/dev/null 2>&1
mysql  mysql< /usr/local/dist/create_user.sql -u'root' 
mysql  ri_db < $1 -u'test' -p'test123' 2>&1 | grep -v "Warning:"
mysql  mysql< /usr/local/dist/destroy_user.sql -u'root'
