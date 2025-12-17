package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bugtracker-backend/internal/db"
	"bugtracker-backend/internal/models"
	"github.com/gorilla/mux"
)

func GetBugLinks(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	links, _ := db.GetBugLinks(id)
	json.NewEncoder(w).Encode(links)
}

func CreateBugLink(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var link models.BugLink
	json.NewDecoder(r.Body).Decode(&link)
	link.BugID = id

	db.CreateBugLink(&link)
	json.NewEncoder(w).Encode(link)
}
