# zapexample
Golang Zap Logging example with `Elastic Search`, `Kibana` &amp; `Filebeat` for centralised logging.


Recently I switched from Java to Golang, and I am finding it easier to learn and fun than other languages I have worked on. It does look a bit complex from a distance, but trust me as soon as you get close, it is really very easy.


## Problem

This is a pretty common question that comes in our minds whether we are looking to implement microservices or loadbalancing in our applications. You can either keep separate files and manage/analyze/debug logs manually or write another service to handle that and collect them at one place, both approaches have their pros and cons.


## Solution

To overcome this problem, I recently presented a 3rd Party System which is capable of providing a perfect solution.
Moreover, it offers stuff even better than that

I am using `Elastic Stack` which goes as:

* [Elasticsearch](https://www.elastic.co/elasticsearch/)
* [Kibana](https://www.elastic.co/kibana/)
* [Filebeat](https://www.elastic.co/beats/filebeat)

### Elasticsearch

`Elasticsearch` is a distributed, RESTful search and analytics engine capable of addressing a growing number of use cases. As the heart of the `Elastic Stack`, it centrally stores your data for lightning fast search, fine‑tuned relevancy, and powerful analytics that scale with ease.


### Kibana

`Kibana` is a powerful analysis on any data from any source, from threat intelligence to search analytics, logs to application monitoring, and much more.


### Filebeat

`Filebeat` helps you keep the simple things simple by offering a lightweight way to forward and centralize logs and files.


## BEWARE

This approach uses ["Uber-Zap Logger"](https://github.com/uber-go/zap) for logging which is Blazing fast, structured, leveled logging in Go.

From that perspective this is pretty basic and still under construction but looking at the code you may get the idea that we are trying to achieve level-based logging as `Log4j` offers and changing logging as per incoming request so we may not end up with enormous log files in **Debug** mode.

And I am running on following versions:

* `Go` v1.20.5
* `Uber-Zap` v1.24.0
* `Elasticsearch` v8.8.2
* `Kibana` v8.8.2
* `Filebeat` v7.15.2

For the sake of simplicity, I have disabled security in **elasticsearch** `../elasticsearch-8.8.2/config/elasticsearch.yml`


## Working

This application is designed to write greetings in log file at given intervals which are collected by `Filebeat`

I am using my local machine as server having `Elasticsearch` & `Kibana` and docker images replication deployment environment having application executable and `Filebeat`


Standalone Application takes **-i** flag to set interval value and can be executed by:

```
./zapexample.exe -i 5
```

Or for Linux environment, use the `linux_executable.bat` file to create linux executable


Running application will produce output like following in the same directory in `log.txt` file:


```
[24/07/2023 07:38:55] info Zap Package level logging example {"i": 5}
[24/07/2023 07:38:55] debug Got a greeting request
[24/07/2023 07:38:55] debug Processing Greeting with interval {"i": 5}
[24/07/2023 07:39:00] info Greetings brother
[24/07/2023 07:39:05] info Greetings brother
[24/07/2023 07:39:10] info Greetings brother
[24/07/2023 07:39:15] info Greetings brother
[24/07/2023 07:39:20] info Greetings brother
```
*Application will keep writing Greetings Brother with provided interval until terminated*

## Docker Configuration

To replicate multiple instances, I have created docker container which spins off from `filebeat` base image, copies `filebeat.yaml` to `filebeat` installation directory (which holds `filebeat` configuration to Run in Docker), copies application executable and copies entry point script (to run at container startup)

```
FROM docker.elastic.co/beats/filebeat:7.15.2

WORKDIR /app

COPY filebeat.yaml /usr/share/filebeat/filebeat.yaml
COPY "zapexample" /app/
COPY entrypoint.sh /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
```

## Filebeat Configuration

Remember, now we are looking from application/client perspective!

since our executable is onside `/app` directory we have to specify log file path. Also, we can specify other fields like `tags`, `host.name` etc which can be used for filtering and analyzing

Other configuration are directing `filebeat` to `elasticsearch` and `kibana` instances


```
filebeat.inputs:
- type: log
paths:
- /app/log.txt
fields_under_root: true
fields:
host.name: "zapexample-${SERVER}"
tags: ["${SERVER}"]

output.elasticsearch:
hosts: ["http://${ELASTICSEARCH_HOST}:9200"]

setup.kibana:
host: "http://${KIBANA_HOST}:5601"
```

### Entrypoint Script

This script executes after container starts, which executes our application in background (since we do not need logs on console), changes directory to `filebeat` installation directory, enables `kibana` support on `filebeat` and runs `filebeat` with specified configuration file by `... -c "filebeat.yaml"`
**-e** flag writes output to the console, you can skip that if you do not want it


```
#!/bin/bash

cd /app/
./zapexample --i $INTERVAL &

cd /usr/share/filebeat
./filebeat modules enable kibana

./filebeat -e -c "filebeat.yaml"
```


### Execution

Create docker image by:

`docker build -t app1_10 .`

To run this docker image:

`docker run -d -e ELASTICSEARCH_HOST=host.docker.internal -e KIBANA_HOST=host.docker.internal -e SERVER="APP1_10" -e INTERVAL=10 app1_10`

`${SERVER}`, `${ELASTICSEARCH_HOST}`, `${KIBANA_HOST}` & `$INTERVAL` are environment variables to make application modifiable from command-line

`host.docker.internal` is used to refer host machine's address from docker image
`$INTERVAL` flag would modify application logging interval

* Run `elasticsearch` from installation directory
* Run `kibana` installation directory
* Run docker images

By default, `Kibana` can be accessed at **http://localhost:5601/**
Go to **Observability > Logs** and start viewing logs submitted by `filebeat`

### Feel free to edit/expand/explore this repository

For feedback and queries, reach me on LinkedIn at [here](https://www.linkedin.com/in/usama28232/?original_referer=)