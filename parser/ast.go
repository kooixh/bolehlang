package parser

// AST Node types

type Node interface {
	nodeType() string
}

type Program struct {
	Body []Node
}

func (p *Program) nodeType() string { return "Program" }

type ImportStatement struct {
	Path string
	Line int
}

func (i *ImportStatement) nodeType() string { return "ImportStatement" }

type Param struct {
	Name     string
	TypeName string // "" means no type annotation
}

type VarDeclaration struct {
	Name     string
	TypeName string // "" means no type annotation
	Value    Node
	Line     int
}

func (v *VarDeclaration) nodeType() string { return "VarDeclaration" }

type ConstDeclaration struct {
	Name     string
	TypeName string // "" means no type annotation
	Value    Node
	Line     int
}

func (c *ConstDeclaration) nodeType() string { return "ConstDeclaration" }

type Assignment struct {
	Name     string
	Value    Node
	Operator string // "=" or "+=" or "-="
	Line     int
}

func (a *Assignment) nodeType() string { return "Assignment" }

type PrintStatement struct {
	Value Node
	Line  int
}

func (p *PrintStatement) nodeType() string { return "PrintStatement" }

type WarnStatement struct {
	Value Node
	Line  int
}

func (w *WarnStatement) nodeType() string { return "WarnStatement" }

type IfStatement struct {
	Condition   Node
	Body        []Node
	ElseIfCases []ElseIfCase
	ElseBody    []Node
	Line        int
}

type ElseIfCase struct {
	Condition Node
	Body      []Node
}

func (i *IfStatement) nodeType() string { return "IfStatement" }

type WhileStatement struct {
	Condition Node
	Body      []Node
	Line      int
}

func (w *WhileStatement) nodeType() string { return "WhileStatement" }

type ForStatement struct {
	VarName string
	Start   Node
	End     Node
	Body    []Node
	Line    int
}

func (f *ForStatement) nodeType() string { return "ForStatement" }

type ForEachStatement struct {
	VarName    string
	Collection Node
	Body       []Node
	Line       int
}

func (f *ForEachStatement) nodeType() string { return "ForEachStatement" }

type FunctionDeclaration struct {
	Name       string
	Params     []Param
	ReturnType string // "" means no type annotation
	Body       []Node
	Line       int
}

func (f *FunctionDeclaration) nodeType() string { return "FunctionDeclaration" }

type ReturnStatement struct {
	Value Node
	Line  int
}

func (r *ReturnStatement) nodeType() string { return "ReturnStatement" }

type BreakStatement struct {
	Line int
}

func (b *BreakStatement) nodeType() string { return "BreakStatement" }

type ContinueStatement struct {
	Line int
}

func (c *ContinueStatement) nodeType() string { return "ContinueStatement" }

type TryCatchStatement struct {
	TryBody     []Node
	CatchVar    string
	CatchBody   []Node
	FinallyBody []Node
	Line        int
}

func (t *TryCatchStatement) nodeType() string { return "TryCatchStatement" }

type ThrowStatement struct {
	Value Node
	Line  int
}

func (t *ThrowStatement) nodeType() string { return "ThrowStatement" }

type ExpressionStatement struct {
	Expression Node
}

func (e *ExpressionStatement) nodeType() string { return "ExpressionStatement" }

// Expressions

type BinaryExpr struct {
	Left     Node
	Operator string
	Right    Node
}

func (b *BinaryExpr) nodeType() string { return "BinaryExpr" }

type UnaryExpr struct {
	Operator string
	Operand  Node
}

func (u *UnaryExpr) nodeType() string { return "UnaryExpr" }

type CallExpr struct {
	Callee Node
	Args   []Node
	Line   int
}

func (c *CallExpr) nodeType() string { return "CallExpr" }

type MemberExpr struct {
	Object   Node
	Property string
}

func (m *MemberExpr) nodeType() string { return "MemberExpr" }

type IndexExpr struct {
	Object Node
	Index  Node
}

func (i *IndexExpr) nodeType() string { return "IndexExpr" }

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) nodeType() string { return "NumberLiteral" }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) nodeType() string { return "StringLiteral" }

type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) nodeType() string { return "BoolLiteral" }

type NullLiteral struct{}

func (n *NullLiteral) nodeType() string { return "NullLiteral" }

type ArrayLiteral struct {
	Elements []Node
}

func (a *ArrayLiteral) nodeType() string { return "ArrayLiteral" }

type Identifier struct {
	Name string
	Line int
}

func (i *Identifier) nodeType() string { return "Identifier" }

type ReadStdinExpr struct {
	Line int
}

func (r *ReadStdinExpr) nodeType() string { return "ReadStdinExpr" }

type AnonFunction struct {
	Params     []Param
	ReturnType string // "" means no type annotation
	Body       []Node
	Line       int
}

func (a *AnonFunction) nodeType() string { return "AnonFunction" }

// structs declaration init

type StructDeclaration struct {
	Name string
	Fields []string
	Line int
}

func(k *StructDeclaration) nodeType() string { return "StructDeclaration"}
