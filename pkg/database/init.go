package database

import (
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
	databaseUrl := config.CurrentConfig.DataDirectory + "/centaureissi.bolt"

	dbConn, err := bolt.Open(databaseUrl, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	cdb.db = dbConn

	err = cdb.Migrate()
	if err != nil {
		log.Fatal(err)
	}

}

func (cdb *CentaureissiRepository) Deinitialize() {
	cdb.db.Close()
}
