package imapinterface

import (
	"fmt"
	"sync"

	"github.com/emersion/go-imap/v2/imapserver"
)

type CentaureissiMailbox struct {
	*imapserver.MailboxTracker
}

type CentaureissiMailboxTracker struct {
	mutex    sync.Mutex
	trackers map[string]*CentaureissiMailbox
}

func NewCentaureissiMailboxTracker() *CentaureissiMailboxTracker {
	return &CentaureissiMailboxTracker{
		trackers: make(map[string]*CentaureissiMailbox),
	}
}

func (cmt *CentaureissiMailboxTracker) TrackMailbox(userId string, mboxId string) *CentaureissiMailbox {
	cmt.mutex.Lock()
	defer cmt.mutex.Unlock()

	trackId := fmt.Sprintf("%s:%s", userId, mboxId)
	if cmt.trackers[trackId] == nil {
		cmt.trackers[trackId] = &CentaureissiMailbox{
			MailboxTracker: imapserver.NewMailboxTracker(0),
		}
	}
	return cmt.trackers[trackId]
}
