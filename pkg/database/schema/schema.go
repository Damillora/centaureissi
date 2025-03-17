package schema

import "time"

type User struct {
	ID       string
	Username string
	Password string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Mailbox struct {
	Id          string
	UserId      string
	UidValidity uint32
	Name        string
	Subscribed  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Message struct {
	Uid       uint32
	Hash      string
	MailboxId string
	Size      uint64
	Flags     map[string]bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Headers struct {
	Header string
}

type MessageBody struct {
}
