package imapinterface

import (
	"fmt"
	"sync"

	"github.com/Damillora/centaureissi/pkg/services"
	"github.com/emersion/go-imap/v2/imapserver"
)

type CentaureissiMailbox struct {
	*imapserver.MailboxTracker
}

type CentaureissiMailboxTracker struct {
	services *services.CentaureissiService

	mutex    sync.Mutex
	trackers map[string]*CentaureissiMailbox
}

func NewCentaureissiMailboxTracker(cs *services.CentaureissiService) *CentaureissiMailboxTracker {
	return &CentaureissiMailboxTracker{
		services: cs,
		trackers: make(map[string]*CentaureissiMailbox),
	}
}

func (cmt *CentaureissiMailboxTracker) TrackMailbox(userId string, mboxId string) *CentaureissiMailbox {
	cmt.mutex.Lock()
	defer cmt.mutex.Unlock()

	countMsgs, _ := cmt.services.CounterMessagesInMailbox(mboxId)

	trackId := fmt.Sprintf("%s:%s", userId, mboxId)
	if cmt.trackers[trackId] == nil {
		cmt.trackers[trackId] = &CentaureissiMailbox{
			MailboxTracker: imapserver.NewMailboxTracker(countMsgs),
		}
	}
	return cmt.trackers[trackId]
}
