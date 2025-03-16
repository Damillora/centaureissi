package services

import (
	"github.com/Damillora/centaureissi/pkg/blob"
	"github.com/Damillora/centaureissi/pkg/database"
	"github.com/Damillora/centaureissi/pkg/search"
)

type CentaureissiService struct {
	repository *database.CentaureissiRepository
	blobs      *blob.CentaureissiBlobRepository
	search     *search.CentaureissiSearchEngine
}

func New(repo *database.CentaureissiRepository, b *blob.CentaureissiBlobRepository, s *search.CentaureissiSearchEngine) *CentaureissiService {
	return &CentaureissiService{
		repository: repo,
		blobs:      b,
		search:     s,
	}
}
