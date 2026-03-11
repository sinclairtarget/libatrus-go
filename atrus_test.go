package atrus_test

import (
	"slices"
	"testing"

	atrus "github.com/sinclairtarget/libatrus-go"
)

const md string = "# Heading\nThis is a paragraph.\n"

func TestParseAndASTOperations(t *testing.T) {
	opts := atrus.ParseOpts{
		ParseLevel: atrus.ParseLevelPost,
	}
	root, err := atrus.Parse(md, opts)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err)
	}

	nodeType := root.Type()
	if nodeType != "root" {
		t.Errorf("expected type to be root, but got \"%s\"", nodeType)
	}

	for child := range root.Children() {
		childType := child.Type()
		if childType != "block" {
			t.Errorf(
				"expected root to only have block children, but got \"%s\"",
				childType,
			)
		}
	}

	headings := slices.Collect(root.All("heading"))
	if len(headings) != 1 {
		t.Errorf(
			"expected AST to contain one heading, but found %d",
			len(headings),
		)
	}

	paragraph := root.First("paragraph")
	if paragraph == nil {
		t.Error("expected paragraph node to be non-nil but was nil")
	}

	nodeType = paragraph.Type()
	if nodeType != "paragraph" {
		t.Errorf(
			"expected node to have type \"paragraph\", but got \"%s\"",
			nodeType,
		)
	}
}

func TestParseAndRenderHTML(t *testing.T) {
	opts := atrus.ParseOpts{
		ParseLevel: atrus.ParseLevelPost,
	}
	root, err := atrus.Parse(md, opts)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err)
	}

	s, err := atrus.RenderHTML(root)
	if err != nil {
		t.Fatalf("render failed with error: %v", err)
	}

	if len(s) <= 0 {
		t.Errorf("html string was empty")
	}
}

func TestParseAndRenderJSON(t *testing.T) {
	opts := atrus.ParseOpts{
		ParseLevel: atrus.ParseLevelPost,
	}
	root, err := atrus.Parse(md, opts)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err)
	}

	renderOpts := atrus.JSONOpts{
		Whitespace: atrus.JSONIndent2,
	}
	s, err := atrus.RenderJSON(root, renderOpts)
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
