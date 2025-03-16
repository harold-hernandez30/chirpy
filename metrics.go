package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)



func (cfg *apiConfig) handlePrintMetrics(res http.ResponseWriter, req *http.Request) {
	header := res.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	content := fmt.Sprintf(`
	<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	res.Write([]byte(content))
}

func (cfg *apiConfig) handleResetMetrics(res http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
}