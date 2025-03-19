package imapinterface

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/models"
	"github.com/Damillora/centaureissi/pkg/services"
	imap "github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
)

func unimplemented() error {
	return &imap.Error{
		Type: imap.StatusResponseTypeNo,
		Code: imap.ResponseCodeCannot,
		Text: "in development!",
	}
}

type CentaureissiImapSession struct {
	services *services.CentaureissiService
	tracker  *CentaureissiMailboxTracker

	mutex   sync.Mutex
	user    *schema.User
	mailbox *CentaureissiImapMailbox
}

func NewCentaureissiImapSession(s *services.CentaureissiService, t *CentaureissiMailboxTracker) *CentaureissiImapSession {
	return &CentaureissiImapSession{
		services: s,
		tracker:  t,
	}
}

// Move implements imapserver.SessionIMAP4rev2.
func (c *CentaureissiImapSession) Move(w *imapserver.MoveWriter, numSet imap.NumSet, dest string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mbox, err := c.services.GetMailboxByUserIdAndName(c.user.ID, dest)
	if err != nil {
		return err
	}
	if mbox == nil {
		return &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Code: imap.ResponseCodeNonExistent,
			Text: "Mailbox does not exist!",
		}
	} else if c.mailbox != nil && c.mailbox.mailboxSchema.Id == mbox.Id {
		return &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Text: "Source and destination mailboxes are identical",
		}
	}

	msgs, err := c.services.ListMessageUidsByMailboxId(c.mailbox.mailboxSchema.Id)
	if err != nil {
		return err
	}

	expunged := make([]string, 0)

	var sourceUIDs, destUIDs imap.UIDSet
	c.mailbox.forEach(msgs, numSet, func(seqNum uint32, msgItem *models.MessageUidListItem) {
		msg := c.hydrateMessage(msgItem)
		if msg == nil {
			return
		}

		appendData, err := c.appendMsg(mbox, msg.buf, &imap.AppendOptions{
			Time:  msg.CreatedAt,
			Flags: flagList(msg.Message),
		})
		if err != nil {
			return
		}
		sourceUIDs.AddNum(imap.UID(msgItem.Uid))
		destUIDs.AddNum(appendData.UID)

		expunged = append(expunged, msgItem.Id)
	})

	seqNums := c.mailbox.expunge(msgs, expunged)

	err = w.WriteCopyData(&imap.CopyData{
		UIDValidity: mbox.UidValidity,
		SourceUIDs:  sourceUIDs,
		DestUIDs:    destUIDs,
	})
	if err != nil {
		return err
	}

	for _, seqNum := range seqNums {
		if err := w.WriteExpunge(c.mailbox.mailboxSession.EncodeSeqNum(seqNum)); err != nil {
			return err
		}
	}

	return nil
}

// Namespace implements imapserver.SessionIMAP4rev2.
func (c *CentaureissiImapSession) Namespace() (*imap.NamespaceData, error) {
	return &imap.NamespaceData{
		Personal: []imap.NamespaceDescriptor{{Delim: mailboxDelim}},
	}, nil
}

// Append implements imapserver.Session.
func (c *CentaureissiImapSession) Append(mailbox string, r imap.LiteralReader, options *imap.AppendOptions) (*imap.AppendData, error) {
	mbox, err := c.services.GetMailboxByUserIdAndName(c.user.ID, mailbox)
	if err != nil {
		return nil, err
	}
	if mbox == nil {
		return nil, &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Code: imap.ResponseCodeNonExistent,
			Text: "Mailbox does not exist!",
		}
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}

	return c.appendMsg(mbox, buf.Bytes(), options)
}

