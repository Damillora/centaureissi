package models

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UserProfileResponse struct {
	Username string `json:"username"`
}

type SearchResponse struct {
	Hits       []*SearchResponseItem `json:"items"`
	Page       int                   `json:"page"`
	TotalPages int                   `json:"totalPages"`
	Count      uint64                `json:"total"`
}

type SearchResponseItem struct {
	Id        string `json:"id"`
	Hash      string `json:"hash"`
	MailboxId string `json:"mailboxId"`

	Sender  string `json:"sender"`
	From    string `json:"from"`
	To      string `json:"to"`
	Cc      string `json:"cc"`
	Bcc     string `json:"bcc"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
}

type CentaureissiStatsResponse struct {
	Version        string `json:"version"`
	DbSize         uint64 `json:"dbSize"`
	MailboxCount   uint64 `json:"mailboxCount"`
	MessageCount   uint64 `json:"messageCount"`
	BlobDbSize     uint64 `json:"blobDbSize"`
	BlobCount      uint64 `json:"blobCount"`
	SearchDocCount uint64 `json:"searchDocCount"`
}
