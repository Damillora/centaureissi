package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Damillora/centaureissi/pkg/database/pb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func (repo *CentaureissiRepository) Migrate() error {
	err := repo.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("migrations"))
		if err != nil {
			return err
		}

		err = migration20250317001_initial(tx, b)
		if err != nil {
			return err
		}

		err = migration20250319001_messageUidIndex(tx, b)
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

func migration20250317001_initial(tx *bolt.Tx, migrations *bolt.Bucket) error {
	migration := migrations.Get([]byte("20250317001_initial"))
	if migration != nil {
		log.Println("migration 20250317001_initial done")
		return nil
	}
	log.Println("running migration 20250317001_initial...")

	// Tables
	_, err := tx.CreateBucketIfNotExists([]byte(bucket_user))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(bucket_mailbox))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(bucket_message))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(bucket_user_mailbox))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(bucket_mailbox_message))
	if err != nil {
		return err
	}

	// Indexes
	_, err = tx.CreateBucketIfNotExists([]byte(index_user_username))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(index_mailbox_user_id_name))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(index_message_mailbox_uid))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}

	// Counters
	_, err = tx.CreateBucketIfNotExists([]byte(counter_uidvalidity))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte(counter_uid))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	return nil
}

func migration20250319001_messageUidIndex(tx *bolt.Tx, migrations *bolt.Bucket) error {
	migration := migrations.Get([]byte("20250319001_messageUidIndex"))
	if migration != nil {
		log.Println("migration 20250319001_messageUidIndex done")
		return nil
	}
	log.Println("running migration 20250319001_messageUidIndex...")

	b := tx.Bucket([]byte(bucket_user))
	messages := tx.Bucket([]byte(bucket_message))

	imiu, err := tx.CreateBucketIfNotExists([]byte(index_message_id_uid))
	if err != nil {
		return err
	}

	c := b.Cursor()

	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		id := string(k)
		mb := tx.Bucket([]byte(bucket_user_mailbox)).Bucket([]byte(id))

		mc := mb.Cursor()

		for mk, _ := mc.First(); mk != nil; mk, _ = mc.Next() {
			mailboxId := string(mk)
			mmb := tx.Bucket([]byte(bucket_mailbox_message)).Bucket([]byte(mailboxId))
			mmc := mmb.Cursor()

			for mmk, _ := mmc.First(); mmk != nil; mmk, _ = mmc.Next() {
				messageId := string(mmk)
				messageData := messages.Get([]byte(messageId))

				// Unmarshal protobuf
				messageProto := &pb.Message{}
				if err := proto.Unmarshal(messageData, messageProto); err != nil {
					return err
				}

				err := imiu.Put([]byte(messageId), []byte(strconv.FormatInt(int64(messageProto.Uid), 10)))
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}
