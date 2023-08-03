# Prometheus otlp example 

An example about prometheus otlp-write-receiver.


# How to usage

## step 1: start app with make.

```
make build && make start
```

## step 2: visit the prometheus web console.

Click [url](http://localhost:9090/graph?g0.expr=http_durations_histogram_seconds_bucket&g0.tab=1&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h) you will see `http_durations_histogram_seconds_bucket` metric values.

