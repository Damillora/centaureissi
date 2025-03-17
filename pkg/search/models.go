package search

import "time"

type CentaureissiSearchDocument struct {
	Id        string
	Hash      string
	UserId    string
	MailboxId string

	Sender  string
	From    string
	To      string
	Cc      string
	Bcc     string
	Subject string
	Date    time.Time
	Content string
}
type CentaureissiSearchResponse struct {
	Hits []*CentaureissiSearchResult
}

type CentaureissiSearchResult struct {
	Id        string
	Hash      string
	MailboxId string

	Sender  string
	From    string
	To      string
	Cc      string
	Bcc     string
	Subject string
	Date    string
	Content string
}