func (c *CentaureissiImapSession) appendMsg(mbox *schema.Mailbox, buf []byte, options *imap.AppendOptions) (*imap.AppendData, error) {

	hash, err := c.services.UploadMessageContent(buf)
	if err != nil {
		return nil, err
	}

	msg := &models.MessageCreateModel{
		Hash:  hash,
		Size:  uint64(len(buf)),
		Flags: make(map[string]bool),
	}
	for _, flag := range options.Flags {
		msg.Flags[string(canonicalFlag(flag))] = true
	}

	msgData, err := c.services.UploadMessage(mbox.Id, msg)
	if err != nil {
		return nil, err
	}

	// Search Index
	hydrated := c.hydrateMessage(msgData)
	indexDoc := hydrated.createSearchDocument()
	indexDoc.MailboxId = mbox.Id

	err = c.services.IndexSearchDocument(indexDoc)
	if err != nil {
		return nil, err
	}

	messageCount, err := c.services.CounterMessagesInMailbox(mbox.Id)
	if err != nil {
		return nil, err
	}
	mboxTracker := c.tracker.TrackMailbox(c.user.ID, mbox.Id)
	mboxTracker.QueueNumMessages(uint32(messageCount))

	appendData := &imap.AppendData{
		UID:         imap.UID(msgData.Uid),
		UIDValidity: mbox.UidValidity,
	}
	return appendData, nil
}
func (c *CentaureissiImapSession) Close() error {
	// We defer closing
	return nil
}

// Copy implements imapserver.Session.
func (c *CentaureissiImapSession) Copy(numSet imap.NumSet, dest string) (*imap.CopyData, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mbox, err := c.services.GetMailboxByUserIdAndName(c.user.ID, dest)
	if err != nil {
		return nil, err
	}
	if mbox == nil {
		return nil, &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Code: imap.ResponseCodeNonExistent,
			Text: "Mailbox does not exist!",
		}
	} else if c.mailbox != nil && c.mailbox.mailboxSchema.Id == mbox.Id {
		return nil, &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Text: "Source and destination mailboxes are identical",
		}
	}

	msgs, err := c.services.ListMessageUidsByMailboxId(c.mailbox.mailboxSchema.Id)
	if err != nil {
		return nil, err
	}

	var sourceUIDs, destUIDs imap.UIDSet
	c.mailbox.forEach(msgs, numSet, func(seqNum uint32, msgItem *models.MessageUidListItem) {
		msg := c.hydrateMessage(msgItem)
		if msg == nil {
			return
		}

		appendData, err := c.appendMsg(mbox, msg.buf, &imap.AppendOptions{
			Time:  msg.CreatedAt,
			Flags: flagList(msg.Message),
		})
		if err != nil {
			return
		}
		sourceUIDs.AddNum(imap.UID(msg.Uid))
		destUIDs.AddNum(appendData.UID)
	})

	return &imap.CopyData{
		UIDValidity: mbox.UidValidity,
		SourceUIDs:  sourceUIDs,
		DestUIDs:    destUIDs,
	}, nil
}

// Create implements imapserver.Session.
func (c *CentaureissiImapSession) Create(mailbox string, options *imap.CreateOptions) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	name := strings.TrimRight(mailbox, string(mailboxDelim))

	mbox, err := c.services.GetMailboxByUserIdAndName(c.user.ID, name)
	if err != nil {
		return err
	}
	if mbox != nil {
		return &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Code: imap.ResponseCodeAlreadyExists,
			Text: "Mailbox already exists",
		}
	}

	c.services.CreateMailbox(c.user.ID, name)
	return nil
}

// Delete implements imapserver.Session.
func (c *CentaureissiImapSession) Delete(mailbox string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.services.DeleteMailbox(c.user.ID, mailbox)
	if err != nil {
		return err
	}
	return nil
}

// Expunge implements imapserver.Session.
func (c *CentaureissiImapSession) Expunge(w *imapserver.ExpungeWriter, uids *imap.UIDSet) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expunged := make([]string, 0)
	msgs, err := c.services.ListMessageUidsByMailboxId(c.mailbox.mailboxSchema.Id)
	if err != nil {
		return err
	}
	for _, msgItem := range msgs {
		msg, err := c.services.GetMessageById(msgItem.Id)
		if err != nil {
			return err
		}

		if uids != nil && !uids.Contains(imap.UID(msg.Uid)) {
			continue
		}
		if _, ok := msg.Flags[string(canonicalFlag(imap.FlagDeleted))]; ok {
			expunged = append(expunged, msg.Id)
		}
	}

	if len(expunged) == 0 {
		return nil
	}
	c.mailbox.expunge(msgs, expunged)
	return nil
}

