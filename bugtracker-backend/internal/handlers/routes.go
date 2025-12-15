package handlers

import "github.com/gorilla/mux"

func RegisterRoutes(r *mux.Router) {

	// Bugs
	r.HandleFunc("/bugs", CreateBug).Methods("POST")
	r.HandleFunc("/bugs", GetAllBugs).Methods("GET")
	r.HandleFunc("/bugs/{id}", GetBug).Methods("GET")
	r.HandleFunc("/bugs/{id}", UpdateBug).Methods("PUT")
	r.HandleFunc("/bugs/{id}", DeleteBug).Methods("DELETE")

	// Comments
	r.HandleFunc("/bugs/{id}/comments", CreateComment).Methods("POST")
	r.HandleFunc("/bugs/{id}/comments", GetComments).Methods("GET")

	// 🔴 SCREENSHOTS (THIS IS WHAT YOU ARE MISSING)
	r.HandleFunc("/bugs/{id}/screenshots", UploadScreenshot).Methods("POST")
	r.HandleFunc("/bugs/{id}/screenshots", GetScreenshots).Methods("GET")
	r.HandleFunc("/screenshots/{id}", DeleteScreenshot).Methods("DELETE")
}
