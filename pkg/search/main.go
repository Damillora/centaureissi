package search

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/Damillora/centaureissi/pkg/config"
	"github.com/blevesearch/bleve/v2"
)

type CentaureissiSearchEngine struct {
	index bleve.Index
}

func New() *CentaureissiSearchEngine {
	search := &CentaureissiSearchEngine{}
	search.Initialize()
	return search
}

func (cse *CentaureissiSearchEngine) Initialize() {
	databaseUrl := config.CurrentConfig.DataDirectory + "/search"
	mapping := bleve.NewIndexMapping()

	_, err := os.Stat(databaseUrl)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatal("Cannot init search!")
	} else if err != nil {
		cse.index, err = bleve.New(databaseUrl, mapping)
		if err != nil {
			log.Fatal("Cannot init search!")
		}
	} else {
		cse.index, err = bleve.Open(databaseUrl)
		if err != nil {
			log.Fatal("Cannot init search!")
		}
	}
}