// Fetch implements imapserver.Session.
func (c *CentaureissiImapSession) Fetch(w *imapserver.FetchWriter, numSet imap.NumSet, options *imap.FetchOptions) error {
	markSeen := false
	for _, bs := range options.BodySection {
		if !bs.Peek {
			markSeen = true
			break
		}
	}
	mboxTracker := c.tracker.TrackMailbox(c.user.ID, c.mailbox.mailboxSchema.Id)

	msgs, err := c.services.ListMessageUidsByMailboxId(c.mailbox.mailboxSchema.Id)
	if err != nil {
		return err
	}

	c.mailbox.forEach(msgs, numSet, func(seqNum uint32, msgItem *models.MessageUidListItem) {
		msg := c.hydrateMessage(msgItem)
		if msg == nil {
			return
		}

		if markSeen {
			msg.Flags[string(canonicalFlag(imap.FlagSeen))] = true
			c.services.UpdateMessageFlags(msg.Id, &models.MessageUpdateFlagsModel{
				Flags: msg.Flags,
			})
			mboxTracker.QueueMessageFlags(seqNum, imap.UID(msgItem.Uid), flagList(msg.Message), nil)
		}

		respWriter := w.CreateMessage(c.mailbox.mailboxSession.EncodeSeqNum(seqNum))

		err = msg.fetch(respWriter, options)
		if err != nil {
			return
		}
	})
	return nil
}

// Idle implements imapserver.Session.
func (c *CentaureissiImapSession) Idle(w *imapserver.UpdateWriter, stop <-chan struct{}) error {
	if c.mailbox != nil {
		return c.mailbox.mailboxSession.Idle(w, stop)
	}
	return nil
}

// List implements imapserver.Session.
func (c *CentaureissiImapSession) List(w *imapserver.ListWriter, ref string, patterns []string, options *imap.ListOptions) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// TODO: fail if ref doesn't exist

	if len(patterns) == 0 {
		return w.WriteList(&imap.ListData{
			Attrs: []imap.MailboxAttr{imap.MailboxAttrNoSelect},
			Delim: mailboxDelim,
		})
	}

	mailboxes, err := c.services.ListMailboxesByUserId(c.user.ID)

	if err != nil {
		return &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Text: "Error getting mailboxes",
		}
	}
	var imapListData []imap.ListData
	for _, mbox := range mailboxes {
		match := false
		for _, pattern := range patterns {
			match = imapserver.MatchList(mbox.Name, mailboxDelim, ref, pattern)
			if match {
				break
			}
		}
		if !match {
			continue
		}

		data := c.createListData(mbox, options)
		if data != nil {
			imapListData = append(imapListData, *data)
		}
	}

	sort.Slice(imapListData, func(i, j int) bool {
		return imapListData[i].Mailbox < imapListData[j].Mailbox
	})

	for _, data := range imapListData {
		if err := w.WriteList(&data); err != nil {
			return err
		}
	}
	return nil
}

func (c *CentaureissiImapSession) createListData(mbox *schema.Mailbox, options *imap.ListOptions) *imap.ListData {
	if options.SelectSubscribed && !mbox.Subscribed {
		return nil
	}
	data := &imap.ListData{
		Mailbox: mbox.Name,
		Delim:   mailboxDelim,
	}

	return data
}

// Login implements imapserver.Session.
func (c *CentaureissiImapSession) Login(username string, password string) error {
	c.user = c.services.Login(username, password)

	if c.user == nil {
		return imapserver.ErrAuthFailed
	}
	return nil
}

