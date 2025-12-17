package models

type BugLink struct {
	ID    int    `json:"id"`
	BugID int    `json:"bug_id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
