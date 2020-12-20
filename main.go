package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/stinkyfingers/reversevideoapi/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, mux()))
}

func mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handlers.UploadHandler)
	mux.HandleFunc("/download", handlers.DownloadHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Print("health called")
		_, err := w.Write([]byte("OK"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		foo := struct {
			Foo string `json:"foo"`
		}{
			"bar",
		}
		err := json.NewEncoder(w).Encode(&foo)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	return mux
}
