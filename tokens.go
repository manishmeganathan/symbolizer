package symbolizer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// TokenKind is an enum for representing token grouping/values.
// For unicode tokens, the TokenKind is equal to its code point value.
// For literal such identifiers and numerics, the TokenKind values descend from 0.
// Note: Custom TokenKind values can be used by external packages for keyword detection
// for special literals, but these values should be below -10 to prevent collisions
type TokenKind int32

const (
	TokenEoF TokenKind = -(iota + 1)
	TokenIdent
	TokenNumber
	TokenString
	TokenBoolean
	TokenHexNumber
)

// String implements the Stringer interface for TokenKind
func (kind TokenKind) String() string {
	if kind > 0 {
		return fmt.Sprintf("<unicode:'%v'>", string(kind))
	}

	switch kind {
	case TokenEoF:
		return "<eof>"
	case TokenIdent:
		return "<ident>"
	case TokenNumber:
		return "<num>"
	case TokenString:
		return "<str>"
	case TokenHexNumber:
		return "<hex>"
	default:
		return fmt.Sprintf("<custom:%d>", kind)
	}
}

// CanValue returns whether the TokenKind can be converted into a value
func (kind TokenKind) CanValue() bool {
	return kind == TokenNumber || kind == TokenString || kind == TokenBoolean || kind == TokenHexNumber
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

// Token represents a lexical Token.
// It may be either a lone unicode character or some literal value
type Token struct {
	Kind     TokenKind
	Literal  string
	Position int
}

// Value returns an object value for the Token.
// If the Token is kind TokenString -> string (literal is returned as is)
// If the Token is kind TokenBoolean -> bool (parsed with strconv.ParseBool)
// If the Token is kind TokenNumber -> uint64/int64 (parsed with strconv depending on if a negative sign is present)
// If the Token is kind TokenHexNumber -> []byte (decoded with hex.DecodeString after trimming the 0x)
// All other Token kinds will return an error if attempted to convert to values
func (token Token) Value() (any, error) {
	switch token.Kind {

	// String Value
	case TokenString:
		return token.Literal, nil

	// Boolean Value
	case TokenBoolean:
		boolean, err := strconv.ParseBool(token.Literal)
		if err != nil {
			return nil, errors.New("invalid boolean token: could not parse as boolean")
		}

		return boolean, nil

	// Hex Value
	case TokenHexNumber:
		data, err := hex.DecodeString(strings.TrimPrefix(token.Literal, "0x"))
		if err != nil {
			return nil, fmt.Errorf("invalid hex token: %w", err)
		}

		return data, nil

	// Numeric Value
	case TokenNumber:
		// Negative Number
		if strings.HasPrefix(token.Literal, "-") {
			number, err := strconv.ParseInt(token.Literal, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid signed numeric token: %w", err)
			}

			return number, nil
		}

		number, err := strconv.ParseUint(token.Literal, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid numeric token: %w", err)
		}

		return number, nil

	default:
		return nil, fmt.Errorf("cannot generate from value from token of kind '%v'", token.Kind)
	}
}

// Enclosure is a tuple of unicode code points that indicate
// start and stop pairs. They cannot be the same.
type Enclosure struct {
	start, stop rune
}

// NewEnclosure generates a new Enclosure set and returns it.
// Throws an error if the start and stop code points are identical
func NewEnclosure(start, stop rune) (Enclosure, error) {
	if start == stop {
		return Enclosure{}, errors.New("enclosure start and stop cannot be the same")
	}

	return Enclosure{start, stop}, nil
}

// EnclosureParens returns an Enclosure set for Parenthesis '()'
func EnclosureParens() Enclosure {
	return Enclosure{'(', ')'}
}

// EnclosureSquare returns an Enclosure set for Square Brackets '[]'
func EnclosureSquare() Enclosure {
	return Enclosure{'[', ']'}
}

// EnclosureCurly returns an Enclosure set for Curly Brackets '{}'
func EnclosureCurly() Enclosure {
	return Enclosure{'{', '}'}
}

// EnclosureAngle returns an Enclosure set for Angle Brackets '<>'
func EnclosureAngle() Enclosure {
	return Enclosure{'<', '>'}
}
