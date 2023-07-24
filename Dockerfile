FROM docker.elastic.co/beats/filebeat:7.15.2

WORKDIR /app

COPY filebeat.yaml /usr/share/filebeat/filebeat.yaml
COPY "zapexample" /app/
COPY entrypoint.sh /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]