package tools

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func SaveUploadedFile(file multipart.File, handler *multipart.FileHeader, userID int64) (string, error) {
	newFileName := time.Now().Format("20060102150405") + "_" + handler.Filename
	dst, err := os.Create("./uploads/" + newFileName)
	if err != nil {
		return "", http.ErrBodyNotAllowed
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", http.ErrBodyNotAllowed
	}

	return newFileName, nil
}
