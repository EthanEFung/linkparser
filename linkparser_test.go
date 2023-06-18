package linkparser

import (
	"os"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestLinkParserUseNode(t *testing.T) {
	parser := New()
	textNode := &html.Node{
		Type: html.TextNode,
		Data: "This is my text",
	}
	node := &html.Node{
		Type:       html.ElementNode,
		DataAtom:   atom.A,
		FirstChild: textNode,
		LastChild:  textNode,
		Attr: []html.Attribute{
			{Key: "href", Val: "#"},
		},
	}
	err := parser.UseNode(node)

	expected := []Link{{"#", "This is my text"}}
	actual, err := parser.Parse()
	if err != nil {
		t.Fatal("could not parse node")
	}

	assertLinks(t, expected, actual)
}

func TestLinkParserUseReader(t *testing.T) {
	parser := New()
	file, _ := os.Open("mockdata/ex1.html")
	defer file.Close()

	err := parser.UseReader(file)
	if err != nil {
		t.Fatal("cannot use file as the reader")
	}

	links, err := parser.Parse()
	if err != nil {
		t.Fatal("cannot parse using the reader")
	}
	if len(links) == 0 {
		t.Fatal("could not get links using reader")
	}
}

func TestLinkParserUseHTMLFile(t *testing.T) {
	parser := New()
	err := parser.UseHTMLFile("mockdata/ex1.html")
	if err != nil {
		t.Fatal("Cannot UseHTMLFile")
	}
	links, err := parser.Parse()
	if err != nil {
		t.Fatal("could not parse after UseHTMLFile was called")
	}
	if len(links) == 0 {
		t.Fatal("could not read links after parsing")
	}
}

func TestLinkParserParse(t *testing.T) {
	parser := New()
	parser.UseHTMLFile("mockdata/ex1.html")

	expected := []Link{
		{"/other-page", "A link to another page"},
	}
	actual, err := parser.Parse()
	if err != nil {
		t.Fatal("attempted to parse ex1.html but to no avail")
	}
	assertLinks(t, expected, actual)

	parser.UseHTMLFile("mockdata/ex2.html")
	expected = []Link{
		{"https://www.twitter.com/joncalhoun", "Check me out on twitter"},
		{"https://github.com/gophercises", "Gophercises is on Github!"},
	}
	actual, err = parser.Parse()
	if err != nil {
		t.Fatal("attempted to parse ex2.html but to no avail")
	}
	assertLinks(t, expected, actual)

	parser.UseHTMLFile("mockdata/ex3.html")
	expected = []Link{
		{"#", "Login"},
		{"/lost", "Lost? Need help?"},
		{"https://twitter.com/marcusolsson", "@marcusolsson"},
	}
	actual, err = parser.Parse()
	if err != nil {
		t.Fatal("attempted to parse ex3.html but to no avail")
	}
	assertLinks(t, expected, actual)

	parser.UseHTMLFile("mockdata/ex4.html")
	expected = []Link{
		{"/dog-cat", "dog cat"},
	}
	actual, err = parser.Parse()
	if err != nil {
		t.Fatal("attempted to parse ex4.html but to no avail")
	}
	assertLinks(t, expected, actual)
}

func assertLinks(t *testing.T, expected, actual []Link) {
	if len(expected) != len(actual) {
		t.Fatalf("expected %d links but got %d", len(expected), len(actual))
	}
	for i := 0; i < len(expected); i++ {
		e, a := expected[i], actual[i]
		if e.Href != a.Href {
			t.Fatalf(`expected link to have href "%s" but got "%s"`, e.Href, a.Href)
		}
		if e.Text != a.Text {
			t.Fatalf(`expected link to have text "%s" but got "%s`, e.Text, a.Text)
		}
	}
}
