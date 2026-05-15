package atrus_test

import (
	"testing"

	atrus "github.com/sinclairtarget/libatrus-go"
)

func parse(t *testing.T, md string) *atrus.ASTNode {
	root, err := atrus.Parse(md, atrus.ParseLevelPost)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err)
	}
	return root
}

func TestTypeAndChildren(t *testing.T) {
	root := parse(t, "# Heading\nThis is a paragraph.\n")

	nodeType := root.Type()
	if nodeType != "root" {
		t.Errorf("expected type to be \"root\", but got \"%s\"", nodeType)
	}

	for _, child := range root.Children() {
		childType := child.Type()
		if childType != "block" {
			t.Errorf(
				"expected root to only have block children, but got \"%s\"",
				childType,
			)
		}
	}
}

func TestHeading(t *testing.T) {
	root := parse(t, "# Foo bar\n")

	node := root.Children()[0].Children()[0]
	nodeType := node.Type()
	if nodeType != "heading" {
		t.Fatalf("expected type to be \"heading\", but got \"%s\"", nodeType)
	}

	heading := node.Heading()
	if heading.Depth != 1 {
		t.Errorf("expected depth to be 1, but got %d", heading.Depth)
	}
}

func TestText(t *testing.T) {
	root := parse(t, "Foo bar.\n")

	node := root.Children()[0].Children()[0].Children()[0]
	nodeType := node.Type()
	if nodeType != "text" {
		t.Fatalf("expected type to be \"text\", but got \"%s\"", nodeType)
	}

	text := node.Text()
	expected := "Foo bar."
	if text.Value != expected {
		t.Errorf(
			"expected text value of \"%s\", but got \"%s\"",
			expected,
			text.Value,
		)
	}
}

func TestCode(t *testing.T) {
	root := parse(t, "```python\ndef foo(): pass\n```\n")

	node := root.Children()[0].Children()[0]
	nodeType := node.Type()
	if nodeType != "code" {
		t.Fatalf("expected type to be \"code\", but got \"%s\"", nodeType)
	}

	code := node.Code()
	expected := "def foo(): pass"
	if code.Value != expected {
		t.Errorf(
			"expected code value of \"%s\", but got \"%s\"",
			expected,
			code.Value,
		)
	}

	expected = "python"
	if code.Lang != expected {
		t.Errorf(
			"expected code lang of \"%s\", but got \"%s\"",
			expected,
			code.Lang,
		)
	}
}

func TestLink(t *testing.T) {
	root := parse(t, "[foo](/foo \"bim bam\")\n")

	node := root.Children()[0].Children()[0].Children()[0]
	nodeType := node.Type()
	if nodeType != "link" {
		t.Fatalf("expected type to be \"link\", but got \"%s\"", nodeType)
	}

	link := node.Link()
	expected := "/foo"
	if link.URL != expected {
		t.Errorf(
			"expected link URL of \"%s\", but got \"%s\"",
			expected,
			link.URL,
		)
	}

	expected = "bim bam"
	if link.Title != expected {
		t.Errorf(
			"expected link title of \"%s\", but got \"%s\"",
			expected,
			link.Title,
		)
	}
}

func TestImage(t *testing.T) {
	root := parse(t, "![foo](/foo.jpg \"bim bam\")\n")

	node := root.Children()[0].Children()[0].Children()[0]
	nodeType := node.Type()
	if nodeType != "image" {
		t.Fatalf("expected type to be \"image\", but got \"%s\"", nodeType)
	}

	image := node.Image()
	expected := "/foo.jpg"
	if image.URL != expected {
		t.Errorf(
			"expected image URL of \"%s\", but got \"%s\"",
			expected,
			image.URL,
		)
	}

	expected = "bim bam"
	if image.Title != expected {
		t.Errorf(
			"expected image title of \"%s\", but got \"%s\"",
			expected,
			image.Title,
		)
	}

	expected = "foo"
	if image.Alt != expected {
		t.Errorf(
			"expected image alt of \"%s\", but got \"%s\"",
			expected,
			image.Alt,
		)
	}
}
