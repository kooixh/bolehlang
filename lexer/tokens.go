package lexer

import "fmt"

type TokenType int

const (
	// Variables
	TOKEN_DEI TokenType = iota
	TOKEN_MEMANG_DEI
	TOKEN_TUKAR

	// Literals
	TOKEN_NUMBER
	TOKEN_STRING
	TOKEN_CAN
	TOKEN_CANNOT
	TOKEN_KOSONG

	// Print
	TOKEN_EH
	TOKEN_WALAO_EH

	// Stdin read
	TOKEN_LAI

	// Import
	TOKEN_TAPAU

	// Control flow
	TOKEN_IS_IT
	TOKEN_OR_IS_IT
	TOKEN_ABUDEN

	// Loops (all use "dari")
	TOKEN_ONE_BY_ONE
	TOKEN_SAMPAI
	TOKEN_DALAM
	TOKEN_GAODIM
	TOKEN_NO_NEED

	// Functions
	TOKEN_JOM
	TOKEN_NAH

	// structs
	TOKEN_KRONI

	// Comparison keywords
	TOKEN_SAMA
	TOKEN_TAK_SAMA
	TOKEN_ALSO
	TOKEN_IF_NOT

	// Error handling
	TOKEN_TRY_SEE
	TOKEN_ALAMAK
	TOKEN_CONFIRM_NEED
	TOKEN_MAMPUS

	// Symbols
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_LBRACKET
	TOKEN_RBRACKET
	TOKEN_SEMICOLON
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_STAR
	TOKEN_SLASH
	TOKEN_PERCENT
	TOKEN_ASSIGN
	TOKEN_LESS
	TOKEN_GREATER
	TOKEN_LESS_EQUAL
	TOKEN_GREATER_EQUAL
	TOKEN_BANG
	TOKEN_PLUS_ASSIGN
	TOKEN_MINUS_ASSIGN
	TOKEN_COLON

	// Other
	TOKEN_IDENTIFIER
	TOKEN_EOF
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%d, %q, %d:%d)", t.Type, t.Value, t.Line, t.Col)
}

// Keywords ordered longest-first for greedy matching
var keywords = []struct {
	text    string
	tokType TokenType
}{
	{"memang dei", TOKEN_MEMANG_DEI},
	{"or is it", TOKEN_OR_IS_IT},
	{"jom", TOKEN_JOM},
	{"nah", TOKEN_NAH},
	{"abuden", TOKEN_ABUDEN},
	{"cannot", TOKEN_CANNOT},
	{"tak sama", TOKEN_TAK_SAMA},
	{"confirm need", TOKEN_CONFIRM_NEED},
	{"gaodim", TOKEN_GAODIM},
	{"try see", TOKEN_TRY_SEE},
	{"no need", TOKEN_NO_NEED},
	{"sampai", TOKEN_SAMPAI},
	{"walao eh", TOKEN_WALAO_EH},
	{"eh", TOKEN_EH},
	{"lai", TOKEN_LAI},
	{"kosong", TOKEN_KOSONG},
	{"mampus", TOKEN_MAMPUS},
	{"alamak", TOKEN_ALAMAK},
	{"is it", TOKEN_IS_IT},
	{"dalam", TOKEN_DALAM},
	{"tapau", TOKEN_TAPAU},
	{"tukar", TOKEN_TUKAR},
	{"can", TOKEN_CAN},
	{"one by one", TOKEN_ONE_BY_ONE},
	{"sama", TOKEN_SAMA},
	{"if not", TOKEN_IF_NOT},
	{"also", TOKEN_ALSO},
	{"dei", TOKEN_DEI},
	{"kroni" , TOKEN_KRONI},
}
