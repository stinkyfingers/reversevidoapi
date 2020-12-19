package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/stinkyfingers/reversevideoapi/handlers"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", mux()))
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

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DOWNLOAD")
	f, err := os.Open("./runlocal/grin.mov")
	if err != nil {
		fmt.Println("E", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "video/mp4")
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("SENT SAMPLE")
}
