package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

func main() {
	filepathRoot := "."
	const port = "8080"
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	mux.HandleFunc("/healthz", handlerHealthz)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerHealthz(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader((http.StatusOK))
	hits := cfg.fileserverHits.Load()
	rw.Write([]byte("Hits: " + strconv.Itoa(int(hits))))
}

func (cfg *apiConfig) handlerReset(rw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	hits := cfg.fileserverHits.Load()
	rw.WriteHeader((http.StatusOK))
	rw.Write([]byte("Hits: " + strconv.Itoa(int(hits))))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})
}
