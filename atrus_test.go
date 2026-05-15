package atrus_test

import (
	"testing"

	atrus "github.com/sinclairtarget/libatrus-go"
)

const md string = "# Heading\nThis is a paragraph.\n"

func TestParse(t *testing.T) {
	_, err := atrus.Parse(md, atrus.ParseLevelPost)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err)
	}
}

func TestParseAndRenderHTML(t *testing.T) {
	root, err := atrus.Parse(md, atrus.ParseLevelPost)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err)
	}

	s, err := atrus.RenderHTML(root)
	if err != nil {
		t.Fatalf("render failed with error: %v", err)
	}

	const expected = `<h1>Heading</h1>
<p>This is a paragraph.</p>
`
	if s != expected {
		t.Errorf("HTML string did not match. Got:\n%s", s)
	}
}

func TestParseAndRenderJSON(t *testing.T) {
	root, err := atrus.Parse(md, atrus.ParseLevelPost)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err)
	}

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
          "type": "paragraph",
          "children": [
            {
              "type": "text",
              "value": "This is a paragraph."
            }
          ]
        }
      ]
    }
  ]
}`
	if s != expected {
		t.Errorf("JSON string did not match. Got:\n%s", s)
	}
}
