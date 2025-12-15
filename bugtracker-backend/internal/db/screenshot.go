package db

import (
	"encoding/json"
	"fmt"
	"time"

	"bugtracker-backend/internal/models"
	"go.etcd.io/bbolt"
)

// CreateScreenshot saves screenshot metadata in BoltDB
func CreateScreenshot(s *models.Screenshot) error {
	return db.Update(func(tx *bbolt.Tx) error {

		counter := tx.Bucket(screenshotCounterBucket)
		last := counter.Get([]byte("lastScreenshotID"))

		nextID := 1
		if last != nil {
			nextID = btoi(last) + 1
		}

		s.ID = nextID
		s.CreatedAt = time.Now()

		data, err := json.Marshal(s)
		if err != nil {
			return err
		}

		if err := tx.Bucket(screenshotsBucket).Put(itob(s.ID), data); err != nil {
			return err
		}

		return counter.Put([]byte("lastScreenshotID"), itob(nextID))
	})
}

// GetScreenshotsByBugID returns all screenshots for one bug
func GetScreenshotsByBugID(bugID int) ([]*models.Screenshot, error) {
	var list []*models.Screenshot

	err := db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(screenshotsBucket).ForEach(func(_, v []byte) error {
			var s models.Screenshot
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			if s.BugID == bugID {
				list = append(list, &s)
			}
			return nil
		})
	})

	return list, err
}

// DeleteScreenshot removes screenshot metadata
func DeleteScreenshot(id int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		if tx.Bucket(screenshotsBucket).Get(itob(id)) == nil {
			return fmt.Errorf("screenshot not found")
		}
		return tx.Bucket(screenshotsBucket).Delete(itob(id))
	})
}
