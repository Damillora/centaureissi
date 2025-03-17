package database

import (
	"github.com/Damillora/centaureissi/pkg/database/pb"
	"github.com/Damillora/centaureissi/pkg/database/schema"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func (repo *CentaureissiRepository) GetUserById(id string) (*schema.User, error) {
	var userData []byte

	// Read data bytes from DB
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_user))
		userData = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if userData == nil {
		return nil, nil
	}

	// Unmarshal protobuf
	userProto := &pb.User{}
	if err := proto.Unmarshal(userData, userProto); err != nil {
		return nil, err
	}

	user := &schema.User{
		ID:       userProto.Id,
		Username: userProto.Username,
		Password: userProto.Password,
	}

	return user, nil
}

func (repo *CentaureissiRepository) GetUserByUsername(username string) (*schema.User, error) {
	var userId string
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(index_user_username))
		userId = string(b.Get([]byte(username)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	user, err := repo.GetUserById(userId)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *CentaureissiRepository) GetMailboxById(id string) (*schema.Mailbox, error) {
	var mailboxData []byte
	// Read data bytes from DB
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_mailbox))
		mailboxData = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if mailboxData == nil {
		return nil, nil
	}

	// Unmarshal protobuf
	mailboxProto := &pb.Mailbox{}
	if err := proto.Unmarshal(mailboxData, mailboxProto); err != nil {
		return nil, err
	}

	mailbox := &schema.Mailbox{
		Id:          mailboxProto.Id,
		UserId:      mailboxProto.UserId,
		UidValidity: mailboxProto.UidValidity,
		Name:        mailboxProto.Name,
		Subscribed:  mailboxProto.Subscribed,
		CreatedAt:   mailboxProto.CreatedAt.AsTime(),
		UpdatedAt:   mailboxProto.UpdatedAt.AsTime(),
	}

	return mailbox, nil
}

func (repo *CentaureissiRepository) GetMailboxByUserIdAndName(userId string, mailboxName string) (*schema.Mailbox, error) {
	var mailboxId string
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(index_mailbox_user_id_name))
		mailboxId = string(b.Get([]byte(formatUserIdAndName(userId, mailboxName))))
		return nil
	})
	if err != nil {
		return nil, err
	}
	mailbox, err := repo.GetMailboxById(mailboxId)

	if err != nil {
		return nil, err
	}

	return mailbox, nil
}

func (repo *CentaureissiRepository) GetMessageById(id string) (*schema.Message, error) {
	var messageData []byte
	// Read data bytes from DB
	err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_message))
		messageData = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if messageData == nil {
		return nil, nil
	}

	// Unmarshal protobuf
	messageProto := &pb.Message{}
	if err := proto.Unmarshal(messageData, messageProto); err != nil {
		return nil, err
	}

	message := &schema.Message{
		Id:        messageProto.Id,
		Hash:      messageProto.Hash,
		MailboxId: messageProto.MailboxId,
		Uid:       messageProto.Uid,
		Size:      messageProto.Size,
		Flags:     messageProto.Flags,
	}

	return message, nil
}
