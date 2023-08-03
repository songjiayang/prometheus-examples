package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	provider := initMeter()
	defer provider.Shutdown(context.Background())

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go timeDuration(provider, sigs, done)
	<-done
}

func initMeter() *metric.MeterProvider {
	ep, _ := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint("prometheus:9090"),
		otlpmetrichttp.WithURLPath("/api/v1/otlp/v1/metrics"),
		otlpmetrichttp.WithInsecure(),
	)

	ep2, _ := stdoutmetric.New()

	ctx := context.Background()
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("service1"),
			semconv.ServiceNamespace("staging"),
			semconv.ServiceVersion("v0.1"),
			semconv.ServiceInstanceID("instance1"),
		),
	)

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(ep, metric.WithInterval(15*time.Second)),
		),
		metric.WithReader(
			metric.NewPeriodicReader(ep2, metric.WithInterval(15*time.Second)),
		),
	)

	return provider
}

func timeDuration(provider *metric.MeterProvider, sigs chan os.Signal, done chan struct{}) {
	meter := provider.Meter("http",
		api.WithInstrumentationVersion("v0.0.1"),
	)

	httpDurationsHistogram, _ := meter.Float64Histogram(
		"http_durations_histogram_seconds",
		api.WithDescription("Http latency distributions."),
	)

	opt := api.WithAttributes(
		attribute.Key("method").String(http.MethodGet),
		attribute.Key("status").String("200"),
	)

	ctx := context.Background()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			elapsed := float64(rand.Intn(100))
			httpDurationsHistogram.Record(
				ctx,
				elapsed,
				opt,
			)
		case <-sigs:
			log.Println("ticker stoping...")
			done <- struct{}{}
			time.Sleep(time.Second)
		}
	}
}
