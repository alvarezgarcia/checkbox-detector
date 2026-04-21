package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func isImage(file multipart.File) bool {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return false
	}
	file.Seek(0, io.SeekStart)
	mime := http.DetectContentType(buf)
	return strings.HasPrefix(mime, "image/")
}
