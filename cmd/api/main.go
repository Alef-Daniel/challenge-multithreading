package main

import (
	http2 "net/http"

	"github.com/alef-daniel/challenge-multithreading/internal/adapters/http"
)

func main() {
	r := http.NewRouter(nil)
	http2.ListenAndServe(":8080", r)
}
