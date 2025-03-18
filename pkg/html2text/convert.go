package html2text

import (
	"bytes"
	"fmt"
	"regexp"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type htmlConvert struct {
	htmlDoc string
	buf     *bytes.Buffer
}

func (hc *htmlConvert) startParsing() error {
	hRead := bytes.NewReader([]byte(hc.htmlDoc))
	node, err := html.Parse(hRead)
	if err != nil {
		return err
	}

	hc.traverse(node, nil)

	return nil
}
func (hc *htmlConvert) traverse(node *html.Node, parent *html.Node) {
	switch node.Type {
	case html.ElementNode:
		switch node.DataAtom {
		// Text
		case atom.P:
			if node.PrevSibling != nil && isPrintNewlineForBlocks(node) {
				fmt.Fprintln(hc.buf)
				fmt.Fprintln(hc.buf)
			}
			// traverse the child nodes
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				hc.traverse(c, node)
			}
		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
			if node.PrevSibling != nil && isPrintNewlineForBlocks(node) {
				fmt.Fprintln(hc.buf)
				fmt.Fprintln(hc.buf)
			}
			// traverse the child nodes
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				hc.traverse(c, node)
			}
		case atom.Br:
			fmt.Fprintln(hc.buf)
		case atom.Ul:
			fmt.Fprintln(hc.buf)
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				fmt.Fprintf(hc.buf, " - ")
				hc.traverse(c, node)
				fmt.Fprintln(hc.buf)
			}
		case atom.Ol:
			fmt.Fprintln(hc.buf)
			i := 1
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				fmt.Fprintf(hc.buf, "%d. ", i)
				hc.traverse(c, node)
				fmt.Fprintln(hc.buf)
				i++
			}
		case atom.Head:
			return
		default:
			// traverse the child nodes
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				hc.traverse(c, node)
			}
		}
	case html.TextNode:
		if node.PrevSibling != nil && isPrintNewlineForBlocks(node) {
			fmt.Fprintln(hc.buf)
			fmt.Fprintln(hc.buf)
		}
		reWhitespaceNorm := regexp.MustCompile(`[\n]+`)
		nodeData := reWhitespaceNorm.ReplaceAll([]byte(node.Data), []byte(" "))
		fmt.Fprintf(hc.buf, "%s", nodeData)
	default:
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			hc.traverse(c, node)
		}
	}
}
