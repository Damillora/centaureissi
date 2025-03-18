package html2text

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func isBlockAtom(nodeAtom atom.Atom) bool {
	return nodeAtom == atom.P ||
		nodeAtom == atom.H1 ||
		nodeAtom == atom.H2 ||
		nodeAtom == atom.H3 ||
		nodeAtom == atom.H4 ||
		nodeAtom == atom.H5 ||
		nodeAtom == atom.H6
}

func isPrintNewlineForBlocks(node *html.Node) bool {
	return node.PrevSibling.Type == html.TextNode ||
		isBlockAtom(node.PrevSibling.DataAtom)
}
