package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

// Logger is a logging middleware. Path, request time and Method are recorded
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Header.Get("X-Forwarded-For"),
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func main() {
	var (
		dir = flag.String("path", "public", "Path to static site")
	)

	fs := http.FileServer(http.Dir(*dir))
	fs = Logger(fs)
	http.Handle("/", fs)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
