package handlers

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
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
	bugID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid bug id", http.StatusBadRequest)
		return
	}

	// ✅ REQUIRED: parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ✅ Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		http.Error(w, "failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf(
		"uploads/%d_%d%s",
		bugID,
		time.Now().UnixNano(),
		filepath.Ext(header.Filename),
	)

	dst, err := os.Create(filename)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}

	screenshot := &models.Screenshot{
		BugID:    bugID,
		FilePath: filename,
	}

	if err := db.CreateScreenshot(screenshot); err != nil {
		http.Error(w, "failed to save screenshot record", http.StatusInternalServerError)
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
