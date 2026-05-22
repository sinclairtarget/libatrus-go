package atrus_test

import (
	"slices"
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

func TestTransform(t *testing.T) {
	root := parse(t, "# Heading\nThis is a paragraph.\n")

	block := root.Children()[0]
	nodeType := block.Type()
	if nodeType != "block" {
		t.Fatalf("expected type to be \"block\", but got \"%s\"", nodeType)
	}

	// Get reference to child we will replace.
	// Not necessary for the replacement but we want to make sure we don't
	// cause a double free by having an active reference to the replaced child.
	paragraph := block.Children()[1]
	nodeType = paragraph.Type()
	if nodeType != "paragraph" {
		t.Fatalf("expected type to be \"paragraph\", but got \"%s\"", nodeType)
	}

	html, err := atrus.CreateHTMLNode("<div class=\"foo\"><p>Hi!</p></div>")
	if err != nil {
		t.Fatalf("node create failed with error: %v", err)
	}

	// This invalidates `text`. The cNode pointer points to garbage.
	block.ReplaceChild(1, html)

	// Create second reference to the html node
	html2 := block.Children()[1]
	nodeType = html2.Type()
	if nodeType != "html" {
		t.Fatalf("expected type to be \"html\", but got \"%s\"", nodeType)
	}
	if html2.HTML().Value != html.HTML().Value {
		t.Errorf("expected html values to be equal")
	}

	// Try rendering
	s, err := atrus.RenderJSON(root, atrus.JSONIndent2)
	if err != nil {
		t.Fatalf("render failed with error: %v", err)
	}

	const expected = `{
  "type": "root",
  "children": [
    {
      "type": "block",
      "children": [
        {
          "type": "heading",
          "depth": 1,
          "children": [
            {
              "type": "text",
              "value": "Heading"
            }
          ]
        },
        {
          "type": "html",
          "value": "<div class=\"foo\"><p>Hi!</p></div>"
        }
      ]
    }
  ]
}`
	if s != expected {
		t.Errorf(
			"JSON string did not match. Expected:\n%s\nGot:\n%s",
			expected,
			s,
		)
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

	if code.Filename != "" {
		t.Errorf("expected code filename to be an empty string")
	}

	if len(code.EmphasizeLines) != 0 {
		t.Errorf("expected emphasize lines to be empty")
	}
}

func TestCodeDirective(t *testing.T) {
	root := parse(t, `:::{code} python
:filename: foo.py
:linenos:
:emphasize-lines: 2-4
def foo():
    x = 1
	y = x + 2
	return x + y
:::
`)

	// root -> mystDirective -> code
	// TODO: Update when libatrus removes redundant directive nodes in post
	// transforms.
	node := root.Children()[0].Children()[0].Children()[0];
	nodeType := node.Type()
	if nodeType != "code" {
		t.Fatalf("expected type to be \"code\", but got \"%s\"", nodeType)
	}

	code := node.Code()
	if !code.ShowLineNumbers {
		t.Errorf("expected ShowLineNumbers to be true")
	}

	expectedFilename := "foo.py"
	if code.Filename != expectedFilename {
		t.Errorf(
			"expected filename to be \"%s\", but got \"%s\"",
			expectedFilename,
			code.Filename,
		)
	}

	expectedEmphasizeLines := []uint{2, 3, 4};
	if !slices.Equal(code.EmphasizeLines, expectedEmphasizeLines) {
		t.Errorf(
			"expected emphasize lines to be %v, but got %v",
			expectedEmphasizeLines,
			code.EmphasizeLines,
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
