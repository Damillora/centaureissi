package imapinterface

import (
	"bytes"
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
	return unimplemented()
}

// Namespace implements imapserver.SessionIMAP4rev2.
func (c *CentaureissiImapSession) Namespace() (*imap.NamespaceData, error) {
	return nil, unimplemented()
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
	hash, err := c.services.UploadMessageContent(buf.Bytes())
	if err != nil {
		return nil, err
	}

	msg := &models.MessageCreateModel{
		Hash:  hash,
		Size:  uint64(buf.Len()),
		Flags: make(map[string]bool),
	}
	for _, flag := range options.Flags {
		msg.Flags[string(canonicalFlag(flag))] = true
	}

	msgData, err := c.services.UploadMessage(mbox.ID, msg)
	if err != nil {
		return nil, err
	}

	// Search Index
	hydrated := c.hydrateMessage(msgData)
	indexDoc := hydrated.createSearchDocument()

	err = c.services.IndexSearchDocument(*indexDoc)
	if err != nil {
		return nil, err
	}

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
	return nil, unimplemented()
}

// Create implements imapserver.Session.
func (c *CentaureissiImapSession) Create(mailbox string, options *imap.CreateOptions) error {

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
	return unimplemented()
}

// Expunge implements imapserver.Session.
func (c *CentaureissiImapSession) Expunge(w *imapserver.ExpungeWriter, uids *imap.UIDSet) error {
	return unimplemented()
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
	mboxTracker := c.tracker.TrackMailbox(c.user.ID, c.mailbox.mailboxSchema.ID)

	var err error
	c.mailbox.forEach(numSet, func(seqNum uint32, msg *schema.Message) {
		if err != nil {
			return
		}

		if markSeen {
			mboxTracker.QueueMessageFlags(seqNum, imap.UID(msg.Uid), flagList(msg), nil)
		}

		respWriter := w.CreateMessage(c.mailbox.mailboxSession.EncodeSeqNum(seqNum))

		heavyweightMsg := c.hydrateMessage(msg)
		if heavyweightMsg == nil {
			return
		}
		err = heavyweightMsg.fetch(respWriter, options)
	})
	return err
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
	return unimplemented()
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
	mboxTracker := c.tracker.TrackMailbox(c.user.ID, mbox.ID)
	c.mailbox = &CentaureissiImapMailbox{
		services:       c.services,
		mailboxSchema:  mbox,
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

	data := imap.StatusData{Mailbox: mbox.Name}
	if options.NumMessages {
		num := uint32(0)
		data.NumMessages = &num
	}
	if options.UIDNext {
		uid, _ := c.services.CounterMailboxUid(mbox.ID)
		data.UIDNext = imap.UID(uid + 1)
	}
	if options.UIDValidity {
		data.UIDValidity = mbox.UidValidity
	}
	if options.NumUnseen {
		num := uint32(0)
		data.NumUnseen = &num
	}
	if options.NumDeleted {
		num := uint32(0)
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
	return unimplemented()
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
	return unimplemented()
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
