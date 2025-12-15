package models

import "time"

// Screenshot represents a screenshot uploaded for a bug
type Screenshot struct {
	ID        int       `json:"id"`
	BugID     int       `json:"bug_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}
