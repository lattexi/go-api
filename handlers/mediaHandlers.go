package handlers

import (
	"encoding/json"
	"fmt"
	"go-api/db"
	"go-api/tools"
	"net/http"
)

// GET /mediaitems
type MediaItem struct {
	MediaID     int    `json:"media_id"`
	UserID      int    `json:"user_id"`
	Filename    string `json:"filename"`
	MediaType   string `json:"media_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func GetMediaItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT media_id, user_id, filename, media_type, title, description, created_at FROM mediaitems")
	if err != nil {
		return
	}
	defer rows.Close()

	var mediaItems []MediaItem
	for rows.Next() {
		var item MediaItem
		err := rows.Scan(&item.MediaID, &item.UserID, &item.Filename, &item.MediaType, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			continue
		}
		mediaItems = append(mediaItems, item)
	}

	j, err := json.Marshal(mediaItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

type MediaUploadRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// POST /mediaitems
func CreateMediaItem(w http.ResponseWriter, r *http.Request) {
	var req MediaUploadRequest
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := tools.GetClaims(r)
	userID := claims.UserID

	req.Title = r.FormValue("title")
	req.Description = r.FormValue("description")

	if req.Title == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := tools.SaveUploadedFile(file, handler, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO mediaitems (user_id, filename, title, description, media_type) VALUES (?, ?, ?, ?, ?)", userID, filename, req.Title, req.Description, "image/jpeg")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, `{"media_id": %d}`, id)
}

// GET /files/{filename}
func ServeFile(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	if filename == "" {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, "./uploads/"+filename)
}

// DELETE /mediaitems
type MediaDeleteRequest struct {
	MediaID int `json:"media_id"`
}

func DeleteMediaItem(w http.ResponseWriter, r *http.Request) {
	var req MediaDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := tools.GetClaims(r)
	userID := claims.UserID

	rows, err := db.DB.Exec("DELETE FROM mediaitems WHERE media_id = ? AND user_id = ?", req.MediaID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if affected, err := rows.RowsAffected(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if affected == 0 {
		http.Error(w, "media item not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PUT /mediaitems
type MediaUpdateRequest struct {
	MediaID     int    `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func UpdateMediaItem(w http.ResponseWriter, r *http.Request) {
	var req MediaUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	claims := tools.GetClaims(r)
	userID := claims.UserID

	rows, err := db.DB.Exec("UPDATE mediaitems SET title = ?, description = ? WHERE media_id = ? AND user_id = ?", req.Title, req.Description, req.MediaID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if affected, err := rows.RowsAffected(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if affected == 0 {
		http.Error(w, "media item not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
