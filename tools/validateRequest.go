package tools

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
)

func ValidateJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	const maxSize int64 = 1 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return false
	}
	return true
}

func ValidateMultipartForm(w http.ResponseWriter, r *http.Request, requiredFields ...string) bool {
	const maxSize int64 = 10 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return false
	}
	defer r.MultipartForm.RemoveAll()

	for _, field := range requiredFields {
		if r.FormValue(field) == "" {
			http.Error(w, "missing field: "+field, http.StatusBadRequest)
			return false
		}
	}

	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return false
	}
	defer file.Close()

	const maxFileSize = 5 << 20
	if fh.Size <= 0 {
		http.Error(w, "empty file", http.StatusBadRequest)
		return false
	}
	if fh.Size > maxFileSize {
		http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
		return false
	}

	allowedTypes := []string{"image/jpeg", "image/png", "video/mp4"}

	mimeType, err := DetectFileType(file)
	if err != nil {
		http.Error(w, "could not detect file type", http.StatusBadRequest)
		return false
	}

	ok := slices.Contains(allowedTypes, mimeType)
	if !ok {
		http.Error(w, "unsupported file type", http.StatusUnsupportedMediaType)
		return false
	}

	return true
}

func DetectFileType(file io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	if n == 0 {
		return "", io.ErrUnexpectedEOF
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}
