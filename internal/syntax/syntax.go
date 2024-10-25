package syntax

type Symbol string

const (
	Bang      Symbol = "!"
	Colon     Symbol = ":"
	SemiColon Symbol = ";"
	LBrace    Symbol = "{"
	RBrace    Symbol = "}"
	LPara     Symbol = "("
	RPara     Symbol = ")"
)

type Keyword string

const (
	// TypeScript
	Export    Keyword = "export"
	Interface Keyword = "interface"
	Enum      Keyword = "enum"

	// GraphQL
	Input    Keyword = "input"
	Type     Keyword = "type"
	Mutation Keyword = "Mutation"
	Queries  Keyword = "Queries"
)
