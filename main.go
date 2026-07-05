package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const (
		port         = ":8080"
		filepathRoot = "."
	)

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	apiCfg := &apiConfig{}
	fsServer := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fsServer)))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /api/metrics", func(cfg *apiConfig) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hits: " + strconv.Itoa(int(cfg.fileserverHits.Load()))))
		}
	}(apiCfg))

	mux.HandleFunc("POST /api/reset", func(cfg *apiConfig) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Store(0)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Reset OK"))
		}
	}(apiCfg))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
