---
sidebar_position: 4
---

# Monitoring Conveyor CI

Conveyor CI exports metrics under the `/metrics` path. It uses Prometheus style metrics.

These can be fetched via `curl` or any http client

```sh
$ curl -L http://localhost:8080/metrics

# HELP active_connections Number of active connections
# TYPE active_connections gauge
active_connections 0
# HELP cpu_usage_percent CPU usage percentage
# TYPE cpu_usage_percent gauge
cpu_usage_percent 5.352644836133662
http_request_duration_seconds_bucket{endpoint="/",method="GET",le="0.005"} 1
# TYPE http_requests_in_flight gauge
http_requests_in_flight 1
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{endpoint="/metrics",method="GET",status_code="200"} 1
http_requests_total{endpoint="/swagger/*",method="GET",status_code="200"} 8
# HELP memory_usage_bytes Memory usage in bytes
# TYPE memory_usage_bytes gauge
memory_usage_bytes{type="heap_alloc"} 1.3269896e+07
memory_usage_bytes{type="heap_idle"} 9.568256e+06
```