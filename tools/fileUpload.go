package tools

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var (
	ErrInvalidFileType = http.ErrNotSupported
	ErrFileTooLarge    = http.ErrContentLength
	ErrSavingFile      = http.ErrBodyNotAllowed
	ErrDatabase        = http.ErrAbortHandler
)

func SaveUploadedFile(file multipart.File, handler *multipart.FileHeader, userID int64) (string, error) {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"video/mp4":  true,
	}
	if !allowedTypes[handler.Header.Get("Content-Type")] {
		return "", ErrInvalidFileType
	}

	const maxFileSize = 10 << 20
	if handler.Size > maxFileSize {
		return "", ErrFileTooLarge
	}

	newFileName := time.Now().Format("20060102150405") + "_" + handler.Filename
	dst, err := os.Create("./uploads/" + newFileName)
	if err != nil {
		return "", ErrSavingFile
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", ErrSavingFile
	}

	return newFileName, nil
}
