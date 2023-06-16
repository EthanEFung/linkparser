package linkparser

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type LinkParser interface {
	UseNode(node *html.Node) error
	UseReader(reader io.Reader) error
	UseHTMLFile(path string) error
	Parse() ([]Link, error)
	Marshal() ([]byte, error)
}

func NewLinkParser() LinkParser {
	return &linkParser{}
}

type linkParser struct {
	node *html.Node
}

func (lp *linkParser) UseHTMLFile(path string) error {
	b, err := os.Open(path)
	if err != nil {
		return err
	}
	defer b.Close()

	node, err := html.Parse(b)
	if err != nil {
		return err
	}
	lp.node = node
	return nil
}

func (lp *linkParser) UseReader(r io.Reader) error {
	node, err := html.Parse(r)
	if err != nil {
		return err
	}
	lp.node = node
	return nil
}

func (lp *linkParser) UseNode(node *html.Node) error {
	lp.node = node
	return nil
}

func (lp *linkParser) Parse() ([]Link, error) {
	node := lp.node
	links := []Link{}

	if node == nil {
		return links, errors.New("html node must be set before running this function")
	}

	if node.DataAtom == atom.A {

		link := Link{}

		for _, attr := range lp.node.Attr {
			if attr.Key == "href" {
				link.Href = attr.Val
			}
		}

		b := strings.Builder{}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			b.WriteString(getText(child))
		}

		link.Text = b.String()
		links = append(links, link)

	} else if node.Type == html.ElementNode || node.Type == html.DocumentNode {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			cParser := NewLinkParser()
			cParser.UseNode(child)
			cLinks, err := cParser.Parse()
			if err != nil {
				return links, err
			}
			links = append(links, cLinks...)
		}

	}

	return links, nil
}

func (lp *linkParser) Marshal() ([]byte, error) {
	links, err := lp.Parse()
	if err != nil {
		return []byte{}, err
	}

	b, err := json.Marshal(links)
	if err != nil {
		return []byte{}, nil
	}

	return b, nil
}

func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		val := node.Data
		val = strings.TrimSpace(val)
		return val
	}

	b := strings.Builder{}

	if node.Type == html.DocumentNode || node.Type == html.ElementNode {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			val := getText(child)
			b.WriteString(" ")
			b.WriteString(val)
		}
	}

	return b.String()
}
