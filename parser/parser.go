package parser

import (
	"fmt"

	"github.com/kooixh/bolehlang/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) Parse() (*Program, error) {
	program := &Program{}

	// Parse body until EOF
	for !p.check(lexer.TOKEN_EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
	}

	return program, nil
}

func (p *Parser) parseStatement() (Node, error) {
	tok := p.current()

	switch tok.Type {
	case lexer.TOKEN_TAPAU:
		return p.parseImport()
	case lexer.TOKEN_DEI:
		return p.parseVarDeclaration()
	case lexer.TOKEN_MEMANG_DEI:
		return p.parseConstDeclaration()
	case lexer.TOKEN_KRONI:
		return p.parseStructDeclaration()
	case lexer.TOKEN_TUKAR:
		return p.parseReassignment()
	case lexer.TOKEN_EH:
		return p.parsePrint()
	case lexer.TOKEN_WALAO_EH:
		return p.parseWarn()
	case lexer.TOKEN_IS_IT:
		return p.parseIf()
	case lexer.TOKEN_ONE_BY_ONE:
		return p.parseFor()
	case lexer.TOKEN_JOM:
		return p.parseFunctionDeclaration()
	case lexer.TOKEN_NAH:
		return p.parseReturn()
	case lexer.TOKEN_GAODIM:
		p.advance()
		p.skipSemicolon()
		return &BreakStatement{Line: tok.Line}, nil
	case lexer.TOKEN_NO_NEED:
		p.advance()
		p.skipSemicolon()
		return &ContinueStatement{Line: tok.Line}, nil
	case lexer.TOKEN_TRY_SEE:
		return p.parseTryCatch()
	case lexer.TOKEN_MAMPUS:
		return p.parseThrow()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseImport() (Node, error) {
	line := p.current().Line
	p.advance() // skip "tapau"

	if !p.check(lexer.TOKEN_STRING) {
		return nil, p.error("expected file path string after 'tapau'")
	}
	path := p.current().Value
	p.advance()

	p.skipSemicolon()
	return &ImportStatement{Path: path, Line: line}, nil
}

func (p *Parser) parseVarDeclaration() (Node, error) {
	line := p.current().Line
	p.advance() // skip "got"

	if !p.check(lexer.TOKEN_IDENTIFIER) {
		return nil, p.error("expected variable name after 'dei'")
	}
	name := p.current().Value
	p.advance()

	typeName := p.parseOptionalType()

	if err := p.expect(lexer.TOKEN_ASSIGN, "="); err != nil {
		return nil, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipSemicolon()
	return &VarDeclaration{Name: name, TypeName: typeName, Value: value, Line: line}, nil
}

func (p *Parser) parseConstDeclaration() (Node, error) {
	line := p.current().Line
	p.advance() // skip "confirm got"

	if !p.check(lexer.TOKEN_IDENTIFIER) {
		return nil, p.error("expected constant name after 'confirm got'")
	}
	name := p.current().Value
	p.advance()

	typeName := p.parseOptionalType()

	if err := p.expect(lexer.TOKEN_ASSIGN, "="); err != nil {
		return nil, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipSemicolon()
	return &ConstDeclaration{Name: name, TypeName: typeName, Value: value, Line: line}, nil
}

func (p *Parser) parseStructDeclaration()(Node , error){
	line := p.current().Line
	p.advance()

	if(!p.check(lexer.TOKEN_IDENTIFIER)){
		return nil , p.error("expected kroni member after 'kroni'")
	}

	name:=p.current().Value
	p.advance()

	if err := p.expect(lexer.TOKEN_LBRACE , "{"); err != nil {
		return nil , err
	}

	var fields[] string

	for !p.check(lexer.TOKEN_RBRACE) && !p.check(lexer.TOKEN_EOF){
		if !p.check(lexer.TOKEN_IDENTIFIER){
			return nil , p.error("expected field name inside kroni")
		}

		fields = append(fields, p.current().Value)
		p.advance()

		if p.check(lexer.TOKEN_COMMA){
			p.advance()
		} else if p.check(lexer.TOKEN_SEMICOLON){
			p.advance()
		}
	}

	if err := p.expect(lexer.TOKEN_RBRACE , "}"); err != nil {
		return nil , err
	}

	p.skipSemicolon()
	return &StructDeclaration{Name: name , Fields: fields , Line: line} , nil
}

func (p *Parser) parseReassignment() (Node, error) {
	line := p.current().Line
	p.advance() // skip "tukar"

	if !p.check(lexer.TOKEN_IDENTIFIER) {
		return nil, p.error("expected variable name after 'tukar'")
	}
	name := p.current().Value
	p.advance()

	op := "="
	if p.check(lexer.TOKEN_ASSIGN) {
		p.advance()
	} else if p.check(lexer.TOKEN_PLUS_ASSIGN) {
		op = "+="
		p.advance()
	} else if p.check(lexer.TOKEN_MINUS_ASSIGN) {
		op = "-="
		p.advance()
	} else {
		return nil, p.error("expected '=', '+=', or '-=' after variable name")
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipSemicolon()
	return &Assignment{Name: name, Value: value, Operator: op, Line: line}, nil
}

func (p *Parser) parsePrint() (Node, error) {
	line := p.current().Line
	p.advance() // skip "cakap"

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipSemicolon()
	return &PrintStatement{Value: value, Line: line}, nil
}

func (p *Parser) parseWarn() (Node, error) {
	line := p.current().Line
	p.advance() // skip "walao eh"

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipSemicolon()
	return &WarnStatement{Value: value, Line: line}, nil
}

func (p *Parser) parseIf() (Node, error) {
	line := p.current().Line
	p.advance() // skip "kalau"

	if err := p.expect(lexer.TOKEN_LPAREN, "("); err != nil {
		return nil, err
	}
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	stmt := &IfStatement{Condition: condition, Body: body, Line: line}

	// Parse "atau kalau" (else if)
	for p.check(lexer.TOKEN_OR_IS_IT) {
		p.advance()
		if err := p.expect(lexer.TOKEN_LPAREN, "("); err != nil {
			return nil, err
		}
		cond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
			return nil, err
		}
		elseIfBody, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		stmt.ElseIfCases = append(stmt.ElseIfCases, ElseIfCase{Condition: cond, Body: elseIfBody})
	}

	// Parse "kalau tak" (else)
	if p.check(lexer.TOKEN_ABUDEN) {
		p.advance()
		elseBody, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		stmt.ElseBody = elseBody
	}

	return stmt, nil
}

func (p *Parser) parseFor() (Node, error) {
	line := p.current().Line
	p.advance() // skip "dari"

	// All forms require parentheses: dari (...) { }
	if err := p.expect(lexer.TOKEN_LPAREN, "("); err != nil {
		return nil, err
	}

	// Check if next token is an identifier (for/foreach) or an expression (while)
	if p.check(lexer.TOKEN_IDENTIFIER) {
		// Peek ahead to determine which form
		varName := p.current().Value

		// Check for Form 2: dari (item dalam collection) { } — foreach-style
		savedPos := p.pos
		p.advance()
		if p.check(lexer.TOKEN_DALAM) {
			p.advance()
			collection, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
				return nil, err
			}
			body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			return &ForEachStatement{VarName: varName, Collection: collection, Body: body, Line: line}, nil
		}

		// Check for Form 3: dari (i = start sampai end) { } — range-style
		if p.check(lexer.TOKEN_ASSIGN) {
			p.advance()
			start, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.TOKEN_SAMPAI, "sampai"); err != nil {
				return nil, err
			}
			end, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
				return nil, err
			}
			body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			return &ForStatement{VarName: varName, Start: start, End: end, Body: body, Line: line}, nil
		}

		// Not foreach or range — backtrack and parse as while-style condition
		p.pos = savedPos
	}

	// Form 1: dari (condition) { } — while-style
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return &WhileStatement{Condition: condition, Body: body, Line: line}, nil
}

func (p *Parser) parseFunctionDeclaration() (Node, error) {
	line := p.current().Line
	p.advance() // skip "buat kerja"

	if !p.check(lexer.TOKEN_IDENTIFIER) {
		// Anonymous function
		return p.parseAnonFunction(line)
	}

	name := p.current().Value
	p.advance()

	if err := p.expect(lexer.TOKEN_LPAREN, "("); err != nil {
		return nil, err
	}

	params, err := p.parseParamList()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
		return nil, err
	}

	returnType := p.parseOptionalType()

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &FunctionDeclaration{Name: name, Params: params, ReturnType: returnType, Body: body, Line: line}, nil
}

