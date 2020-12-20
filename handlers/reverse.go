package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/google/uuid"
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

	id := fmt.Sprintf("%s.mov", uuid.New().String())
	go func() {
		err = video.Reverse(file, id)
		if err != nil {
			video.UpdateLog(id, false, err.Error())
			return
		}
	}()
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
	defer body.Close()
	w.Header().Add("Content-Type", "video/mp4")
	_, err = io.Copy(w, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print("SENT")
	go video.Cleanup(key)
}

func CheckVideoStatus(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	l, err := video.CheckLog(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Ffmpeg(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("ffmpeg", "-version").CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	ffmpeg := struct {
		Version string `json:"version"`
	}{
		string(out),
	}
	err = json.NewEncoder(w).Encode(&ffmpeg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
