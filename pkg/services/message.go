package services

import (
	"fmt"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/models"
	"github.com/Damillora/centaureissi/pkg/search"
	"github.com/google/uuid"
	"golang.org/x/crypto/blake2b"
)

func (cs *CentaureissiService) CounterMessagesInMailbox(id string) (uint32, error) {
	return cs.repository.CounterMessagesInMailbox(id), nil
}

func (cs *CentaureissiService) CounterMessagesInMailboxByFlag(id string, flag string) (uint32, error) {
	return cs.repository.CounterMessagesInMailboxByFlag(id, flag), nil
}

func (cs *CentaureissiService) ListMessageByMailboxId(id string) ([]*schema.Message, error) {
	msgs, err := cs.repository.ListMessagesByMailboxId(id)
	if err != nil {
		return nil, err
	}
	return msgs, err
}
func (cs *CentaureissiService) GetMessageContent(hash string) ([]byte, error) {
	blob, err := cs.blobs.GetBlob(hash)
	if err != nil {
		return nil, err
	}
	return blob, nil
}
func (cs *CentaureissiService) UploadMessage(mailboxId string, msg *models.MessageCreateModel) (*schema.Message, error) {
	uid, err := cs.IncrementMailboxUid(mailboxId)
	msgData := &schema.Message{
		Id:        uuid.NewString(),
		Hash:      msg.Hash,
		MailboxId: mailboxId,
		Uid:       uid,
		Size:      msg.Size,
		Flags:     msg.Flags,
	}
	if err != nil {
		return nil, err
	}

	err = cs.repository.CreateMessage(msgData)
	if err != nil {
		return nil, err
	}
	return msgData, nil
}
func (cs *CentaureissiService) UploadMessageContent(content []byte) (string, error) {
	sum := blake2b.Sum512(content)
	hash := fmt.Sprintf("%x", sum)

	_, err := cs.blobs.SetBlob(hash, content)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (cs *CentaureissiService) IndexSearchDocument(msg search.CentaureissiSearchDocument) error {
	err := cs.search.Index(msg)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CentaureissiService) UpdateMessageFlags(messageId string, msgSchema *models.MessageUpdateFlagsModel) error {
	msg, err := cs.repository.GetMessageById(messageId)
	if err != nil {
		return err
	}
	msg.Flags = msgSchema.Flags

	err = cs.repository.UpdateMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CentaureissiService) DeleteMessage(msgId string) error {
	err := cs.repository.DeleteMessage(msgId)
	if err != nil {
		return err
	}
	err = cs.search.Unindex(msgId)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CentaureissiService) SearchMessages(userId string, q string) (*models.SearchResponse, error) {
	result, err := cs.search.Search(userId, q)
	if err != nil {
		return nil, err
	}
	hits := make([]*models.SearchResponseItem, 0)
	for _, hit := range result.Hits {
		item := &models.SearchResponseItem{
			Id:        hit.Id,
			Hash:      hit.Hash,
			MailboxId: hit.MailboxId,
			Sender:    hit.Sender,
			From:      hit.From,
			To:        hit.To,
			Cc:        hit.Cc,
			Bcc:       hit.Bcc,
			Subject:   hit.Subject,
			Date:      hit.Date,
			Content:   hit.Content,
		}
		hits = append(hits, item)
	}
	response := &models.SearchResponse{
		Hits: hits,
	}
	return response, nil
}