func (p *Parser) parseAnonFunction(line int) (Node, error) {
	if err := p.expect(lexer.TOKEN_LPAREN, "("); err != nil {
		return nil, err
	}

	params, err := p.parseParamList()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
		return nil, err
	}

	returnType := p.parseOptionalType()

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &AnonFunction{Params: params, ReturnType: returnType, Body: body, Line: line}, nil
}

func (p *Parser) parseParamList() ([]Param, error) {
	var params []Param
	if p.check(lexer.TOKEN_RPAREN) {
		return params, nil
	}

	if !p.check(lexer.TOKEN_IDENTIFIER) {
		return nil, p.error("expected parameter name")
	}
	name := p.current().Value
	p.advance()
	typeName := p.parseOptionalType()
	params = append(params, Param{Name: name, TypeName: typeName})

	for p.check(lexer.TOKEN_COMMA) {
		p.advance()
		if !p.check(lexer.TOKEN_IDENTIFIER) {
			return nil, p.error("expected parameter name after ','")
		}
		name = p.current().Value
		p.advance()
		typeName = p.parseOptionalType()
		params = append(params, Param{Name: name, TypeName: typeName})
	}

	return params, nil
}

func (p *Parser) parseReturn() (Node, error) {
	line := p.current().Line
	p.advance() // skip "bagi balik"

	// Check if there's a value to return
	if p.check(lexer.TOKEN_SEMICOLON) || p.check(lexer.TOKEN_RBRACE) || p.check(lexer.TOKEN_EOF) {
		p.skipSemicolon()
		return &ReturnStatement{Value: &NullLiteral{}, Line: line}, nil
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipSemicolon()
	return &ReturnStatement{Value: value, Line: line}, nil
}

func (p *Parser) parseTryCatch() (Node, error) {
	line := p.current().Line
	p.advance() // skip "cuba lah"

	tryBody, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	stmt := &TryCatchStatement{TryBody: tryBody, Line: line}

	// Parse "kantoi" (catch)
	if p.check(lexer.TOKEN_ALAMAK) {
		p.advance()
		if err := p.expect(lexer.TOKEN_LPAREN, "("); err != nil {
			return nil, err
		}
		if !p.check(lexer.TOKEN_IDENTIFIER) {
			return nil, p.error("expected error variable name")
		}
		stmt.CatchVar = p.current().Value
		p.advance()
		if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
			return nil, err
		}
		catchBody, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		stmt.CatchBody = catchBody
	}

	// Parse "mesti buat" (finally)
	if p.check(lexer.TOKEN_CONFIRM_NEED) {
		p.advance()
		finallyBody, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		stmt.FinallyBody = finallyBody
	}

	return stmt, nil
}

