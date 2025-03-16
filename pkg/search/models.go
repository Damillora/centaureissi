package search

import "time"

type CentaureissiSearchDocument struct {
	Hash    string
	Sender  string
	From    string
	To      string
	Content string
	Cc      string
	Bcc     string
	Subject string
	Date    time.Time
}
