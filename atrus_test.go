package atrus_test

import (
	"testing"

	atrus "github.com/sinclairtarget/libatrus-go"
)

func TestParseAndRender(t *testing.T) {
	const md = "# Heading\nThis is a paragraph.\n";

	astNode, err := atrus.ParseAST(md)
	if err != nil {
		t.Fatalf("parse failed with error: %v", err);
	}

	s, err := atrus.RenderJSON(astNode)
	if err != nil {
		t.Fatalf("render failed with error: %v", err);
	}

	if len(s) <= 0 {
		t.Errorf("json string was empty");
	}
}
