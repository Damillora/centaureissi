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

func (repo *CentaureissiRepository) UpdateMailbox(mailboxSchema *schema.Mailbox) error {
	existingMbox, err := repo.ExistsMailboxById(mailboxSchema.ID)
	if err != nil {
		return err
	}
	if !existingMbox {
		return errors.New("user does not exists")
	}

	mailboxProto := &pb.Mailbox{
		Id:          mailboxSchema.ID,
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

		err := b.Put([]byte(mailboxProto.Id), userData)
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
