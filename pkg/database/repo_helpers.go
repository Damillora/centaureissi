package database

import (
	"fmt"
	"strconv"
)

func stringToUint32(num string) uint32 {
	convertedUint, _ := strconv.ParseUint(string(num), 10, 32)
	return uint32(convertedUint)
}
func uint32ToString(num uint32) string {
	return strconv.FormatUint(uint64(num), 10)
}

func formatUserIdAndName(userId string, mboxName string) string {
	return fmt.Sprintf("%s:%s", userId, mboxName)
}

func formatMailboxIdAndUid(mailboxId string, uid uint32) string {
	uidString := uint32ToString(uid)
	return fmt.Sprintf("%s:%s", mailboxId, uidString)
}
