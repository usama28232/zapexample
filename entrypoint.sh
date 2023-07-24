#!/bin/bash

cd /app/
./zapexample --i $INTERVAL &

cd /usr/share/filebeat
./filebeat modules enable kibana

./filebeat -e -c "filebeat.yaml"
