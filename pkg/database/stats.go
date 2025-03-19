package database

import (
	"os"

	"github.com/Damillora/centaureissi/pkg/config"
	bolt "go.etcd.io/bbolt"
)

func (repo *CentaureissiRepository) StatsDbSize() uint64 {

	databaseUrl := config.CurrentConfig.DataDirectory + "/centaureissi.bolt"

	info, _ := os.Stat(databaseUrl)
	return uint64(info.Size())
}

func (repo *CentaureissiRepository) StatsMailboxes() uint64 {
	var keyCount uint64
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_mailbox))
		keyCount = uint64(b.Stats().KeyN)
		return nil
	})
	if err != nil {
		return 0
	}

	return keyCount
}

func (repo *CentaureissiRepository) StatsMessages() uint64 {
	var keyCount uint64
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_message))
		keyCount = uint64(b.Stats().KeyN)
		return nil
	})
	if err != nil {
		return 0
	}

	return keyCount
}
