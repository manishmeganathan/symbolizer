package symbolizer

// TokenKind is an enum for representing token grouping/values.
// For unicode tokens, the TokenKind is equal to its code point value.
// For literal such identifiers and numerics, the TokenKind values descend from 0.
// Note: Custom TokenKind values can be used by external packages for keyword detection
// for special literals, but these values should be below -10 to prevent collisions
type TokenKind int32

const (
	TokenEoF TokenKind = -(iota + 1)
	TokenIdentifier
	TokenNumber
	TokenString
	TokenHexNumber
)

// Token represents a lexical Token.
// It may be either a lone unicode character or some literal value
type Token struct {
	Kind     TokenKind
	Literal  string
	Position int
}

// UnicodeToken returns a Token for a given rune character.
// The TokenKind of the returned Token has the same value as it's unicode code point.
func UnicodeToken(char rune, pos int) Token {
	return Token{TokenKind(char), string(char), pos}
}

// EOFToken returns an End of File Token
func EOFToken(pos int) Token {
	return Token{TokenEoF, "", pos}
}