func (p *Parser) parseThrow() (Node, error) {
	line := p.current().Line
	p.advance() // skip "mampus"

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.skipSemicolon()
	return &ThrowStatement{Value: value, Line: line}, nil
}

func (p *Parser) parseExpressionStatement() (Node, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	p.skipSemicolon()
	return &ExpressionStatement{Expression: expr}, nil
}

func (p *Parser) parseBlock() ([]Node, error) {
	if err := p.expect(lexer.TOKEN_LBRACE, "{"); err != nil {
		return nil, err
	}

	var stmts []Node
	for !p.check(lexer.TOKEN_RBRACE) && !p.check(lexer.TOKEN_EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	if err := p.expect(lexer.TOKEN_RBRACE, "}"); err != nil {
		return nil, err
	}

	return stmts, nil
}

// Expression parsing (precedence climbing)

func (p *Parser) parseExpression() (Node, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (Node, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.check(lexer.TOKEN_IF_NOT) {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: "if not", Right: right}
	}

	return left, nil
}

func (p *Parser) parseAnd() (Node, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.check(lexer.TOKEN_ALSO) {
		p.advance()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: "also", Right: right}
	}

	return left, nil
}

func (p *Parser) parseEquality() (Node, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.check(lexer.TOKEN_SAMA) || p.check(lexer.TOKEN_TAK_SAMA) {
		op := p.current().Value
		p.advance()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseComparison() (Node, error) {
	left, err := p.parseAddition()
	if err != nil {
		return nil, err
	}

	for p.check(lexer.TOKEN_LESS) || p.check(lexer.TOKEN_GREATER) ||
		p.check(lexer.TOKEN_LESS_EQUAL) || p.check(lexer.TOKEN_GREATER_EQUAL) {
		op := p.current().Value
		p.advance()
		right, err := p.parseAddition()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseAddition() (Node, error) {
	left, err := p.parseMultiplication()
	if err != nil {
		return nil, err
	}

	for p.check(lexer.TOKEN_PLUS) || p.check(lexer.TOKEN_MINUS) {
		op := p.current().Value
		p.advance()
		right, err := p.parseMultiplication()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseMultiplication() (Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.check(lexer.TOKEN_STAR) || p.check(lexer.TOKEN_SLASH) || p.check(lexer.TOKEN_PERCENT) {
		op := p.current().Value
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseUnary() (Node, error) {
	if p.check(lexer.TOKEN_MINUS) || p.check(lexer.TOKEN_BANG) {
		op := p.current().Value
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Operator: op, Operand: operand}, nil
	}
	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (Node, error) {
	node, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		if p.check(lexer.TOKEN_LPAREN) {
			// Function call
			p.advance()
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
				return nil, err
			}
			node = &CallExpr{Callee: node, Args: args, Line: p.current().Line}
		} else if p.check(lexer.TOKEN_LBRACKET) {
			// Index access
			p.advance()
			index, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.TOKEN_RBRACKET, "]"); err != nil {
				return nil, err
			}
			node = &IndexExpr{Object: node, Index: index}
		} else if p.check(lexer.TOKEN_DOT) {
			// Member access
			p.advance()
			if !p.check(lexer.TOKEN_IDENTIFIER) {
				return nil, p.error("expected property name after '.'")
			}
			prop := p.current().Value
			p.advance()
			node = &MemberExpr{Object: node, Property: prop}
		} else {
			break
		}
	}

	return node, nil
}

func (p *Parser) parseArgList() ([]Node, error) {
	var args []Node
	if p.check(lexer.TOKEN_RPAREN) {
		return args, nil
	}

	arg, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	args = append(args, arg)

	for p.check(lexer.TOKEN_COMMA) {
		p.advance()
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return args, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	tok := p.current()

	switch tok.Type {
	case lexer.TOKEN_NUMBER:
		p.advance()
		var val float64
		fmt.Sscanf(tok.Value, "%f", &val)
		return &NumberLiteral{Value: val}, nil

	case lexer.TOKEN_STRING:
		p.advance()
		return &StringLiteral{Value: tok.Value}, nil

	case lexer.TOKEN_CAN:
		p.advance()
		return &BoolLiteral{Value: true}, nil

	case lexer.TOKEN_CANNOT:
		p.advance()
		return &BoolLiteral{Value: false}, nil

	case lexer.TOKEN_KOSONG:
		p.advance()
		return &NullLiteral{}, nil

	case lexer.TOKEN_IDENTIFIER:
		p.advance()
		return &Identifier{Name: tok.Value, Line: tok.Line}, nil

	case lexer.TOKEN_LPAREN:
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.TOKEN_RPAREN, ")"); err != nil {
			return nil, err
		}
		return expr, nil

	case lexer.TOKEN_LBRACKET:
		return p.parseArrayLiteral()

	case lexer.TOKEN_JOM:
		// Anonymous function in expression position
		line := tok.Line
		p.advance()
		return p.parseAnonFunction(line)

	case lexer.TOKEN_LAI:
		line := tok.Line
		p.advance()
		return &ReadStdinExpr{Line: line}, nil

	default:
		return nil, p.error(fmt.Sprintf("unexpected token: %s", tok.Value))
	}
}

func (p *Parser) parseArrayLiteral() (Node, error) {
	p.advance() // skip [

	var elements []Node
	if !p.check(lexer.TOKEN_RBRACKET) {
		elem, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)

		for p.check(lexer.TOKEN_COMMA) {
			p.advance()
			if p.check(lexer.TOKEN_RBRACKET) {
				break // trailing comma
			}
			elem, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
		}
	}

	if err := p.expect(lexer.TOKEN_RBRACKET, "]"); err != nil {
		return nil, err
	}

	return &ArrayLiteral{Elements: elements}, nil
}

// Helper methods

func (p *Parser) current() lexer.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return lexer.Token{Type: lexer.TOKEN_EOF}
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) check(tt lexer.TokenType) bool {
	return p.current().Type == tt
}

func (p *Parser) expect(tt lexer.TokenType, name string) error {
	if p.current().Type != tt {
		return p.error(fmt.Sprintf("expected '%s', got '%s'", name, p.current().Value))
	}
	p.advance()
	return nil
}

func (p *Parser) skipSemicolon() {
	if p.check(lexer.TOKEN_SEMICOLON) {
		p.advance()
	}
}

var validTypes = map[string]bool{
	"number":   true,
	"string":   true,
	"bool":     true,
	"array":    true,
	"function": true,
}

func (p *Parser) parseOptionalType() string {
	if !p.check(lexer.TOKEN_COLON) {
		return ""
	}
	p.advance() // skip ':'
	if !p.check(lexer.TOKEN_IDENTIFIER) {
		return ""
	}
	name := p.current().Value
	if !validTypes[name] {
		return ""
	}
	p.advance()
	return name
}

func (p *Parser) error(msg string) error {
	tok := p.current()
	return fmt.Errorf("line %d, col %d: %s", tok.Line, tok.Col, msg)
}
