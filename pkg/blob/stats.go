package blob

import (
	"os"

	"github.com/Damillora/centaureissi/pkg/config"
	bolt "go.etcd.io/bbolt"
)

func (cbr *CentaureissiBlobRepository) StatsDbSize() uint64 {

	databaseUrl := config.CurrentConfig.DataDirectory + "/blobs.bolt"

	info, _ := os.Stat(databaseUrl)
	return uint64(info.Size())
}

func (cbr *CentaureissiBlobRepository) StatsBlobs() uint64 {
	var keyCount uint64
	err := cbr.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_blob))
		keyCount = uint64(b.Stats().KeyN)
		return nil
	})
	if err != nil {
		return 0
	}

	return keyCount
}
