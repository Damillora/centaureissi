package database

import (
	"fmt"
	"log"

	"github.com/Damillora/centaureissi/pkg/config"
	bolt "go.etcd.io/bbolt"
)

type CentaureissiRepository struct {
	db *bolt.DB
}

func New() *CentaureissiRepository {
	repo := &CentaureissiRepository{}
	repo.Initialize()
	return repo
}

func (cdb *CentaureissiRepository) Initialize() {
	databaseUrl := config.CurrentConfig.DataDirectory + "/centaureissi.db"

	dbConn, err := bolt.Open(databaseUrl, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Auto create buckets
	err = dbConn.Update(func(tx *bolt.Tx) error {
		// Tables
		_, err := tx.CreateBucketIfNotExists([]byte(bucket_user))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(bucket_mailbox))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(bucket_message))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		// Indexes
		_, err = tx.CreateBucketIfNotExists([]byte(index_user_username))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(index_mailbox_user_id_name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(index_message_mailbox_uid))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		// Counters
		_, err = tx.CreateBucketIfNotExists([]byte(counter_uidvalidity))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(counter_uid))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	cdb.db = dbConn
}

func (cdb *CentaureissiRepository) Deinitialize() {
	cdb.db.Close()
}
