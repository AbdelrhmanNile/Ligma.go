package token

// Token is a string that represents a lexical token.
type TokenType string;

type Token struct {
	Type   TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"
	FLOAT = "FLOAT"
	STRING = "STRING"
	BOOL = "BOOL"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"
	MOD	  	 = "%"
	POW		 = "**"

	LT = "<"
	GT = ">"
	EQ = "=="
	NOT_EQ = "!="
	GTE = ">="
	LTE = "<="

	// Logical Operators
	AND = "and"
	OR  = "or"
	NOT = "not"
	
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON	 = ":"
	DOT = "."

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	DEF     = "DEF"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	NULL	 = "NULL"
	CLASS = "CLASS"
	SELF = "SELF"
	SUPER = "SUPER"
	IMPORT = "IMPORT"

	// Control Flow
	FOR = "FOR"
	WHILE = "WHILE"
)


var keywords = map[string]TokenType{
	"def": DEF,
	"func":  FUNCTION,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
	"for": FOR,
	"while": WHILE,
	"and": AND,
	"or": OR,
	"not": NOT,
	"null": NULL,
	"class": CLASS,
	"self": SELF,
	"super": SUPER,
	"import": IMPORT,
}

// LookupIdent checks if a given identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}