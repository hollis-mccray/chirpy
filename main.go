package main

import (
	"log"
	"sync/atomic"
	"net/http"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main()  {
	const filepathRoot = "."
	const port = ":8080"
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz",handlerReady)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	mux.HandleFunc("GET /admin/metrics",apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset",apiCfg.handlerReset)
	server := http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
	log.Printf("In the end...")
}