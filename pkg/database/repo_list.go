package database

import (
	"errors"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	bolt "go.etcd.io/bbolt"
)

func (repo *CentaureissiRepository) ListMailboxesByUserId(id string) ([]*schema.Mailbox, error) {
	exists, err := repo.ExistsUserById(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	mailboxes := make([]*schema.Mailbox, 0)
	repo.db.View(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(bucket_user_mailbox)).Bucket([]byte(id))

		c := mb.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			mailboxId := string(k)
			mailbox, err := repo.GetMailboxById(mailboxId)
			if err != nil {
				return err
			}
			mailboxes = append(mailboxes, mailbox)
		}
		return nil
	})
	return mailboxes, nil
}

func (repo *CentaureissiRepository) ListMessagesByMailboxId(id string) ([]*schema.Message, error) {
	exists, err := repo.ExistsMailboxById(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("mailbox not found")
	}

	messages := make([]*schema.Message, 0)
	repo.db.View(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(bucket_mailbox_message)).Bucket([]byte(id))

		c := mb.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			messageId := string(k)
			message, err := repo.GetMessageById(messageId)
			if err != nil {
				return err
			}
			if message != nil {
				messages = append(messages, message)
			}
		}
		return nil
	})
	return messages, nil
}
