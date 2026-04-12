package tests

import (
	"testing"

	"github.com/kooixh/bolehlang/lexer"
)

func tokenTypes(tokens []lexer.Token) []lexer.TokenType {
	types := make([]lexer.TokenType, len(tokens))
	for i, t := range tokens {
		types[i] = t.Type
	}
	return types
}

func tokenize(t *testing.T, source string) []lexer.Token {
	t.Helper()
	lex := lexer.New(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("unexpected tokenize error: %s", err)
	}
	return tokens
}

func TestTokenizeNumber(t *testing.T) {
	tokens := tokenize(t, "42")
	if len(tokens) != 2 { // number + EOF
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != lexer.TOKEN_NUMBER || tokens[0].Value != "42" {
		t.Errorf("expected NUMBER 42, got %s", tokens[0])
	}
}

func TestTokenizeDecimalNumber(t *testing.T) {
	tokens := tokenize(t, "3.14")
	if tokens[0].Type != lexer.TOKEN_NUMBER || tokens[0].Value != "3.14" {
		t.Errorf("expected NUMBER 3.14, got %s", tokens[0])
	}
}

func TestTokenizeString(t *testing.T) {
	tokens := tokenize(t, `"hello world"`)
	if tokens[0].Type != lexer.TOKEN_STRING || tokens[0].Value != "hello world" {
		t.Errorf("expected STRING hello world, got %s", tokens[0])
	}
}

func TestTokenizeStringEscapes(t *testing.T) {
	tokens := tokenize(t, `"line1\nline2"`)
	if tokens[0].Value != "line1\nline2" {
		t.Errorf("expected newline escape, got %q", tokens[0].Value)
	}
}

func TestTokenizeUnterminatedString(t *testing.T) {
	lex := lexer.New(`"unterminated`)
	_, err := lex.Tokenize()
	if err == nil {
		t.Fatal("expected error for unterminated string")
	}
}

func TestTokenizeKeywords(t *testing.T) {
	tests := []struct {
		source   string
		expected lexer.TokenType
	}{
		{"dei", lexer.TOKEN_DEI},
		{"memang dei", lexer.TOKEN_MEMANG_DEI},
		{"tukar", lexer.TOKEN_TUKAR},
		{"eh", lexer.TOKEN_EH},
		{"walao eh", lexer.TOKEN_WALAO_EH},
		{"tapau", lexer.TOKEN_TAPAU},
		{"is it", lexer.TOKEN_IS_IT},
		{"or is it", lexer.TOKEN_OR_IS_IT},
		{"abuden", lexer.TOKEN_ABUDEN},
		{"one by one", lexer.TOKEN_ONE_BY_ONE},
		{"sampai", lexer.TOKEN_SAMPAI},
		{"dalam", lexer.TOKEN_DALAM},
		{"gaodim", lexer.TOKEN_GAODIM},
		{"no need", lexer.TOKEN_NO_NEED},
		{"jom", lexer.TOKEN_JOM},
		{"nah", lexer.TOKEN_NAH},
		{"sama", lexer.TOKEN_SAMA},
		{"tak sama", lexer.TOKEN_TAK_SAMA},
		{"also", lexer.TOKEN_ALSO},
		{"if not", lexer.TOKEN_IF_NOT},
		{"try see", lexer.TOKEN_TRY_SEE},
		{"alamak", lexer.TOKEN_ALAMAK},
		{"confirm need", lexer.TOKEN_CONFIRM_NEED},
		{"mampus", lexer.TOKEN_MAMPUS},
		{"can", lexer.TOKEN_CAN},
		{"cannot", lexer.TOKEN_CANNOT},
		{"kosong", lexer.TOKEN_KOSONG},
		{"lai", lexer.TOKEN_LAI},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			tokens := tokenize(t, tt.source+";")
			if tokens[0].Type != tt.expected {
				t.Errorf("expected token type %d for %q, got %d (%s)", tt.expected, tt.source, tokens[0].Type, tokens[0])
			}
		})
	}
}

func TestTokenizeKeywordBoundary(t *testing.T) {
	// "candle" should be an identifier, not TOKEN_CAN + "dle"
	tokens := tokenize(t, "candle")
	if tokens[0].Type != lexer.TOKEN_IDENTIFIER || tokens[0].Value != "candle" {
		t.Errorf("expected IDENTIFIER candle, got %s", tokens[0])
	}
}

func TestTokenizeIdentifier(t *testing.T) {
	tokens := tokenize(t, "myVar_123")
	if tokens[0].Type != lexer.TOKEN_IDENTIFIER || tokens[0].Value != "myVar_123" {
		t.Errorf("expected IDENTIFIER myVar_123, got %s", tokens[0])
	}
}

