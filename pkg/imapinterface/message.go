package imapinterface

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/search"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-message/textproto"
)

type CentaureissiMessage struct {
	*schema.Message
	buf []byte
}

func (msg *CentaureissiMessage) fetch(w *imapserver.FetchResponseWriter, options *imap.FetchOptions) error {
	w.WriteUID(imap.UID(msg.Uid))

	if options.Flags {
		w.WriteFlags(flagList(msg.Message))
	}
	if options.InternalDate {
		w.WriteInternalDate(msg.CreatedAt)
	}
	if options.RFC822Size {
		w.WriteRFC822Size(int64(msg.Size))
	}
	if options.Envelope {
		w.WriteEnvelope(msg.envelope())
	}
	if options.BodyStructure != nil {
		w.WriteBodyStructure(imapserver.ExtractBodyStructure(bytes.NewReader(msg.buf)))
	}

	for _, bs := range options.BodySection {
		buf := imapserver.ExtractBodySection(bytes.NewReader(msg.buf), bs)
		wc := w.WriteBodySection(bs, int64(len(buf)))
		_, writeErr := wc.Write(buf)
		closeErr := wc.Close()
		if writeErr != nil {
			return writeErr
		}
		if closeErr != nil {
			return closeErr
		}
	}

	for _, bs := range options.BinarySection {
		buf := imapserver.ExtractBinarySection(bytes.NewReader(msg.buf), bs)
		wc := w.WriteBinarySection(bs, int64(len(buf)))
		_, writeErr := wc.Write(buf)
		closeErr := wc.Close()
		if writeErr != nil {
			return writeErr
		}
		if closeErr != nil {
			return closeErr
		}
	}

	for _, bss := range options.BinarySectionSize {
		n := imapserver.ExtractBinarySectionSize(bytes.NewReader(msg.buf), bss)
		w.WriteBinarySectionSize(bss, n)
	}

	return w.Close()
}

func (msg *CentaureissiMessage) envelope() *imap.Envelope {
	br := bufio.NewReader(bytes.NewReader(msg.buf))
	header, err := textproto.ReadHeader(br)
	if err != nil {
		return nil
	}
	return imapserver.ExtractEnvelope(header)
}

func (msg *CentaureissiMessage) createSearchDocument() *search.CentaureissiSearchDocument {
	envelope := msg.envelope()
	var sender []string
	for _, senderAddr := range envelope.Sender {
		sender = append(sender, senderAddr.Addr())
	}
	var from []string
	for _, fromAddr := range envelope.From {
		from = append(from, fromAddr.Addr())
	}
	var to []string
	for _, toAddr := range envelope.To {
		to = append(to, toAddr.Addr())
	}
	var cc []string
	for _, ccAddr := range envelope.To {
		cc = append(cc, ccAddr.Addr())
	}
	var bcc []string
	for _, bccAddr := range envelope.To {
		bcc = append(bcc, bccAddr.Addr())
	}

	contents := ""

	readMail, err := mail.CreateReader(bytes.NewReader(msg.buf))
	if err != nil {
		return nil
	}
	for {
		p, err := readMail.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			b, _ := io.ReadAll(p.Body)
			contents += string(b)
			contents += "\n"
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			log.Printf("Got attachment: %v\n", filename)
		}
	}

	return &search.CentaureissiSearchDocument{
		Hash:    msg.Hash,
		Sender:  strings.Join(sender[:], ", "),
		From:    strings.Join(from[:], ", "),
		To:      strings.Join(to[:], ", "),
		Cc:      strings.Join(cc[:], ", "),
		Bcc:     strings.Join(bcc[:], ", "),
		Subject: envelope.Subject,
		Date:    envelope.Date,
		Content: contents,
	}
}
