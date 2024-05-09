package lexer

import "ligma/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination

	line int

	// for error handling
	errors []string
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	line := 1
	l := &Lexer{input: input, line: line}
	l.readChar()
	return l
}

// readChar reads the next character in the input and advances the position in the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { // check if we've reached the end of the input
		l.ch = 0 // ASCII code for "NUL" (null terminator)
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// NextToken returns the next token in the input
func (l *Lexer) NextToken() token.Token {
	
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
		case '+':
			tok = newTokenChar(token.PLUS, l.ch)
		case '-':
			tok = newTokenChar(token.MINUS, l.ch)
		case '*':
			if l.peekChar() == '*' {
				ch := l.ch
				l.readChar()
				tok = newTokenStr(token.POW, string(ch) + string(l.ch))
			} else {
				tok = newTokenChar(token.ASTERISK, l.ch)
			}
		case '/':
			tok = newTokenChar(token.SLASH, l.ch)
		case '%':
			tok = newTokenChar(token.MOD, l.ch)
		case '!':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = newTokenStr(token.NOT_EQ, string(ch) + string(l.ch))
			} else {
				tok = newTokenChar(token.BANG, l.ch)
			}
		case '<':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = newTokenStr(token.LTE, string(ch) + string(l.ch))
			} else {
				tok = newTokenChar(token.LT, l.ch)
			}
		case '>':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = newTokenStr(token.GTE, string(ch) + string(l.ch))
			} else {
				tok = newTokenChar(token.GT, l.ch)
			}
		case '=':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = newTokenStr(token.EQ, string(ch) + string(l.ch))
			} else {
				tok = newTokenChar(token.ASSIGN, l.ch)
			}
		case ';':
			tok = newTokenChar(token.SEMICOLON, l.ch)
		case ':':
			tok = newTokenChar(token.COLON, l.ch)
		case ',':
			tok = newTokenChar(token.COMMA, l.ch)
		case '.':
			tok = newTokenChar(token.DOT, l.ch)
		case '(':
			tok = newTokenChar(token.LPAREN, l.ch)
		case ')':
			tok = newTokenChar(token.RPAREN, l.ch)
		case '{':
			tok = newTokenChar(token.LBRACE, l.ch)
		case '}':
			tok = newTokenChar(token.RBRACE, l.ch)
		case '[':
			tok = newTokenChar(token.LBRACKET, l.ch)
		case ']':
			tok = newTokenChar(token.RBRACKET, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
		default:
			if isLetter(l.ch) {
				tok.Literal = l.readIdentifier()
				tok.Type = token.LookupIdent(tok.Literal)
				return tok
			} else if isDigit(l.ch) {
				return l.readNumber()
			} else if l.ch == '"' {
				tok.Type = token.STRING
				tok.Literal = l.readString()
				return tok
			} else {
				tok = newTokenChar(token.ILLEGAL, l.ch)
			}
	}
	l.readChar()
	return tok
}

// readIdentifier reads an identifier from the input
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isAlphanumeric(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// readNumber reads a number from the input
func (l *Lexer) readNumber() token.Token {
	position := l.position
	dotCount := 0

	for isDigit(l.ch) || l.ch == '.' {
		
		if l.ch == '.' {
			dotCount++
		}

		if dotCount > 1 { // if there are more than one dots in the number, it's invalid
			invalidNumber := l.input[position:l.position]

			for l.ch != ' ' && l.ch != '\t' && l.ch != '\n' && l.ch != '\r' && l.ch != ';' {
				invalidNumber += string(l.ch)
				l.readChar()
			}

			tok := newTokenStr(token.ILLEGAL, invalidNumber)
			return tok

		}

		l.readChar()
	}

	if dotCount == 1 {
		return newTokenStr(token.FLOAT, l.input[position:l.position])
	}

	return newTokenStr(token.INT, l.input[position:l.position])
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			l.readChar()
			break
		}
	}
	return l.input[position:l.position - 1]
}

// skipWhitespace skips any whitespace characters in the input
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

// peekChar returns the next character in the input without advancing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// newToken creates a new token with the given type and literal
func newTokenChar(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newTokenStr(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

// isLetter checks if a given character is a letter or an underscore
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if a given character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlphanumeric(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}