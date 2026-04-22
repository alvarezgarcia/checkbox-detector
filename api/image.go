package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
)

func isImage(file multipart.File) bool {
	_, _, err := image.DecodeConfig(file)
	file.Seek(0, io.SeekStart)
	return err == nil
}