func TestTokenizeSymbols(t *testing.T) {
	tokens := tokenize(t, "( ) { } [ ] ; , . + - * / % = < > !")
	expected := []lexer.TokenType{
		lexer.TOKEN_LPAREN, lexer.TOKEN_RPAREN, lexer.TOKEN_LBRACE, lexer.TOKEN_RBRACE,
		lexer.TOKEN_LBRACKET, lexer.TOKEN_RBRACKET, lexer.TOKEN_SEMICOLON, lexer.TOKEN_COMMA,
		lexer.TOKEN_DOT, lexer.TOKEN_PLUS, lexer.TOKEN_MINUS, lexer.TOKEN_STAR, lexer.TOKEN_SLASH,
		lexer.TOKEN_PERCENT, lexer.TOKEN_ASSIGN, lexer.TOKEN_LESS, lexer.TOKEN_GREATER, lexer.TOKEN_BANG,
		lexer.TOKEN_EOF,
	}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(types))
	}
	for i := range expected {
		if types[i] != expected[i] {
			t.Errorf("token %d: expected %d, got %d", i, expected[i], types[i])
		}
	}
}

func TestTokenizeTwoCharSymbols(t *testing.T) {
	tokens := tokenize(t, "<= >= += -=")
	expected := []lexer.TokenType{lexer.TOKEN_LESS_EQUAL, lexer.TOKEN_GREATER_EQUAL, lexer.TOKEN_PLUS_ASSIGN, lexer.TOKEN_MINUS_ASSIGN, lexer.TOKEN_EOF}
	types := tokenTypes(tokens)
	for i := range expected {
		if types[i] != expected[i] {
			t.Errorf("token %d: expected %d, got %d", i, expected[i], types[i])
		}
	}
}

func TestTokenizeSingleLineComment(t *testing.T) {
	tokens := tokenize(t, "// eh: this is a comment\n42")
	if tokens[0].Type != lexer.TOKEN_NUMBER || tokens[0].Value != "42" {
		t.Errorf("expected NUMBER 42 after comment, got %s", tokens[0])
	}
}

func TestTokenizeMultiLineComment(t *testing.T) {
	tokens := tokenize(t, "alamak mula this is\na block comment alamak habis 42")
	if tokens[0].Type != lexer.TOKEN_NUMBER || tokens[0].Value != "42" {
		t.Errorf("expected NUMBER 42 after block comment, got %s", tokens[0])
	}
}

func TestTokenizeUnterminatedBlockComment(t *testing.T) {
	lex := lexer.New("alamak mula this never ends")
	_, err := lex.Tokenize()
	if err == nil {
		t.Fatal("expected error for unterminated block comment")
	}
}

func TestTokenizeLineTracking(t *testing.T) {
	tokens := tokenize(t, "dei x = 1;\ndei y = 2;")
	// "dei" on line 1, "dei" on line 2
	if tokens[0].Line != 1 {
		t.Errorf("expected line 1, got %d", tokens[0].Line)
	}
	// Find second "dei"
	found := false
	for _, tok := range tokens {
		if tok.Type == lexer.TOKEN_DEI && tok.Line == 2 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected second 'dei' on line 2")
	}
}

func TestTokenizeVarDeclaration(t *testing.T) {
	tokens := tokenize(t, `dei nama = "Ali";`)
	expected := []lexer.TokenType{lexer.TOKEN_DEI, lexer.TOKEN_IDENTIFIER, lexer.TOKEN_ASSIGN, lexer.TOKEN_STRING, lexer.TOKEN_SEMICOLON, lexer.TOKEN_EOF}
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(types))
	}
	for i := range expected {
		if types[i] != expected[i] {
			t.Errorf("token %d: expected %d, got %d", i, expected[i], types[i])
		}
	}
}

func TestTokenizeUnexpectedCharacter(t *testing.T) {
	lex := lexer.New("@")
	_, err := lex.Tokenize()
	if err == nil {
		t.Fatal("expected error for unexpected character")
	}
}

func TestTokenizeComplexExpression(t *testing.T) {
	tokens := tokenize(t, "is it (x > 0 also y < 10) { eh x; }")
	// Should parse without error and contain the expected keyword tokens
	hasIsIt := false
	hasAlso := false
	for _, tok := range tokens {
		if tok.Type == lexer.TOKEN_IS_IT {
			hasIsIt = true
		}
		if tok.Type == lexer.TOKEN_ALSO {
			hasAlso = true
		}
	}
	if !hasIsIt || !hasAlso {
		t.Error("expected IS_IT and ALSO tokens in complex expression")
	}
}
