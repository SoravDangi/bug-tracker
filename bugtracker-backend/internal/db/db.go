package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"bugtracker-backend/internal/models"

	"go.etcd.io/bbolt"
)

var (
	db          *bbolt.DB
	initialized bool

	bugsBucket     = []byte("bugs")
	commentsBucket = []byte("comments")
	counterBucket  = []byte("counter")
    

    // bug link bucket
	bugLinksBucket = []byte("bug_links")

     

	// Screenshot buckets
	screenshotsBucket       = []byte("screenshots")
	screenshotCounterBucket = []byte("screenshot_counter")

	databasePath = getDBPath()
)

func getDBPath() string {
	if path := os.Getenv("DB_PATH"); path != "" {
		return path
	}
	return "bugs.db"
}

func Init() error {
	if initialized {
		return fmt.Errorf("database already initialized")
	}

	var err error
	db, err = bbolt.Open(databasePath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {

		// Buckets
		if _, err := tx.CreateBucketIfNotExists(bugsBucket); err != nil {
			return fmt.Errorf("create bugs bucket: %w", err)
		}

		if _, err := tx.CreateBucketIfNotExists(commentsBucket); err != nil {
			return fmt.Errorf("create comments bucket: %w", err)
		}

		counter, err := tx.CreateBucketIfNotExists(counterBucket)
		if err != nil {
			return fmt.Errorf("create counter bucket: %w", err)
		}
		if counter.Get([]byte("lastBugID")) == nil {
			if err := counter.Put([]byte("lastBugID"), itob(0)); err != nil {
				return fmt.Errorf("init bug counter: %w", err)
			}
		}

		if _, err := tx.CreateBucketIfNotExists(screenshotsBucket); err != nil {
			return fmt.Errorf("create screenshots bucket: %w", err)
		}

		sCounter, err := tx.CreateBucketIfNotExists(screenshotCounterBucket)
		if err != nil {
			return fmt.Errorf("create screenshot counter bucket: %w", err)
		}
		if sCounter.Get([]byte("lastScreenshotID")) == nil {
			if err := sCounter.Put([]byte("lastScreenshotID"), itob(0)); err != nil {
				return fmt.Errorf("init screenshot counter: %w", err)
			}
		}
		if _, err := tx.CreateBucketIfNotExists(bugLinksBucket); err != nil {
	        return fmt.Errorf("create bug links bucket: %w", err) 
        }


		return nil
	})

	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	initialized = true
	return nil
}

func CreateBug(bug *models.Bug) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugsBucket)

		id, err := getNextBugID(tx)
		if err != nil {
			return err
		}
		bug.ID = id

		data, err := json.Marshal(bug)
		if err != nil {
			return err
		}

		return b.Put(itob(bug.ID), data)
	})
}

func GetBug(id int) (*models.Bug, error) {
	var bug models.Bug

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugsBucket)
		data := b.Get(itob(id))
		if data == nil {
			return fmt.Errorf("bug not found")
		}
		return json.Unmarshal(data, &bug)
	})

	if err != nil {
		return nil, err
	}
	return &bug, nil
}

func GetAllBugs() ([]*models.Bug, error) {
	var bugs []*models.Bug

	err := db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bugsBucket).ForEach(func(_, v []byte) error {
			var bug models.Bug
			if err := json.Unmarshal(v, &bug); err != nil {
				return err
			}
			bugs = append(bugs, &bug)
			return nil
		})
	})

	return bugs, err
}

func UpdateBug(bug *models.Bug) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugsBucket)

		if b.Get(itob(bug.ID)) == nil {
			return fmt.Errorf("bug not found")
		}

		bug.UpdatedAt = time.Now()

		data, err := json.Marshal(bug)
		if err != nil {
			return err
		}

		return b.Put(itob(bug.ID), data)
	})
}

func DeleteBug(id int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugsBucket)
		if b.Get(itob(id)) == nil {
			return fmt.Errorf("bug not found")
		}
		return b.Delete(itob(id))
	})
}

func DeleteAllBugs() (int, error) {
	var count int

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bugsBucket)
		count = b.Stats().KeyN

		if err := tx.DeleteBucket(bugsBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucket(bugsBucket); err != nil {
			return err
		}

		counter := tx.Bucket(counterBucket)
		return counter.Put([]byte("lastBugID"), itob(0))
	})

	return count, err
}

func Cleanup() {
	if db != nil {
		db.Close()
		db = nil
	}
	initialized = false
}

func getNextBugID(tx *bbolt.Tx) (int, error) {
	b := tx.Bucket(counterBucket)
	id := btoi(b.Get([]byte("lastBugID"))) + 1
	if err := b.Put([]byte("lastBugID"), itob(id)); err != nil {
		return 0, err
	}
	return id, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	if b == nil {
		return 0
	}
	return int(binary.BigEndian.Uint64(b))
}

func CleanupTestDB() error {
	if db == nil {
		return nil
	}

	return db.Update(func(tx *bbolt.Tx) error {
		// Delete existing buckets (ignore errors if they don't exist)
		_ = tx.DeleteBucket(bugsBucket)
		_ = tx.DeleteBucket(commentsBucket)
		_ = tx.DeleteBucket(counterBucket)
		_ = tx.DeleteBucket(screenshotsBucket)
		_ = tx.DeleteBucket(screenshotCounterBucket)
		_ = tx.DeleteBucket(bugLinksBucket)

		// Recreate buckets
		if _, err := tx.CreateBucket(bugsBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucket(commentsBucket); err != nil {
			return err
		}

		counter, err := tx.CreateBucket(counterBucket)
		if err != nil {
			return err
		}
		if err := counter.Put([]byte("lastBugID"), itob(0)); err != nil {
			return err
		}

		if _, err := tx.CreateBucket(screenshotsBucket); err != nil {
			return err
		}

		sCounter, err := tx.CreateBucket(screenshotCounterBucket)
		if err != nil {
			return err
		}
		if err := sCounter.Put([]byte("lastScreenshotID"), itob(0)); err != nil {
			return err
		}
		if _, err := tx.CreateBucket(bugLinksBucket); err != nil {
	        return err
        }

		return nil
	})
}
