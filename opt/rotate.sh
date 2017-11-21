cp /tmp/kataribe.log /tmp/kataribe.log.1
/opt/kataribe -conf /opt/kataribe.toml < /var/log/nginx/access.log > /tmp/kataribe.log
cp /var/log/nginx/access.log /var/log/nginx/access.log.1
>/var/log/nginx/access.log
