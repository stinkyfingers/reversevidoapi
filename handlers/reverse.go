package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/stinkyfingers/reversevideoapi/video"
)

type UploadResponse struct {
	Uri string `json:"uri"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("videoFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	id, err := video.Reverse(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&UploadResponse{Uri: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	body, err := video.GetVideo(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "video/mp4")
	_, err = io.Copy(w, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("SENT VIDEO")
}
