package models

import "time"

type UserCreateModel struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateModel struct {
	Username string `json:"username" validate:"required"`
}

type UserUpdatePasswordModel struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type MailboxCreateModel struct {
	Name string `json:"name" validate:"required"`
}

type MessageCreateModel struct {
	Hash  string
	Size  uint64
	Flags map[string]bool
}

type MessageUpdateFlagsModel struct {
	Flags map[string]bool
}

type MessageIndexModel struct {
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
