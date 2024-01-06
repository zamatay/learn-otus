package internalhttp

import (
	"log"
	"net/http"
	"time"
)

func logMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer log.Printf("address:%s; method:%s; uri:%s; time: %v", req.RemoteAddr, req.Method, req.RequestURI, time.Since(time.Now()))
		next.ServeHTTP(w, req)
	})
}
