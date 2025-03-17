package database

import (
	bolt "go.etcd.io/bbolt"
)

func (repo *CentaureissiRepository) ExistsUserById(id string) (bool, error) {
	var userData []byte
	// Read data bytes from DB
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_user))
		userData = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		return false, err
	}

	return userData != nil, nil
}
func (repo *CentaureissiRepository) ExistsUserByUsername(username string) (bool, error) {
	var userId []byte
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(index_user_username))
		userId = b.Get([]byte(username))
		return nil
	})
	if err != nil {
		return false, err
	}

	return userId != nil, nil
}

func (repo *CentaureissiRepository) ExistsMailboxById(id string) (bool, error) {
	var mailboxData []byte
	// Read data bytes from DB
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_mailbox))
		mailboxData = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		return false, err
	}

	return mailboxData != nil, nil
}

func (repo *CentaureissiRepository) ExistsMailboxByUserIdAndName(userId string, mailboxName string) (bool, error) {
	var mailboxId []byte
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(index_mailbox_user_id_name))
		mailboxId = b.Get([]byte(formatUserIdAndName(userId, mailboxName)))
		return nil
	})
	if err != nil {
		return false, err
	}

	return mailboxId != nil, nil
}

func (repo *CentaureissiRepository) ExistsMessageById(id string) (bool, error) {
	var messageId []byte
	// Read data bytes from DB
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_message))
		messageId = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		return false, err
	}

	return messageId != nil, nil
}
