package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func detectHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if !assert(err == nil, w, "missing image field", http.StatusBadRequest) {
		return
	}
	defer file.Close()

	if !assert(isImage(file), w, "file is not an image", http.StatusUnprocessableEntity) {
		return
	}

	debug := r.URL.Query().Get("debug") == "true"
	result, err := runDetect(file, header.Filename, debug)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getDetectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := readJobResult(id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getDetectImageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	img, contentType, err := readJobImage(id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))
	w.Write(img)
}

func main() {
	godotenv.Load()

	if bucket := os.Getenv("BUCKET_DIR"); bucket != "" {
		if _, err := os.Stat(bucket); os.IsNotExist(err) {
			log.Fatalf("BUCKET_DIR %q does not exist", bucket)
		}
	}

	http.HandleFunc("POST /detect", detectHandler)
	http.HandleFunc("GET /detect/{id}", getDetectHandler)
	http.HandleFunc("GET /detect/{id}/image", getDetectImageHandler)
	http.ListenAndServe(":8080", nil)
}
