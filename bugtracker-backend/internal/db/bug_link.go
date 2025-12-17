package db

import (
	"encoding/json"
	"fmt"

	"bugtracker-backend/internal/models"
	"go.etcd.io/bbolt"
)

func GetBugLinksByBugID(bugID int) ([]*models.BugLink, error) {
	var links []*models.BugLink

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugLinksBucket)
		if b == nil {
			return nil
		}

		return b.ForEach(func(_, v []byte) error {
			var link models.BugLink
			if err := json.Unmarshal(v, &link); err != nil {
				return err
			}
			if link.BugID == bugID {
				links = append(links, &link)
			}
			return nil
		})
	})

	return links, err
}

func CreateBugLink(link *models.BugLink) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugLinksBucket)
		if b == nil {
			return fmt.Errorf("bug links bucket not found")
		}

		data, err := json.Marshal(link)
		if err != nil {
			return err
		}

		return b.Put(itob(link.ID), data)
	})
}

func DeleteBugLink(id int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugLinksBucket)
		if b == nil {
			return nil
		}
		return b.Delete(itob(id))
	})
}

func DeleteBugLinksByBugID(bugID int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugLinksBucket)
		if b == nil {
			return nil
		}

		var toDelete [][]byte

		_ = b.ForEach(func(k, v []byte) error {
			var link models.BugLink
			if err := json.Unmarshal(v, &link); err == nil {
				if link.BugID == bugID {
					toDelete = append(toDelete, k)
				}
			}
			return nil
		})

		for _, k := range toDelete {
			_ = b.Delete(k)
		}

		return nil
	})
}
