package database

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

func (repo *CentaureissiRepository) CounterUidValidity(id string) (uint32, error) {
	userExists, err := repo.ExistsUserById(id)
	if err != nil {
		return 0, err
	}
	if !userExists {
		return 0, errors.New("user does not exist")
	}

	var uidValidity uint32
	err = repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(counter_uidvalidity))
		counterData := b.Get([]byte(id))
		if counterData == nil {
			counterData = []byte(uint32ToString(0))
		}
		uidValidity = stringToUint32(string(counterData))
		return nil
	})
	if err != nil {
		return 0, err
	}

	return uidValidity, nil
}
func (repo *CentaureissiRepository) IncrementUidValidity(id string) (uint32, error) {
	userExists, err := repo.ExistsUserById(id)
	if err != nil {
		return 0, err
	}
	if !userExists {
		return 0, errors.New("user does not exist")
	}

	var uidValidity uint32
	err = repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(counter_uidvalidity))
		counterData := b.Get([]byte(id))
		if counterData == nil {
			counterData = []byte(uint32ToString(0))
		}
		uidValidity = stringToUint32(string(counterData))
		uidValidity++
		err := b.Put([]byte(id), []byte(uint32ToString(uidValidity)))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return uidValidity, nil
}

func (repo *CentaureissiRepository) CounterMailboxUid(id string) (uint32, error) {
	mailboxExists, err := repo.ExistsMailboxById(id)
	if err != nil {
		return 0, err
	}
	if !mailboxExists {
		return 0, errors.New("user does not exist")
	}

	var uid uint32
	err = repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(counter_uidvalidity))
		counterData := b.Get([]byte(id))
		if counterData == nil {
			counterData = []byte(uint32ToString(0))
		}
		uid = stringToUint32(string(counterData))
		return nil
	})
	if err != nil {
		return 0, err
	}

	return uid, nil
}

func (repo *CentaureissiRepository) IncrementMailboxUid(id string) (uint32, error) {
	mailboxExists, err := repo.ExistsMailboxById(id)
	if err != nil {
		return 0, err
	}
	if !mailboxExists {
		return 0, errors.New("user does not exist")
	}

	var uid uint32
	err = repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(counter_uidvalidity))
		counterData := b.Get([]byte(id))
		if counterData == nil {
			counterData = []byte(uint32ToString(0))
		}
		uid = stringToUint32(string(counterData))
		uid++
		err := b.Put([]byte(id), []byte(uint32ToString(uid)))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return uid, nil
}

func (repo *CentaureissiRepository) CounterMessagesInMailbox(mailboxId string) uint32 {
	var keyCount uint32
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_mailbox))
		bmm := b.Bucket([]byte(bucket_mailbox_message))
		keyCount = uint32(bmm.Stats().KeyN)
		return nil
	})
	if err != nil {
		return 0
	}

	return keyCount
}
