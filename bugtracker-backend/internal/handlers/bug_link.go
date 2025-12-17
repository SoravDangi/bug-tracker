package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"bugtracker-backend/internal/db"
	"bugtracker-backend/internal/models"
)

func GetBugLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	bugID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid bug id", http.StatusBadRequest)
		return
	}

	links, err := db.GetBugLinksByBugID(bugID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(links)
}

func CreateBugLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	bugID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid bug id", http.StatusBadRequest)
		return
	}

	var body struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	link := &models.BugLink{
		ID:        int(time.Now().UnixNano()), // simple unique ID
		BugID:     bugID,
		Title:     body.Title,
		URL:       body.URL,
		
	}

	if err := db.CreateBugLink(link); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

func DeleteBugLink(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	linkID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid link id", http.StatusBadRequest)
		return
	}

	if err := db.DeleteBugLink(linkID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
