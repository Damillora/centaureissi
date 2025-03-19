package services

import "github.com/Damillora/centaureissi/pkg/models"

func (cs *CentaureissiService) Stats() *models.CentaureissiStatsResponse {
	return &models.CentaureissiStatsResponse{
		Version:        "0.1.0",
		DbSize:         cs.repository.StatsDbSize(),
		MailboxCount:   cs.repository.StatsMailboxes(),
		MessageCount:   cs.repository.StatsMessages(),
		BlobDbSize:     cs.blobs.StatsDbSize(),
		BlobCount:      cs.blobs.StatsBlobs(),
		SearchDocCount: cs.search.StatsDocuments(),
	}
}
