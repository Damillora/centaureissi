package imapinterface

import (
	"github.com/Damillora/centaureissi/pkg/models"
)

func (c *CentaureissiImapSession) hydrateMessage(msg *models.MessageUidListItem) *CentaureissiMessage {
	message, err := c.services.GetMessageById(msg.Id)
	if err != nil {
		return nil
	}

	blobs, err := c.services.GetMessageContent(message.Hash)
	if err != nil {
		return nil
	}
	return &CentaureissiMessage{
		Message: message,
		buf:     blobs,
	}
}

// func (c *Cen)
