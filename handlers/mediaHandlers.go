package handlers

import (
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
	rows, err := db.DB.Query("SELECT media_id, user_id, filename, media_type, title, description, created_at FROM mediaitems ORDER BY created_at DESC")
	if err != nil {
		return
	}
	defer rows.Close()

	var response []MediaItem
	for rows.Next() {
		var item MediaItem
		err := rows.Scan(&item.MediaID, &item.UserID, &item.Filename, &item.MediaType, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			continue
		}
		response = append(response, item)
	}

	tools.JSONResponse(w, http.StatusOK, response)
}

// GET /mediaitems?title={title}
func GetMediaItemsByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "missing title query parameter", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(
		"SELECT media_id, user_id, filename, media_type, title, description, created_at FROM mediaitems WHERE title LIKE ? ORDER BY created_at DESC",
		"%"+title+"%",
	)
	if err != nil {
		return
	}
	defer rows.Close()

	var response []MediaItem
	for rows.Next() {
		var item MediaItem
		err := rows.Scan(&item.MediaID, &item.UserID, &item.Filename, &item.MediaType, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			continue
		}
		response = append(response, item)
	}

	tools.JSONResponse(w, http.StatusOK, response)
}

// POST /mediaitems
type MediaCreateResponse struct {
	MediaID int64 `json:"media_id"`
}

func CreateMediaItem(w http.ResponseWriter, r *http.Request) {
	if !tools.ValidateMultipartForm(w, r, "title") {
		return
	}

	claims, ok := tools.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := tools.SaveUploadedFile(file, fh, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	mimetype, _ := tools.DetectFileType(file)

	result, err := db.DB.Exec("INSERT INTO mediaitems (user_id, filename, title, description, media_type) VALUES (?, ?, ?, ?, ?)", userID, filename, title, description, mimetype)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := MediaCreateResponse{MediaID: id}
	tools.JSONResponse(w, http.StatusCreated, response)
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
	tools.ValidateJSON(w, r, &req)

	claims, ok := tools.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID

	rows, err := db.DB.ExecContext(r.Context(), "DELETE FROM mediaitems WHERE media_id = ? AND user_id = ?", req.MediaID, userID)
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
	tools.ValidateJSON(w, r, &req)

	claims, ok := tools.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	rows, err := db.DB.Exec("UPDATE mediaitems SET title = ?, description = ? WHERE media_id = ? AND user_id = ?", req.Title, req.Description, req.MediaID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, err := rows.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n == 0 {
		http.Error(w, "media item not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
