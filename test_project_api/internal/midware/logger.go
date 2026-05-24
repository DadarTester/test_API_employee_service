package middleware

import (
    "net/http"
    "time"
    log "github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        lw := &loggingWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(lw, r)
        log.WithFields(log.Fields{
            "method":   r.Method,
            "uri":      r.RequestURI,
            "status":   lw.statusCode,
            "duration": time.Since(start),
        }).Info("request")
    })
}

type loggingWriter struct {
    http.ResponseWriter
    statusCode int
}

func (lw *loggingWriter) WriteHeader(code int) {
    lw.statusCode = code
    lw.ResponseWriter.WriteHeader(code)
}