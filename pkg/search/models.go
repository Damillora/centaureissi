package search

import "time"

type CentaureissiSearchDocument struct {
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
