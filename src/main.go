package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type interceptingWriter struct {
	code int
	http.ResponseWriter
}

func (iw *interceptingWriter) WriteHeader(code int) {
	iw.code = code
	iw.ResponseWriter.WriteHeader(code)
}

// Log is the structure for logging
type Log struct {
	Time     time.Time `json:"time"`
	IP       string    `json:"ip"`
	Method   string    `json:"method"`
	Status   int       `json:"status"`
	Duration string    `json:"duration"`
	Path     string    `json:"path"`
	Referrer string    `json:"referrer"`
}

// Logger is a logging middleware. Path, request time and Method are recorded
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		iw := interceptingWriter{200, w}

		inner.ServeHTTP(&iw, r)

		d := fmt.Sprintf("%s", time.Since(start))
		l := Log{
			Time:     time.Now(),
			IP:       r.Header.Get("X-Forwarded-For"),
			Method:   r.Method,
			Status:   iw.code,
			Duration: d,
			Path:     r.RequestURI,
			Referrer: r.Referer(),
		}

		b, _ := json.Marshal(l)
		fmt.Println(string(b))
	})
}

// Cache is a caching middleware
func Cache(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner.ServeHTTP(w, r)

		w.Header().Add("Cache-Control", "public")
		if strings.Compare(w.Header().Get("Content-Type"), "text/html") != 0 {
			w.Header().Add("Cache-Control", "maxage=604800") // 1 Week
		} else {
			w.Header().Add("Cache-Control", "maxage=86400") // 1 Day
		}
	})
}

func main() {
	var (
		dir = flag.String("path", "public", "Path to static site")
	)

	fs := http.FileServer(http.Dir(*dir))
	fs = Logger(fs)
	fs = Cache(fs)
	http.Handle("/", fs)
	http.ListenAndServe(":3000", nil)
}
