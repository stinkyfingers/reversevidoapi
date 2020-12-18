package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stinkyfingers/reversevideoapi/video"
)

type UploadResponse struct {
	Uri string `json:"uri"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("videoFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	err = video.Reverse(file, header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// var buf bytes.Buffer
	// _, err = io.Copy(&buf, file)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer buf.Reset()
	// // TODO Handle video
	// tmp, err := ioutil.TempFile("", "")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// _, err = buf.WriteTo(tmp)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Println(len(buf.Bytes()), header.Filename, header.Size)
	// err = exec.Command("ffmpeg", "-i", tmp.Name(), "-vf", "reverse", "reversed.mp4").Run()
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Println("SENBDING")
	// w.Header().Add("Content-Type", "video/mp4")
	// _, err = io.Copy(w, &buf)
	// http.ServeFile(w, r, tmp.Name())
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	json.NewEncoder(w).Encode(&UploadResponse{Uri: header.Filename})
	fmt.Println("SENT")
}
