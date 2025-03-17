package database

import (
	"errors"

	"github.com/Damillora/centaureissi/pkg/database/pb"
	"github.com/Damillora/centaureissi/pkg/database/schema"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (repo *CentaureissiRepository) CreateUser(userSchema *schema.User) error {
	existingUser, err := repo.ExistsUserById(userSchema.ID)
	if err != nil {
		return err
	}
	if existingUser {
		return errors.New("user already exists")
	}

	userProto := &pb.User{
		Id:        userSchema.ID,
		Username:  userSchema.Username,
		Password:  userSchema.Password,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}
	userData, err := proto.Marshal(userProto)
	if err != nil {
		return err
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_user))
		bm := tx.Bucket([]byte(bucket_user_mailbox))
		err = b.Put([]byte(userProto.Id), userData)
		if err != nil {
			return err
		}

		index := tx.Bucket([]byte(index_user_username))
		err = index.Put([]byte(userProto.Username), []byte(userProto.Id))
		if err != nil {
			return err
		}
		// Create bucket for mailboxes
		_, err = bm.CreateBucketIfNotExists([]byte(userSchema.ID))
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

func (repo *CentaureissiRepository) CreateMailbox(mailboxSchema *schema.Mailbox) error {
	user, err := repo.GetUserById(mailboxSchema.UserId)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user does not exists")
	}

	mailboxProto := &pb.Mailbox{
		Id:          mailboxSchema.Id,
		UserId:      mailboxSchema.UserId,
		Name:        mailboxSchema.Name,
		UidValidity: mailboxSchema.UidValidity,
		Subscribed:  mailboxSchema.Subscribed,
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	}
	mailboxData, err := proto.Marshal(mailboxProto)
	if err != nil {
		return err
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		bm := tx.Bucket([]byte(bucket_mailbox))
		bum := tx.Bucket([]byte(bucket_user_mailbox)).Bucket([]byte(mailboxSchema.UserId))
		bmm := tx.Bucket([]byte(bucket_mailbox_message))
		imuin := tx.Bucket([]byte(index_mailbox_user_id_name))

		// Insert into Mailbox
		err := bm.Put([]byte(mailboxSchema.Id), mailboxData)
		if err != nil {
			return err
		}
		// Add into User's mailbox list
		err = bum.Put([]byte(mailboxSchema.Id), []byte(mailboxSchema.UserId))
		if err != nil {
			return err
		}
		// Map user ID and mbox name into index
		err = imuin.Put([]byte(formatUserIdAndName(mailboxSchema.UserId, mailboxSchema.Name)), []byte(mailboxSchema.Id))
		if err != nil {
			return err
		}
		// Create bucket for messages
		_, err = bmm.CreateBucketIfNotExists([]byte(mailboxSchema.Id))
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

func (repo *CentaureissiRepository) CreateMessage(messageSchema *schema.Message) error {
	mailbox, err := repo.GetMailboxById(messageSchema.MailboxId)
	if err != nil {
		return err
	}
	if mailbox == nil {
		return errors.New("mailbox does not exists")
	}
	messageProto := &pb.Message{
		Hash:      messageSchema.Hash,
		MailboxId: messageSchema.MailboxId,
		Uid:       messageSchema.Uid,
		Size:      messageSchema.Size,
		Flags:     messageSchema.Flags,
	}
	messageData, err := proto.Marshal(messageProto)
	if err != nil {
		return err
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		bm := tx.Bucket([]byte(bucket_message))
		bmm := tx.Bucket([]byte(bucket_mailbox_message)).Bucket([]byte(messageSchema.MailboxId))
		immuid := tx.Bucket([]byte(index_message_mailbox_uid))

		// Insert into Message List
		err := bm.Put([]byte(messageSchema.Hash), messageData)
		if err != nil {
			return err
		}
		// Add into mailbox list
		err = bmm.Put([]byte(messageSchema.Hash), []byte{})
		if err != nil {
			return err
		}
		// Map user ID and mbox name into index
		err = immuid.Put([]byte(formatMailboxIdAndUid(messageSchema.MailboxId, messageSchema.Uid)), []byte(messageSchema.Hash))
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
