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

		return nil
	})
}
