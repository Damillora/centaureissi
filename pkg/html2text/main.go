package html2text

import (
	"bytes"
)

func NewConverter(htmlDoc string) *htmlConvert {
	buf := bytes.NewBufferString("")
	htmlConvert := &htmlConvert{
		htmlDoc: htmlDoc,
		buf:     buf,
	}
	return htmlConvert
}
func Parse(htmlDoc string) string {
	htmlConvert := NewConverter(htmlDoc)

	err := htmlConvert.startParsing()
	if err != nil {
		return ""
	}

	result := htmlConvert.buf.String()
	return result
}
