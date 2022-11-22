package symbolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	customKeywords := map[string]TokenKind{"classes": -10, "age": -11}

	tests := []struct {
		input          string
		standardOutput []Token
		noSpaceOutput  []Token
		customOutput   []Token
	}{
		{
			"hello123,^",
			[]Token{
				{TokenIdentifier, "hello123"},
				{TokenKind(','), ","},
				{TokenKind('^'), "^"},
				EOFToken(),
			},
			[]Token{
				{TokenIdentifier, "hello123"},
				{TokenKind(','), ","},
				{TokenKind('^'), "^"},
				EOFToken(),
			},
			[]Token{
				{TokenIdentifier, "hello123"},
				{TokenKind(','), ","},
				{TokenKind('^'), "^"},
				EOFToken(),
			},
		},
		{
			"classes:: MyClass",
			[]Token{
				{TokenIdentifier, "classes"},
				{TokenKind(':'), ":"},
				{TokenKind(':'), ":"},
				UnicodeToken(' '),
				{TokenIdentifier, "MyClass"},
				EOFToken(),
			},
			[]Token{
				{TokenIdentifier, "classes"},
				{TokenKind(':'), ":"},
				{TokenKind(':'), ":"},
				{TokenIdentifier, "MyClass"},
				EOFToken(),
			},
			[]Token{
				{-10, "classes"},
				{TokenKind(':'), ":"},
				{TokenKind(':'), ":"},
				UnicodeToken(' '),
				{TokenIdentifier, "MyClass"},
				EOFToken(),
			},
		},
		{
			`"this is the text" -> "hello"`,
			[]Token{
				{TokenString, "this is the text"},
				UnicodeToken(' '),
				UnicodeToken('-'),
				UnicodeToken('>'),
				UnicodeToken(' '),
				{TokenString, "hello"},
				EOFToken(),
			},
			[]Token{
				{TokenString, "this is the text"},
				UnicodeToken('-'),
				UnicodeToken('>'),
				{TokenString, "hello"},
				EOFToken(),
			},
			[]Token{
				{TokenString, "this is the text"},
				UnicodeToken(' '),
				UnicodeToken('-'),
				UnicodeToken('>'),
				UnicodeToken(' '),
				{TokenString, "hello"},
				EOFToken(),
			},
		},
		{
			"12345. 2231",
			[]Token{
				{TokenNumber, "12345"},
				UnicodeToken('.'),
				UnicodeToken(' '),
				{TokenNumber, "2231"},
				EOFToken(),
			},
			[]Token{
				{TokenNumber, "12345"},
				UnicodeToken('.'),
				{TokenNumber, "2231"},
				EOFToken(),
			},
			[]Token{
				{TokenNumber, "12345"},
				UnicodeToken('.'),
				UnicodeToken(' '),
				{TokenNumber, "2231"},
				EOFToken(),
			},
		},
		{
			"person.age = 0x18",
			[]Token{
				{TokenIdentifier, "person"},
				UnicodeToken('.'),
				{TokenIdentifier, "age"},
				UnicodeToken(' '),
				UnicodeToken('='),
				UnicodeToken(' '),
				{TokenHexNumber, "18"},
				EOFToken(),
			},
			[]Token{
				{TokenIdentifier, "person"},
				UnicodeToken('.'),
				{TokenIdentifier, "age"},
				UnicodeToken('='),
				{TokenHexNumber, "18"},
				EOFToken(),
			},
			[]Token{
				{TokenIdentifier, "person"},
				UnicodeToken('.'),
				{-11, "age"},
				UnicodeToken(' '),
				UnicodeToken('='),
				UnicodeToken(' '),
				{TokenHexNumber, "18"},
				EOFToken(),
			},
		},
	}

	t.Run("Standard Lexer", func(t *testing.T) {
		for _, test := range tests {
			lexer, err := NewLexer(test.input)

			assert.Nil(t, err)
			assert.Equal(t, test.standardOutput, lexer.Tokens())
		}
	})

	t.Run("No Spaces Lexer", func(t *testing.T) {
		for _, test := range tests {
			lexer, err := NewLexer(test.input, IgnoreWhitespaces())

			assert.Nil(t, err)
			assert.Equal(t, test.noSpaceOutput, lexer.Tokens())
		}
	})

	t.Run("Custom Keyword Lexer", func(t *testing.T) {
		for _, test := range tests {
			lexer, err := NewLexer(test.input, Keywords(customKeywords))

			assert.Nil(t, err)
			assert.Equal(t, test.customOutput, lexer.Tokens())
		}
	})
}
