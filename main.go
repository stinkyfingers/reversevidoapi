package main

import (
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
	mux.HandleFunc("/check", handlers.CheckVideoStatus)
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/ffmpeg", handlers.Ffmpeg)
	return mux
}
