package parser

import (
	"testing"

	"github.com/simonwater/gopression/values"
)

func tokenTypes(tokens []values.Token) []values.TokenType {
	types := make([]values.TokenType, 0, len(tokens))
	for _, t := range tokens {
		types = append(types, t.Type)
	}
	return types
}

func TestScanTokens_SimpleSymbols(t *testing.T) {
	src := "( ) { } , . - + * ; % ! = > < /"
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()
	expected := []values.TokenType{
		values.LEFT_PAREN, values.RIGHT_PAREN, values.LEFT_BRACE, values.RIGHT_BRACE,
		values.COMMA, values.DOT, values.MINUS, values.PLUS, values.STAR,
		values.SEMICOLON, values.PERCENT, values.BANG, values.EQUAL,
		values.GREATER, values.LESS, values.SLASH, values.EOF,
	}
	got := tokenTypes(tokens)
	if len(got) != len(expected) {
		t.Fatalf("token count mismatch: got %d, want %d", len(got), len(expected))
	}
	for i, typ := range expected {
		if got[i] != typ {
			t.Errorf("token %d: got %v, want %v", i, got[i], typ)
		}
	}
}

func TestScanTokens_Operators(t *testing.T) {
	src := "== != >= <= ** && || //"
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()
	expected := []values.TokenType{
		values.EQUAL_EQUAL, values.BANG_EQUAL, values.GREATER_EQUAL, values.LESS_EQUAL,
		values.STARSTAR, values.AND, values.OR, values.EOF,
	}
	got := tokenTypes(tokens)
	for i, typ := range expected {
		if got[i] != typ {
			t.Errorf("token %d: got %v, want %v", i, got[i], typ)
		}
	}
}

func TestScanTokens_Numbers(t *testing.T) {
	src := "123 45.67"
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()
	if tokens[0].Type != values.NUMBER || tokens[1].Type != values.NUMBER {
		t.Errorf("expected NUMBER tokens")
	}
	if tokens[0].Literal.AsInteger() != 123 {
		t.Errorf("expected 123, got %v", tokens[0].Literal)
	}
	if tokens[1].Literal.AsDouble() != 45.67 {
		t.Errorf("expected 45.67, got %v", tokens[1].Literal)
	}
}

func TestScanTokens_Strings(t *testing.T) {
	src := `"hello" "中文"`
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()
	if tokens[0].Type != values.STRING || tokens[1].Type != values.STRING {
		t.Errorf("expected STRING tokens")
	}
	if tokens[0].Literal.String() != "hello" {
		t.Errorf("expected 'hello', got %v", tokens[0].Literal.String())
	}
	if tokens[1].Literal.String() != "中文" {
		t.Errorf("expected '中文', got %v", tokens[1].Literal.String())
	}
}

func TestScanTokens_KeywordsAndIdentifiers(t *testing.T) {
	src := "class fun var if else true false null abc _var 中文"
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()
	expected := []values.TokenType{
		values.CLASS, values.FUN, values.VAR, values.IF, values.ELSE,
		values.TRUE, values.FALSE, values.NULL,
		values.IDENTIFIER, values.IDENTIFIER, values.IDENTIFIER, values.EOF,
	}
	got := tokenTypes(tokens)
	for i, typ := range expected {
		if got[i] != typ {
			t.Errorf("token %d: got %v, want %v", i, got[i], typ)
		}
	}
}

func TestScanTokens_CommentsAndWhitespace(t *testing.T) {
	src := `
    // this is a comment
    var x = 1 // another comment
    `
	scanner := NewScanner(src)
	tokens := scanner.ScanTokens()
	// Should parse: var IDENTIFIER = NUMBER EOF
	if tokens[0].Type != values.VAR || tokens[1].Type != values.IDENTIFIER ||
		tokens[2].Type != values.EQUAL || tokens[3].Type != values.NUMBER ||
		tokens[4].Type != values.EOF {
		t.Errorf("unexpected tokens: %+v", tokens)
	}
}

func TestScanTokens_UnterminatedString(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unterminated string")
		}
	}()
	src := `"unterminated`
	scanner := NewScanner(src)
	scanner.ScanTokens()
}

func TestScanTokens_UnknownCharacter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unknown character")
		}
	}()
	src := "@"
	scanner := NewScanner(src)
	scanner.ScanTokens()
}
