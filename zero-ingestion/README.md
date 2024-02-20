# Prometheus created-timestamp-zero-ingestion example 

An example about prometheus `created-timestamp-zero-ingestion`.


# How to usage

## step 1: start app with make.

```
make build && make start
```

## step 2: visit the prometheus web console.

- enable created-timestamp-zero-ingestion then click [url](http://localhost:9090/graph?g0.expr=rate(http_requests_total%7Bcode%3D%22500%22%7D%5B1m%5D)&g0.tab=1&g0.display_mode=lines&g0.show_exemplars=0&g0.range_input=1h) you will see:

```
{code="500", instance="zero-ingestion:8080", job="example-app"} 0.21265284423179157
```

- disable created-timestamp-zero-ingestion then click [url](http://localhost:9091/graph?g0.expr=rate(http_requests_total%7Bcode%3D%22500%22%7D%5B1m%5D)&g0.tab=1&g0.display_mode=lines&g0.show_exemplars=0&g0.range_input=1h) you will see:

```
{code="500", instance="zero-ingestion:8080", job="example-app"} 0
```

