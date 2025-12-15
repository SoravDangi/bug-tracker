package handlers

import (
	"encoding/json"
	"fmt"

	"bugtracker-backend/internal/models"
	"go.etcd.io/bbolt"
)

var screenshotsBucket = []byte("screenshots")

func GetScreenshotsByBugID(bugID int) ([]models.Screenshot, error) {
	var screenshots []models.Screenshot

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(screenshotsBucket)
		if b == nil {
			return nil
		}

		return b.ForEach(func(_, v []byte) error {
			var s models.Screenshot
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			if s.BugID == bugID {
				screenshots = append(screenshots, s)
			}
			return nil
		})
	})

	return screenshots, err
}

func DeleteScreenshotsByBugID(bugID int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(screenshotsBucket)
		if b == nil {
			return nil
		}

		var keysToDelete [][]byte

		b.ForEach(func(k, v []byte) error {
			var s models.Screenshot
			if json.Unmarshal(v, &s) == nil && s.BugID == bugID {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		for _, k := range keysToDelete {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}
