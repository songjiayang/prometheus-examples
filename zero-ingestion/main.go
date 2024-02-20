package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	requestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total amount of HTTP requests",
	}, []string{"code"})
	reg := prometheus.NewRegistry()
	reg.MustRegister(requestsTotal)

	go func() {
		// 模拟程序运行一段时间后 500 徒增到 10，然后保持稳定不变
		time.Sleep(10 * time.Second)
		requestsTotal.WithLabelValues("500").Add(10)
	}()

	go func() {
		for {
			requestsTotal.WithLabelValues("200").Inc()
			time.Sleep(time.Second)
		}
	}()

	http.Handle(
		"/metrics", promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{}),
	)

	// To test: curl -H 'Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited;q=0.8' localhost:8080/metrics
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
