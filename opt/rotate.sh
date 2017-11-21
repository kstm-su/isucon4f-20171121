cp /tmp/kataribe.log /tmp/kataribe.log.1
/opt/kataribe -conf /opt/kataribe.toml < /var/log/nginx/access_log > /tmp/kataribe.log
cp /var/log/nginx/access_log /var/log/nginx/access_log.1
>/var/log/nginx/access_log
