package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"os"

	"github.com/gorilla/mux"

	"bugtracker-backend/internal/db"
	"bugtracker-backend/internal/models"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/bugs", CreateBug).Methods("POST")
	r.HandleFunc("/bugs", GetBugs).Methods("GET")
	r.HandleFunc("/bugs", DeleteAllBugs).Methods("DELETE")
	r.HandleFunc("/bugs/{id}", GetBug).Methods("GET")
	r.HandleFunc("/bugs/{id}", UpdateBug).Methods("PUT")
	r.HandleFunc("/bugs/{id}", DeleteBug).Methods("DELETE")
	r.HandleFunc("/bugs/{id}/screenshots", UploadScreenshot).Methods("POST")
    r.HandleFunc("/bugs/{id}/screenshots", GetScreenshots).Methods("GET")
    r.HandleFunc("/screenshots/{id}", DeleteScreenshot).Methods("DELETE")

	RegisterCommentRoutes(r)
}

func CreateBug(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateBug called from %s", r.RemoteAddr)
	log.Printf("Request headers: %v", r.Header)
	log.Printf("Request method: %s", r.Method)

	var req models.CreateBugRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode create bug request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "title is required",
		})
		return
	}


	bug := &models.Bug{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 🔥 This MUST assign an ID, or GET /bugs/:id WILL FAIL
	if err := db.CreateBug(bug); err != nil {
		log.Printf("Failed to create bug: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(bug)
}


func GetBugs(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetBugs called from %s", r.RemoteAddr)
	log.Printf("Request headers: %v", r.Header)
	bugs, err := db.GetAllBugs()
	if err != nil {
		log.Printf("Error getting bugs: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully retrieved %d bugs", len(bugs))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bugs)
}

func GetBug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid bug ID",
		})
		return
	}

	bug, err := db.GetBug(idInt)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bug not found" {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(bug)
}

func UpdateBug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid bug ID",
		})
		return
	}

	var req models.CreateBugRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	existingBug, err := db.GetBug(idInt)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bug not found" {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	existingBug.Title = req.Title
	existingBug.Description = req.Description
	existingBug.Status = req.Status
	existingBug.Priority = req.Priority
	existingBug.UpdatedAt = time.Now()

	if err := db.UpdateBug(existingBug); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(existingBug)
}

func DeleteBug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid bug ID", http.StatusBadRequest)
		return
	}

	// 1. Get screenshots linked to this bug
	screenshots, err := db.GetScreenshotsByBugID(idInt)
	if err != nil {
		http.Error(w, "failed to fetch screenshots", http.StatusInternalServerError)
		return
	}

	// 2. Delete screenshot files from disk
	for _, s := range screenshots {
		_ = os.Remove(s.FilePath) // ignore error if file does not exist
	}

	// 3. Delete screenshot records from DB
	if err := db.DeleteScreenshotsByBugID(idInt); err != nil {
		http.Error(w, "failed to delete screenshots", http.StatusInternalServerError)
		return
	}

	// 4. Delete the bug itself
	if err := db.DeleteBug(idInt); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}



func DeleteAllBugs(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteAllBugs called from %s", r.RemoteAddr)
	
	count, err := db.DeleteAllBugs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{
		"deleted": count,
	})
}
