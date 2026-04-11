package lexer

import (
	"testing"
)

func tokenTypes(tokens []Token) []TokenType {
	types := make([]TokenType, len(tokens))
	for i, t := range tokens {
		types[i] = t.Type
	}
	return types
}

func tokenize(t *testing.T, source string) []Token {
	t.Helper()
	lex := New(source)
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
	if tokens[0].Type != TOKEN_NUMBER || tokens[0].Value != "42" {
		t.Errorf("expected NUMBER 42, got %s", tokens[0])
	}
}

func TestTokenizeDecimalNumber(t *testing.T) {
	tokens := tokenize(t, "3.14")
	if tokens[0].Type != TOKEN_NUMBER || tokens[0].Value != "3.14" {
		t.Errorf("expected NUMBER 3.14, got %s", tokens[0])
	}
}

func TestTokenizeString(t *testing.T) {
	tokens := tokenize(t, `"hello world"`)
	if tokens[0].Type != TOKEN_STRING || tokens[0].Value != "hello world" {
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
	lex := New(`"unterminated`)
	_, err := lex.Tokenize()
	if err == nil {
		t.Fatal("expected error for unterminated string")
	}
}

func TestTokenizeKeywords(t *testing.T) {
	tests := []struct {
		source   string
		expected TokenType
	}{
		{"dei", TOKEN_DEI},
		{"memang dei", TOKEN_MEMANG_DEI},
		{"tukar", TOKEN_TUKAR},
		{"eh", TOKEN_EH},
		{"walao eh", TOKEN_WALAO_EH},
		{"tapau", TOKEN_TAPAU},
		{"is it", TOKEN_IS_IT},
		{"or is it", TOKEN_OR_IS_IT},
		{"abuden", TOKEN_ABUDEN},
		{"one by one", TOKEN_ONE_BY_ONE},
		{"sampai", TOKEN_SAMPAI},
		{"dalam", TOKEN_DALAM},
		{"gaodim", TOKEN_GAODIM},
		{"no need", TOKEN_NO_NEED},
		{"jom", TOKEN_JOM},
		{"nah", TOKEN_NAH},
		{"sama", TOKEN_SAMA},
		{"tak sama", TOKEN_TAK_SAMA},
		{"dan", TOKEN_DAN},
		{"atau", TOKEN_ATAU},
		{"try see", TOKEN_TRY_SEE},
		{"alamak", TOKEN_ALAMAK},
		{"confirm need", TOKEN_CONFIRM_NEED},
		{"mampus", TOKEN_MAMPUS},
		{"can", TOKEN_CAN},
		{"cannot", TOKEN_CANNOT},
		{"kosong", TOKEN_KOSONG},
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
	if tokens[0].Type != TOKEN_IDENTIFIER || tokens[0].Value != "candle" {
		t.Errorf("expected IDENTIFIER candle, got %s", tokens[0])
	}
}

func TestTokenizeIdentifier(t *testing.T) {
	tokens := tokenize(t, "myVar_123")
	if tokens[0].Type != TOKEN_IDENTIFIER || tokens[0].Value != "myVar_123" {
		t.Errorf("expected IDENTIFIER myVar_123, got %s", tokens[0])
	}
}

func TestTokenizeSymbols(t *testing.T) {
	tokens := tokenize(t, "( ) { } [ ] ; , . + - * / % = < > !")
	expected := []TokenType{
		TOKEN_LPAREN, TOKEN_RPAREN, TOKEN_LBRACE, TOKEN_RBRACE,
		TOKEN_LBRACKET, TOKEN_RBRACKET, TOKEN_SEMICOLON, TOKEN_COMMA,
		TOKEN_DOT, TOKEN_PLUS, TOKEN_MINUS, TOKEN_STAR, TOKEN_SLASH,
		TOKEN_PERCENT, TOKEN_ASSIGN, TOKEN_LESS, TOKEN_GREATER, TOKEN_BANG,
		TOKEN_EOF,
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
	expected := []TokenType{TOKEN_LESS_EQUAL, TOKEN_GREATER_EQUAL, TOKEN_PLUS_ASSIGN, TOKEN_MINUS_ASSIGN, TOKEN_EOF}
	types := tokenTypes(tokens)
	for i := range expected {
		if types[i] != expected[i] {
			t.Errorf("token %d: expected %d, got %d", i, expected[i], types[i])
		}
	}
}

func TestTokenizeSingleLineComment(t *testing.T) {
	tokens := tokenize(t, "// eh: this is a comment\n42")
	if tokens[0].Type != TOKEN_NUMBER || tokens[0].Value != "42" {
		t.Errorf("expected NUMBER 42 after comment, got %s", tokens[0])
	}
}

func TestTokenizeMultiLineComment(t *testing.T) {
	tokens := tokenize(t, "alamak mula this is\na block comment alamak habis 42")
	if tokens[0].Type != TOKEN_NUMBER || tokens[0].Value != "42" {
		t.Errorf("expected NUMBER 42 after block comment, got %s", tokens[0])
	}
}

func TestTokenizeUnterminatedBlockComment(t *testing.T) {
	lex := New("alamak mula this never ends")
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
		if tok.Type == TOKEN_DEI && tok.Line == 2 {
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
	expected := []TokenType{TOKEN_DEI, TOKEN_IDENTIFIER, TOKEN_ASSIGN, TOKEN_STRING, TOKEN_SEMICOLON, TOKEN_EOF}
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
	lex := New("@")
	_, err := lex.Tokenize()
	if err == nil {
		t.Fatal("expected error for unexpected character")
	}
}

func TestTokenizeComplexExpression(t *testing.T) {
	tokens := tokenize(t, "is it (x > 0 dan y < 10) { eh x; }")
	// Should parse without error and contain the expected keyword tokens
	hasIsIt := false
	hasDan := false
	for _, tok := range tokens {
		if tok.Type == TOKEN_IS_IT {
			hasIsIt = true
		}
		if tok.Type == TOKEN_DAN {
			hasDan = true
		}
	}
	if !hasIsIt || !hasDan {
		t.Error("expected IS_IT and DAN tokens in complex expression")
	}
}
