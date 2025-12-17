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
	api := r.PathPrefix("/api").Subrouter()

	// Bugs
	api.HandleFunc("/bugs", CreateBug).Methods("POST")
	api.HandleFunc("/bugs", GetBugs).Methods("GET")
	api.HandleFunc("/bugs", DeleteAllBugs).Methods("DELETE")
	api.HandleFunc("/bugs/{id}", GetBug).Methods("GET")
	api.HandleFunc("/bugs/{id}", UpdateBug).Methods("PUT")
	api.HandleFunc("/bugs/{id}", DeleteBug).Methods("DELETE")

	// Screenshots
	api.HandleFunc("/bugs/{id}/screenshots", UploadScreenshot).Methods("POST")
	api.HandleFunc("/bugs/{id}/screenshots", GetScreenshots).Methods("GET")
	api.HandleFunc("/screenshots/{id}", DeleteScreenshot).Methods("DELETE")

	// ✅ Bug links (THIS FIXES 404)
	api.HandleFunc("/bugs/{id}/links", GetBugLinks).Methods("GET")
	api.HandleFunc("/bugs/{id}/links", CreateBugLink).Methods("POST")
	api.HandleFunc("/bugs/links/{id}", DeleteBugLink).Methods("DELETE")

	// Comments
	RegisterCommentRoutes(api)
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
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid bug ID",
		})
		return
	}

	// 1️⃣ Fetch screenshots
	screenshots, err := db.GetScreenshotsByBugID(idInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to fetch screenshots",
		})
		return
	}

	// 2️⃣ Delete screenshot files
	for _, s := range screenshots {
		_ = os.Remove(s.FilePath)
	}

	// 3️⃣ Delete screenshot records
	if err := db.DeleteScreenshotsByBugID(idInt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to delete screenshots",
		})
		return
	}

	// 4️⃣ Delete comments
	_ = db.DeleteCommentsByBugID(idInt)

	// 5️⃣ Delete bug links
	_ = db.DeleteBugLinksByBugID(idInt)

	// 6️⃣ Delete bug
	if err := db.DeleteBug(idInt); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "bug not found",
		})
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