// Poll implements imapserver.Session.
func (c *CentaureissiImapSession) Poll(w *imapserver.UpdateWriter, allowExpunge bool) error {
	if c.mailbox != nil {
		return c.mailbox.mailboxSession.Poll(w, allowExpunge)
	}
	return nil
}

// Rename implements imapserver.Session.
func (c *CentaureissiImapSession) Rename(mailbox string, newName string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, err := c.services.GetMailboxByUserIdAndName(c.user.ID, mailbox)
	if err != nil {
		return err
	}
	c.services.UpdateMailboxName(c.user.ID, mailbox, newName)

	return nil
}

// Search implements imapserver.Session.
func (c *CentaureissiImapSession) Search(kind imapserver.NumKind, criteria *imap.SearchCriteria, options *imap.SearchOptions) (*imap.SearchData, error) {
	return nil, unimplemented()
}

// Select implements imapserver.Session.
func (c *CentaureissiImapSession) Select(mailbox string, options *imap.SelectOptions) (*imap.SelectData, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mbox, err := c.services.GetMailboxByUserIdAndName(c.user.ID, mailbox)
	if err != nil {
		return nil, err
	}
	if mbox == nil {
		// name := strings.TrimRight(mailbox, string(mailboxDelim))
		// mbox, err = c.services.CreateMailbox(c.user.ID, name)

		return nil, &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Code: imap.ResponseCodeNonExistent,
			Text: "Mailbox does not exist!",
		}
	}
	mboxTracker := c.tracker.TrackMailbox(c.user.ID, mbox.Id)
	c.mailbox = &CentaureissiImapMailbox{
		services:       c.services,
		mailboxSchema:  mbox,
		mailboxTracker: mboxTracker.MailboxTracker,
		mailboxSession: mboxTracker.NewSession(),
	}
	selectInfo := c.mailbox.GetSelectInfo()
	return selectInfo, nil
}

// Status implements imapserver.Session.
func (c *CentaureissiImapSession) Status(mailbox string, options *imap.StatusOptions) (*imap.StatusData, error) {
	mbox, err := c.services.GetMailboxByUserIdAndName(c.user.ID, mailbox)
	if err != nil {
		return nil, err
	}
	if mbox == nil {
		return nil, &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Code: imap.ResponseCodeNonExistent,
			Text: "Mailbox does not exist!",
		}
	}

	numMessages, err := c.services.CounterMessagesInMailbox(mbox.Id)
	if err != nil {
		return nil, err
	}
	numRead, err := c.services.CounterMessagesInMailboxByFlag(mbox.Id, string(canonicalFlag(imap.FlagSeen)))
	if err != nil {
		return nil, err
	}
	numDeleted, err := c.services.CounterMessagesInMailboxByFlag(mbox.Id, string(canonicalFlag(imap.FlagDeleted)))
	if err != nil {
		return nil, err
	}

	data := imap.StatusData{Mailbox: mbox.Name}
	if options.NumMessages {
		num := numMessages
		data.NumMessages = &num
	}
	if options.UIDNext {
		uid, _ := c.services.CounterMailboxUid(mbox.Id)
		data.UIDNext = imap.UID(uid + 1)
	}
	if options.UIDValidity {
		data.UIDValidity = mbox.UidValidity
	}
	if options.NumUnseen {
		num := numMessages - numRead
		data.NumUnseen = &num
	}
	if options.NumDeleted {
		num := numDeleted
		data.NumDeleted = &num
	}
	if options.Size {
		size := int64(0)
		data.Size = &size
	}
	return &data, nil
}

