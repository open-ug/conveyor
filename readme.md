# CRANE (THE CRANOM PLATFORM ENGINE)

## Introduction

Crane is the core engine of the Cranom platform. It is a lightweight, high-performance, and scalable platform that is designed to be used in a variety of applications. 

Learn more about the Crane Engine at [https://cranom.tech/platform-tools/crane-engine](https://cranom.tech/platform-tools/crane-engine).

## Steps


```bash
docker volume create prometheus-data
# Start Prometheus container
docker run \
    -p 9090:9090 \
    -v /home/junior/dev/cranom/crane/yaml/prometheus.yml:/etc/prometheus/prometheus.yml \
    -v prometheus-data:/prometheus \
    --detach=true   --name=promserver \
    prom/prometheus

```