package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Server hits reset"))
}