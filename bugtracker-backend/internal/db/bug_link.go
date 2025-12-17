package db

import (
	"encoding/json"
	"bugtracker-backend/internal/models"
	"go.etcd.io/bbolt"
)

func CreateBugLink(link *models.BugLink) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugLinksBucket)
		id := btoi(b.Get([]byte("lastID"))) + 1
		b.Put([]byte("lastID"), itob(id))

		link.ID = id
		data, _ := json.Marshal(link)
		return b.Put(itob(id), data)
	})
}

func GetBugLinks(bugID int) ([]models.BugLink, error) {
	var links []models.BugLink

	err := db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bugLinksBucket).ForEach(func(_, v []byte) error {
			var l models.BugLink
			json.Unmarshal(v, &l)
			if l.BugID == bugID {
				links = append(links, l)
			}
			return nil
		})
	})

	return links, err
}

func DeleteBugLinksByBugID(bugID int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugLinksBucket)
		var del [][]byte

		b.ForEach(func(k, v []byte) error {
			var l models.BugLink
			json.Unmarshal(v, &l)
			if l.BugID == bugID {
				del = append(del, k)
			}
			return nil
		})

		for _, k := range del {
			b.Delete(k)
		}
		return nil
	})
}
