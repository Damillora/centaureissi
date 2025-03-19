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

func (repo *CentaureissiRepository) DeleteMessage(messageId string) error {
	msg, err := repo.GetMessageById(messageId)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message does not exists")
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		bm := tx.Bucket([]byte(bucket_message))
		bmm := tx.Bucket([]byte(bucket_mailbox_message)).Bucket([]byte(msg.MailboxId))
		immuid := tx.Bucket([]byte(index_message_mailbox_uid))
		imiu := tx.Bucket([]byte(index_message_id_uid))

		// Delete message
		err := bm.Delete([]byte(msg.Id))
		if err != nil {
			return err
		}
		// Delete message in mailbox
		err = bmm.Delete([]byte(msg.Id))
		if err != nil {
			return err
		}
		// Delete user ID and mbox name index
		err = immuid.Delete([]byte(formatMailboxIdAndUid(msg.MailboxId, msg.Uid)))
		if err != nil {
			return err
		}
		// Delete message id and UID index
		err = imiu.Delete([]byte(msg.Id))
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
