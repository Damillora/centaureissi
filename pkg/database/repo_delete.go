package database

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

func (repo *CentaureissiRepository) DeleteMailbox(mailboxId string) error {
	existingMbox, err := repo.GetMailboxById(mailboxId)
	if err != nil {
		return err
	}
	if existingMbox == nil {
		return errors.New("user does not exists")
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		bm := tx.Bucket([]byte(bucket_mailbox))
		bum := tx.Bucket([]byte(bucket_user_mailbox)).Bucket([]byte(existingMbox.UserId))
		imuin := tx.Bucket([]byte(index_mailbox_user_id_name))

		// Delete mailbox
		err := bm.Delete([]byte(mailboxId))
		if err != nil {
			return err
		}
		// Delete in user's mailbox list
		err = bum.Delete([]byte(mailboxId))
		if err != nil {
			return err
		}
		// Delete user ID and mbox name index
		err = imuin.Delete([]byte(formatUserIdAndName(existingMbox.UserId, existingMbox.Name)))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
