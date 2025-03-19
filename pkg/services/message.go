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
	msgs, err := cs.repository.ListMessageIdsByMailboxId(id)
	messages := make([]*schema.Message, 0)
	for _, msgId := range msgs {
		msg, err := cs.repository.GetMessageById(msgId)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err != nil {
		return nil, err
	}
	return messages, err
}

func (cs *CentaureissiService) ListMessageUidsByMailboxId(id string) ([]*models.MessageUidListItem, error) {
	msgs, err := cs.repository.ListMessageIdsByMailboxId(id)
	messages := make([]*models.MessageUidListItem, 0)
	for _, msgId := range msgs {
		uid, err := cs.repository.GetMessageUidById(msgId)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &models.MessageUidListItem{
			Id:  msgId,
			Uid: uid,
		})
	}
	if err != nil {
		return nil, err
	}
	return messages, err
}

func (cs *CentaureissiService) GetMessageById(id string) (*schema.Message, error) {
	return cs.repository.GetMessageById(id)
}

func (cs *CentaureissiService) GetMessageContent(hash string) ([]byte, error) {
	blob, err := cs.blobs.GetBlob(hash)
	if err != nil {
		return nil, err
	}
	return blob, nil
}
func (cs *CentaureissiService) UploadMessage(mailboxId string, msg *models.MessageCreateModel) (*models.MessageUidListItem, error) {
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
	return &models.MessageUidListItem{
		Id:  msgData.Id,
		Uid: uid,
	}, nil
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

func (cs *CentaureissiService) IndexSearchDocument(msg *models.MessageIndexModel) error {
	doc := &search.CentaureissiSearchDocument{
		Id:        msg.Id,
		Hash:      msg.Hash,
		UserId:    msg.UserId,
		MailboxId: msg.MailboxId,
		Sender:    msg.Sender,
		From:      msg.From,
		To:        msg.To,
		Cc:        msg.Cc,
		Bcc:       msg.Bcc,
		Subject:   msg.Subject,
		Date:      msg.Date,
		Content:   msg.Content,
	}
	err := cs.search.Index(*doc)
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

func (cs *CentaureissiService) SearchMessages(userId string, q string, page int, perPage int) (*models.SearchResponse, error) {

	result, err := cs.search.Search(userId, q, page, perPage)
	if err != nil {
		return nil, err
	}
	hits := make([]*models.SearchResponseItem, 0)
	hitCount := result.Total

	totalPages := (hitCount / uint64(perPage))
	if hitCount%uint64(perPage) > 0 {
		totalPages++
	}

	lowerBound := uint64((page - 1) * perPage)
	upperBound := uint64(page * perPage)
	if lowerBound <= hitCount {
		if upperBound > hitCount {
			upperBound = hitCount
		}
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
			}
			hits = append(hits, item)
		}
	}

	response := &models.SearchResponse{
		Hits:       hits,
		Page:       page,
		TotalPages: int(totalPages),
		Count:      result.Total,
	}
	return response, nil
}
