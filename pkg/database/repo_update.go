package database

import (
	"errors"

	"github.com/Damillora/centaureissi/pkg/database/pb"
	"github.com/Damillora/centaureissi/pkg/database/schema"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (repo *CentaureissiRepository) UpdateUser(userSchema *schema.User) error {
	existingUser, err := repo.ExistsUserById(userSchema.ID)
	if err != nil {
		return err
	}
	if !existingUser {
		return errors.New("user does not exists")
	}

	userProto := &pb.User{
		Id:        userSchema.ID,
		Username:  userSchema.Username,
		Password:  userSchema.Password,
		CreatedAt: timestamppb.New(userSchema.CreatedAt),
		UpdatedAt: timestamppb.Now(),
	}
	userData, err := proto.Marshal(userProto)
	if err != nil {
		return err
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_user))

		err := b.Put([]byte(userProto.Id), userData)
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

func (repo *CentaureissiRepository) UpdateMailboxOldNameIndex(mailboxId string, oldName string) error {
	mbox, err := repo.GetMailboxById(mailboxId)
	if err != nil {
		return err
	}
	if mbox == nil {
		return errors.New("mailbox does not exists")
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		imuin := tx.Bucket([]byte(index_mailbox_user_id_name))

		// Map user ID and mbox name into index
		err = imuin.Delete([]byte(formatUserIdAndName(mbox.UserId, mbox.Name)))
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
func (repo *CentaureissiRepository) UpdateMailbox(mailboxSchema *schema.Mailbox) error {
	existingMbox, err := repo.ExistsMailboxById(mailboxSchema.Id)
	if err != nil {
		return err
	}
	if !existingMbox {
		return errors.New("user does not exists")
	}

	mailboxProto := &pb.Mailbox{
		Id:          mailboxSchema.Id,
		UserId:      mailboxSchema.UserId,
		Name:        mailboxSchema.Name,
		UidValidity: mailboxSchema.UidValidity,
		Subscribed:  mailboxSchema.Subscribed,
		CreatedAt:   timestamppb.New(mailboxSchema.CreatedAt),
		UpdatedAt:   timestamppb.Now(),
	}
	userData, err := proto.Marshal(mailboxProto)
	if err != nil {
		return err
	}

	err = repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_mailbox))
		imuin := tx.Bucket([]byte(index_mailbox_user_id_name))

		err := b.Put([]byte(mailboxProto.Id), userData)
		if err != nil {
			return err
		}

		// Map user ID and mbox name into index
		err = imuin.Put([]byte(formatUserIdAndName(mailboxSchema.UserId, mailboxSchema.Name)), []byte(mailboxSchema.Id))
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

func (repo *CentaureissiRepository) UpdateMessage(messageSchema *schema.Message) error {
	msg, err := repo.ExistsMessageById(messageSchema.Id)
	if err != nil {
		return err
	}
	if !msg {
		return errors.New("message does not exists")
	}
	messageProto := &pb.Message{
		Id:        messageSchema.Id,
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
		err := bm.Put([]byte(messageSchema.Id), messageData)
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
