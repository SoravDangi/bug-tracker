package db

import (
	"encoding/json"
	"fmt"

	"bugtracker-backend/internal/models"

	"go.etcd.io/bbolt"
)

func CreateScreenshot(s *models.Screenshot) error {
	return db.Update(func(tx *bbolt.Tx) error {
		sb := tx.Bucket(screenshotsBucket)
		cb := tx.Bucket(screenshotCounterBucket)

		id := btoi(cb.Get([]byte("lastScreenshotID"))) + 1
		if err := cb.Put([]byte("lastScreenshotID"), itob(id)); err != nil {
			return err
		}

		s.ID = id

		data, err := json.Marshal(s)
		if err != nil {
			return err
		}

		return sb.Put(itob(s.ID), data)
	})
}

func GetScreenshotsByBugID(bugID int) ([]models.Screenshot, error) {
	var result []models.Screenshot

	err := db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(screenshotsBucket).ForEach(func(_, v []byte) error {
			var s models.Screenshot
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			if s.BugID == bugID {
				result = append(result, s)
			}
			return nil
		})
	})

	return result, err
}

func DeleteScreenshot(id int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(screenshotsBucket).Delete(itob(id))
	})
}

func DeleteScreenshotsByBugID(bugID int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		sb := tx.Bucket(screenshotsBucket)
		var toDelete [][]byte

		sb.ForEach(func(k, v []byte) error {
			var s models.Screenshot
			json.Unmarshal(v, &s)
			if s.BugID == bugID {
				toDelete = append(toDelete, k)
			}
			return nil
		})

		for _, k := range toDelete {
			sb.Delete(k)
		}
		return nil
	})
}
