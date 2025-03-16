package imapinterface

import (
	"strings"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/emersion/go-imap/v2"
)

func flagList(msg *schema.Message) []imap.Flag {
	var flags []imap.Flag
	for flag := range msg.Flags {
		flags = append(flags, imap.Flag(flag))
	}
	return flags
}

func canonicalFlag(flag imap.Flag) imap.Flag {
	return imap.Flag(strings.ToLower(string(flag)))
}
