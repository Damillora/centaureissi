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
	Hits []*SearchResponseItem `json:"items"`
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
