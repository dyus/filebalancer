package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Debug().Msg(fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, time.Since(start)))
	})
}
