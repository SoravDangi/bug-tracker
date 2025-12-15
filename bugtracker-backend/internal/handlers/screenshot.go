package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bugtracker-backend/internal/db"
	"bugtracker-backend/internal/models"
)

func UploadScreenshot(w http.ResponseWriter, r *http.Request) {
	bugID, _ := strconv.Atoi(r.FormValue("bugId"))

	filePath := r.FormValue("filePath") // already saved by middleware

	s := &models.Screenshot{
		BugID:    bugID,
		FilePath: filePath,
	}

	if err := db.CreateScreenshot(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(s)
}

func GetScreenshots(w http.ResponseWriter, r *http.Request) {
	bugID, _ := strconv.Atoi(r.URL.Query().Get("bugId"))

	data, err := db.GetScreenshotsByBugID(bugID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func DeleteScreenshot(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if err := db.DeleteScreenshot(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
