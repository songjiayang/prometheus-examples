package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	nodeExportEndpoint string
	cacheHeader        *http.Header
	cacheMetrics       []byte
	rw                 sync.RWMutex
	client             *http.Client
)

func main() {
	nodeExportEndpoint = os.Getenv("NODE_EXPORTER_ENDPORINT")
	if nodeExportEndpoint == "" {
		nodeExportEndpoint = "localhost"
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		rw.RLock()
		defer rw.RUnlock()

		for k, v := range *cacheHeader {
			w.Header().Add(k, v[0])
		}
		w.Write(cacheMetrics)
	})

	go http.ListenAndServe(":9100", http.DefaultServeMux)

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go syncNodeExporerMetrics(sigs, done)

	<-done
}

func syncNodeExporerMetrics(sigs chan os.Signal, done chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-sigs:
			log.Println("ticker stoping...")
			done <- struct{}{}
		case <-ticker.C:
			fetchAndUpdateData()
		}
	}
}

func fetchAndUpdateData() {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:9100/metrics", nodeExportEndpoint), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("fetch node exporter metrics with error: %v \n", err)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	header := resp.Header.Clone()

	rw.Lock()
	defer rw.Unlock()

	cacheHeader = &header
	cacheMetrics = data
}
