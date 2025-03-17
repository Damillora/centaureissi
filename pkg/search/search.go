package search

import (
	"github.com/blevesearch/bleve/v2"
)

func (cse *CentaureissiSearchEngine) Index(msg CentaureissiSearchDocument) error {
	err := cse.index.Index(msg.Id, msg)
	if err != nil {
		return err
	}
	return nil
}

func (cse *CentaureissiSearchEngine) Unindex(msgId string) error {
	err := cse.index.Delete(msgId)
	if err != nil {
		return err
	}
	return nil
}

func (cse *CentaureissiSearchEngine) Search(userId string, q string) (*CentaureissiSearchResponse, error) {
	allFields, _ := cse.index.Fields()
	query := bleve.NewQueryStringQuery(q)
	searchReq := bleve.NewSearchRequest(query)
	searchReq.Fields = allFields
	result, err := cse.index.Search(searchReq)
	if err != nil {
		return nil, err
	}

	searchResult := make([]*CentaureissiSearchResult, 0)
	for _, hit := range result.Hits {
		item := &CentaureissiSearchResult{
			Id:        hit.Fields["Id"].(string),
			Hash:      hit.Fields["Hash"].(string),
			MailboxId: hit.Fields["MailboxId"].(string),

			Sender:  hit.Fields["Sender"].(string),
			From:    hit.Fields["From"].(string),
			To:      hit.Fields["To"].(string),
			Cc:      hit.Fields["Cc"].(string),
			Bcc:     hit.Fields["Bcc"].(string),
			Subject: hit.Fields["Subject"].(string),
			Date:    hit.Fields["Date"].(string),
			Content: hit.Fields["Content"].(string),
		}
		searchResult = append(searchResult, item)
	}

	response := &CentaureissiSearchResponse{
		Hits: searchResult,
	}
	return response, nil
}
