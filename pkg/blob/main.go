package blob

import (
	"fmt"
	"log"

	"github.com/Damillora/centaureissi/pkg/config"
	bolt "go.etcd.io/bbolt"
)

const bucket_blob = "blobs"

type CentaureissiBlobRepository struct {
	db *bolt.DB
}

func New() *CentaureissiBlobRepository {
	cbr := &CentaureissiBlobRepository{}
	cbr.Initialize()
	return cbr
}

func (cbr *CentaureissiBlobRepository) Initialize() {
	databaseUrl := config.CurrentConfig.DataDirectory + "/blobs.db"

	dbConn, err := bolt.Open(databaseUrl, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Auto create buckets
	err = dbConn.Update(func(tx *bolt.Tx) error {
		// Tables
		_, err := tx.CreateBucketIfNotExists([]byte(bucket_blob))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	cbr.db = dbConn
}
func (cbr *CentaureissiBlobRepository) GetBlob(hash string) ([]byte, error) {
	var blob []byte
	err := cbr.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_blob))
		blob = b.Get([]byte(hash))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func (cbr *CentaureissiBlobRepository) SetBlob(hash string, content []byte) (string, error) {
	err := cbr.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_blob))
		b.Put([]byte(hash), content)
		return nil
	})
	if err != nil {
		return "", err
	}
	return hash, nil

}
