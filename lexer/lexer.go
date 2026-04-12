package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	source []rune
	pos    int
	line   int
	col    int
}

func New(source string) *Lexer {
	return &Lexer{
		source: []rune(source),
		pos:    0,
		line:   1,
		col:    1,
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for l.pos < len(l.source) {
		// Skip whitespace
		if unicode.IsSpace(l.source[l.pos]) {
			l.skipWhitespace()
			continue
		}

		// Single-line comments: // eh:
		if l.pos+5 <= len(l.source) && string(l.source[l.pos:l.pos+5]) == "// eh" {
			l.skipSingleLineComment()
			continue
		}

		// Multi-line comments: alamak mula ... alamak habis
		if l.matchAt("alamak mula") {
			err := l.skipMultiLineComment()
			if err != nil {
				return nil, err
			}
			continue
		}

		// String literals
		if l.source[l.pos] == '"' {
			tok, err := l.readString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
			continue
		}

		// Number literals
		if unicode.IsDigit(l.source[l.pos]) {
			tok := l.readNumber()
			tokens = append(tokens, tok)
			continue
		}

		// Two-character symbols
		if l.pos+1 < len(l.source) {
			twoChar := string(l.source[l.pos : l.pos+2])
			switch twoChar {
			case "<=":
				tokens = append(tokens, Token{TOKEN_LESS_EQUAL, "<=", l.line, l.col})
				l.advance()
				l.advance()
				continue
			case ">=":
				tokens = append(tokens, Token{TOKEN_GREATER_EQUAL, ">=", l.line, l.col})
				l.advance()
				l.advance()
				continue
			case "+=":
				tokens = append(tokens, Token{TOKEN_PLUS_ASSIGN, "+=", l.line, l.col})
				l.advance()
				l.advance()
				continue
			case "-=":
				tokens = append(tokens, Token{TOKEN_MINUS_ASSIGN, "-=", l.line, l.col})
				l.advance()
				l.advance()
				continue
			}
		}

		// Keywords (greedy, longest-first)
		if matched, tok := l.matchKeyword(); matched {
			tokens = append(tokens, tok)
			continue
		}

		// Identifiers
		if unicode.IsLetter(l.source[l.pos]) || l.source[l.pos] == '_' {
			tok := l.readIdentifier()
			tokens = append(tokens, tok)
			continue
		}

		// Single-character punctuation
		if tok, ok := l.readPunctuation(); ok {
			tokens = append(tokens, tok)
			continue
		}

		return nil, fmt.Errorf("line %d, col %d: unexpected character '%c'", l.line, l.col, l.source[l.pos])
	}

	tokens = append(tokens, Token{TOKEN_EOF, "", l.line, l.col})
	return tokens, nil
}

func (l *Lexer) advance() {
	if l.pos < len(l.source) {
		if l.source[l.pos] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.pos++
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.source) && unicode.IsSpace(l.source[l.pos]) {
		l.advance()
	}
}

func (l *Lexer) skipSingleLineComment() {
	for l.pos < len(l.source) && l.source[l.pos] != '\n' {
		l.advance()
	}
}

func (l *Lexer) matchAt(s string) bool {
	runes := []rune(s)
	if l.pos+len(runes) > len(l.source) {
		return false
	}
	for i, r := range runes {
		if l.source[l.pos+i] != r {
			return false
		}
	}
	return true
}

func (l *Lexer) skipMultiLineComment() error {
	startLine := l.line
	// Skip "alamak mula"
	for range []rune("alamak mula") {
		l.advance()
	}
	for l.pos < len(l.source) {
		if l.matchAt("alamak habis") {
			for range []rune("alamak habis") {
				l.advance()
			}
			return nil
		}
		l.advance()
	}
	return fmt.Errorf("line %d: unterminated block comment (alamak mula without alamak habis)", startLine)
}

func (l *Lexer) readString() (Token, error) {
	startLine := l.line
	startCol := l.col
	l.advance() // skip opening quote
	var sb strings.Builder
	for l.pos < len(l.source) {
		ch := l.source[l.pos]
		if ch == '"' {
			l.advance() // skip closing quote
			return Token{TOKEN_STRING, sb.String(), startLine, startCol}, nil
		}
		if ch == '\\' && l.pos+1 < len(l.source) {
			l.advance()
			switch l.source[l.pos] {
			case 'n':
				sb.WriteRune('\n')
			case 't':
				sb.WriteRune('\t')
			case '\\':
				sb.WriteRune('\\')
			case '"':
				sb.WriteRune('"')
			default:
				sb.WriteRune('\\')
				sb.WriteRune(l.source[l.pos])
			}
			l.advance()
			continue
		}
		if ch == '\n' {
			return Token{}, fmt.Errorf("line %d: unterminated string", startLine)
		}
		sb.WriteRune(ch)
		l.advance()
	}
	return Token{}, fmt.Errorf("line %d: unterminated string", startLine)
}

func (l *Lexer) readNumber() Token {
	startCol := l.col
	var sb strings.Builder
	hasDecimal := false
	for l.pos < len(l.source) {
		ch := l.source[l.pos]
		if unicode.IsDigit(ch) {
			sb.WriteRune(ch)
			l.advance()
		} else if ch == '.' && !hasDecimal {
			hasDecimal = true
			sb.WriteRune(ch)
			l.advance()
		} else {
			break
		}
	}
	return Token{TOKEN_NUMBER, sb.String(), l.line, startCol}
}

func (l *Lexer) matchKeyword() (bool, Token) {
	remaining := string(l.source[l.pos:])
	for _, kw := range keywords {
		if strings.HasPrefix(remaining, kw.text) {
			// Check word boundary: the character after the keyword must not be a letter/digit/underscore
			endPos := l.pos + len([]rune(kw.text))
			if endPos < len(l.source) {
				next := l.source[endPos]
				if unicode.IsLetter(next) || unicode.IsDigit(next) || next == '_' {
					continue
				}
			}
			tok := Token{kw.tokType, kw.text, l.line, l.col}
			for range []rune(kw.text) {
				l.advance()
			}
			return true, tok
		}
	}
	return false, Token{}
}

func (l *Lexer) readIdentifier() Token {
	startCol := l.col
	var sb strings.Builder
	for l.pos < len(l.source) {
		ch := l.source[l.pos]
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			sb.WriteRune(ch)
			l.advance()
		} else {
			break
		}
	}
	return Token{TOKEN_IDENTIFIER, sb.String(), l.line, startCol}
}