// Store implements imapserver.Session.
func (c *CentaureissiImapSession) Store(w *imapserver.FetchWriter, numSet imap.NumSet, flags *imap.StoreFlags, options *imap.StoreOptions) error {
	msgs, err := c.services.ListMessageUidsByMailboxId(c.mailbox.mailboxSchema.Id)
	if err != nil {
		return err
	}

	mboxTracker := c.tracker.TrackMailbox(c.user.ID, c.mailbox.mailboxSchema.Id)
	c.mailbox.forEach(msgs, numSet, func(seqNum uint32, msgItem *models.MessageUidListItem) {
		msg, err := c.services.GetMessageById(msgItem.Id)
		if err != nil {
			return
		}
		if msg.Flags == nil {
			msg.Flags = make(map[string]bool)
		}

		switch flags.Op {
		case imap.StoreFlagsSet:
			msg.Flags = make(map[string]bool)
			fallthrough
		case imap.StoreFlagsAdd:
			for _, flag := range flags.Flags {
				msg.Flags[string(canonicalFlag(flag))] = true
			}
		case imap.StoreFlagsDel:
			for _, flag := range flags.Flags {
				delete(msg.Flags, string(canonicalFlag(flag)))
			}
		default:
			panic(fmt.Errorf("unknown STORE flag operation: %v", flags.Op))
		}

		c.services.UpdateMessageFlags(msg.Id, &models.MessageUpdateFlagsModel{
			Flags: msg.Flags,
		})
		mboxTracker.QueueMessageFlags(seqNum, imap.UID(msg.Uid), flagList(msg), nil)
	})
	if !flags.Silent {
		return c.Fetch(w, numSet, &imap.FetchOptions{Flags: true})
	}
	return nil
}

// Subscribe implements imapserver.Session.
func (c *CentaureissiImapSession) Subscribe(mailbox string) error {
	if c.user != nil {
		_, err := c.services.UpdateMailboxSubscribeStatus(c.user.ID, mailbox, true)
		if err != nil {
			return err
		}
		return nil
	}
	return &imap.Error{
		Type: imap.StatusResponseTypeNo,
		Text: "not yet logged in?",
	}
}

// Unselect implements imapserver.Session.
func (c *CentaureissiImapSession) Unselect() error {
	if c.mailbox == nil {
		return &imap.Error{
			Type: imap.StatusResponseTypeNo,
			Text: "Mailbox unselected!",
		}
	}
	c.mailbox.mailboxSession.Close()
	c.mailbox = nil
	return nil
}

// Unsubscribe implements imapserver.Session.
func (c *CentaureissiImapSession) Unsubscribe(mailbox string) error {
	if c.user != nil {
		_, err := c.services.UpdateMailboxSubscribeStatus(c.user.ID, mailbox, false)
		if err != nil {
			return err
		}
		return nil
	}
	return &imap.Error{
		Type: imap.StatusResponseTypeNo,
		Text: "not yet logged in?",
	}
}

var _ imapserver.SessionIMAP4rev2 = (*CentaureissiImapSession)(nil)

// Close() error

// // Not authenticated state
// Login(username, password string) error

// // Authenticated state
// Select(mailbox string, options *imap.SelectOptions) (*imap.SelectData, error)
// Create(mailbox string, options *imap.CreateOptions) error
// Delete(mailbox string) error
// Rename(mailbox, newName string) error
// Subscribe(mailbox string) error
// Unsubscribe(mailbox string) error
// List(w *ListWriter, ref string, patterns []string, options *imap.ListOptions) error
// Status(mailbox string, options *imap.StatusOptions) (*imap.StatusData, error)
// Append(mailbox string, r imap.LiteralReader, options *imap.AppendOptions) (*imap.AppendData, error)
// Poll(w *UpdateWriter, allowExpunge bool) error
// Idle(w *UpdateWriter, stop <-chan struct{}) error

// // Selected state
// Unselect() error
// Expunge(w *ExpungeWriter, uids *imap.UIDSet) error
// Search(kind NumKind, criteria *imap.SearchCriteria, options *imap.SearchOptions) (*imap.SearchData, error)
// Fetch(w *FetchWriter, numSet imap.NumSet, options *imap.FetchOptions) error
// Store(w *FetchWriter, numSet imap.NumSet, flags *imap.StoreFlags, options *imap.StoreOptions) error
// Copy(numSet imap.NumSet, dest string) (*imap.CopyData, error)
