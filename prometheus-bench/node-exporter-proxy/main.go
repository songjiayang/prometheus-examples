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
)

func main() {
	nodeExportEndpoint = os.Getenv("NODE_EXPORTER_ENDPOTINT")
	if nodeExportEndpoint == "" {
		nodeExportEndpoint = "node_exporter"
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		rw.RLock()
		for k, v := range *cacheHeader {
			w.Header().Add(k, v[0])
		}
		w.Write(cacheMetrics)
	})

	go http.ListenAndServe(":8080", http.DefaultServeMux)

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go syncNodeExporerData(sigs, done)

	<-done
}

func syncNodeExporerData(sigs chan os.Signal, done chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-sigs:
			log.Println("ticker stoping...")
			done <- struct{}{}
			time.Sleep(time.Second)
		case <-ticker.C:
			fetchAndUpdateData()
		}
	}
}

func fetchAndUpdateData() {
	http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:9090/metrics", nodeExportEndpoint), nil)
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