func (l *Lexer) readPunctuation() (Token, bool) {
	ch := l.source[l.pos]
	line, col := l.line, l.col
	var tokType TokenType
	switch ch {
	case '(':
		tokType = TOKEN_LPAREN
	case ')':
		tokType = TOKEN_RPAREN
	case '{':
		tokType = TOKEN_LBRACE
	case '}':
		tokType = TOKEN_RBRACE
	case '[':
		tokType = TOKEN_LBRACKET
	case ']':
		tokType = TOKEN_RBRACKET
	case ';':
		tokType = TOKEN_SEMICOLON
	case ',':
		tokType = TOKEN_COMMA
	case '.':
		tokType = TOKEN_DOT
	case '+':
		tokType = TOKEN_PLUS
	case '-':
		tokType = TOKEN_MINUS
	case '*':
		tokType = TOKEN_STAR
	case '/':
		tokType = TOKEN_SLASH
	case '%':
		tokType = TOKEN_PERCENT
	case '=':
		tokType = TOKEN_ASSIGN
	case '<':
		tokType = TOKEN_LESS
	case '>':
		tokType = TOKEN_GREATER
	case '!':
		tokType = TOKEN_BANG
	case ':':
		tokType = TOKEN_COLON
	default:
		return Token{}, false
	}
	l.advance()
	return Token{tokType, string(ch), line, col}, true
}
