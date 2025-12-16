package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"bugtracker-backend/internal/db"
	"bugtracker-backend/internal/models"
)

func UploadScreenshot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	bugIDStr := vars["id"]

	bugID, err := strconv.Atoi(bugIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid bug ID",
		})
		return
	}

	// Ensure bug exists
	if _, err := db.GetBug(bugID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "bug not found",
		})
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to parse form",
		})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "file is required",
		})
		return
	}
	defer file.Close()

	// Create uploads dir if not exists
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + handler.Filename
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	screenshot := &models.Screenshot{
		BugID:    bugID,
		FilePath: filePath,
	}

	if err := db.CreateScreenshot(screenshot); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(screenshot)
}

func GetScreenshots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bugID, _ := strconv.Atoi(vars["id"])

	screenshots, err := db.GetScreenshotsByBugID(bugID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(screenshots)
}

func DeleteScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	s, err := db.GetScreenshot(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = os.Remove(s.FilePath)
	_ = db.DeleteScreenshot(id)

	w.WriteHeader(http.StatusNoContent)
}
