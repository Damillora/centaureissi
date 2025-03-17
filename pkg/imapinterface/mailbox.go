package imapinterface

import (
	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/services"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
)

type CentaureissiImapMailbox struct {
	services *services.CentaureissiService

	searchRes      imap.UIDSet
	mailboxSchema  *schema.Mailbox
	mailboxSession *imapserver.SessionTracker
}

func (cim *CentaureissiImapMailbox) GetSelectInfo() *imap.SelectData {
	flags := []imap.Flag{}

	permanentFlags := make([]imap.Flag, len(flags))
	copy(permanentFlags, flags)
	permanentFlags = append(permanentFlags, imap.FlagWildcard)

	currentUID, err := cim.services.CounterMailboxUid(cim.mailboxSchema.Id)
	if err != nil {
		return nil
	}

	return &imap.SelectData{
		Flags:          flags,
		PermanentFlags: permanentFlags,
		NumMessages:    uint32(0),
		UIDNext:        imap.UID(currentUID + 1),
		UIDValidity:    cim.mailboxSchema.UidValidity,
	}
}

func (cim *(CentaureissiImapMailbox)) forEach(numSet imap.NumSet, f func(seqNum uint32, msg *schema.Message)) {
	// TODO: optimize

	numSet = cim.staticNumSet(numSet)

	messages, err := cim.services.ListMessageByMailboxId(cim.mailboxSchema.Id)
	if err != nil {

	}
	for i, msg := range messages {
		seqNum := uint32(i) + 1

		var contains bool
		switch numSet := numSet.(type) {
		case imap.SeqSet:
			seqNum := cim.mailboxSession.EncodeSeqNum(seqNum)
			contains = seqNum != 0 && numSet.Contains(seqNum)
		case imap.UIDSet:
			contains = numSet.Contains(imap.UID(msg.Uid))
		}
		if !contains {
			continue
		}

		f(seqNum, msg)
	}
}

// staticNumSet converts a dynamic sequence set into a static one.
//
// This is necessary to properly handle the special symbol "*", which
// represents the maximum sequence number or UID in the mailbox.
//
// This function also handles the special SEARCHRES marker "$".
func (cim *CentaureissiImapMailbox) staticNumSet(numSet imap.NumSet) imap.NumSet {
	if imap.IsSearchRes(numSet) {
		return cim.searchRes
	}

	switch numSet := numSet.(type) {
	case imap.SeqSet:
		maxSeqSet, _ := cim.services.CounterMessagesInMailbox(cim.mailboxSchema.Id)
		max := uint32(maxSeqSet)
		for i := range numSet {
			r := &numSet[i]
			staticNumRange(&r.Start, &r.Stop, max)
		}
	case imap.UIDSet:
		maxIdSet, _ := cim.services.CounterMailboxUidNext(cim.mailboxSchema.Id)
		max := uint32(maxIdSet) - 1
		for i := range numSet {
			r := &numSet[i]
			staticNumRange((*uint32)(&r.Start), (*uint32)(&r.Stop), max)
		}
	}

	return numSet
}

func staticNumRange(start, stop *uint32, max uint32) {
	dyn := false
	if *start == 0 {
		*start = max
		dyn = true
	}
	if *stop == 0 {
		*stop = max
		dyn = true
	}
	if dyn && *start > *stop {
		*start, *stop = *stop, *start
	}
}
