package imapinterface

import "github.com/Damillora/centaureissi/pkg/database/schema"

func (c *CentaureissiImapSession) hydrateMessage(msg *schema.Message) *CentaureissiMessage {
	blobs, err := c.services.GetMessageContent(msg.Hash)
	if err != nil {
		return nil
	}
	return &CentaureissiMessage{
		Message: msg,
		buf:     blobs,
	}
}

// func (c *Cen)
