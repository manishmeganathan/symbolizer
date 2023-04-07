package symbolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	customKeywords := map[string]TokenKind{"classes": -10, "age": -11, "mark": -12}

	tests := []struct {
		input          string
		standardOutput []Token
		noSpaceOutput  []Token
		customOutput   []Token
	}{
		{
			"hello123,^",
			[]Token{
				{TokenIdentifier, "hello123", 0},
				{TokenKind(','), ",", 8},
				{TokenKind('^'), "^", 9},
				EOFToken(10),
			},
			[]Token{
				{TokenIdentifier, "hello123", 0},
				{TokenKind(','), ",", 8},
				{TokenKind('^'), "^", 9},
				EOFToken(10),
			},
			[]Token{
				{TokenIdentifier, "hello123", 0},
				{TokenKind(','), ",", 8},
				{TokenKind('^'), "^", 9},
				EOFToken(10),
			},
		},
		{
			"classes:: MyClass",
			[]Token{
				{TokenIdentifier, "classes", 0},
				{TokenKind(':'), ":", 7},
				{TokenKind(':'), ":", 8},
				UnicodeToken(' ', 9),
				{TokenIdentifier, "MyClass", 10},
				EOFToken(17),
			},
			[]Token{
				{TokenIdentifier, "classes", 0},
				{TokenKind(':'), ":", 7},
				{TokenKind(':'), ":", 8},
				{TokenIdentifier, "MyClass", 10},
				EOFToken(17),
			},
			[]Token{
				{-10, "classes", 0},
				{TokenKind(':'), ":", 7},
				{TokenKind(':'), ":", 8},
				UnicodeToken(' ', 9),
				{TokenIdentifier, "MyClass", 10},
				EOFToken(17),
			},
		},
		{
			`"this is the text" -> "hello"`,
			[]Token{
				{TokenString, "this is the text", 0},
				UnicodeToken(' ', 18),
				UnicodeToken('-', 19),
				UnicodeToken('>', 20),
				UnicodeToken(' ', 21),
				{TokenString, "hello", 22},
				EOFToken(29),
			},
			[]Token{
				{TokenString, "this is the text", 0},
				UnicodeToken('-', 19),
				UnicodeToken('>', 20),
				{TokenString, "hello", 22},
				EOFToken(29),
			},
			[]Token{
				{TokenString, "this is the text", 0},
				UnicodeToken(' ', 18),
				UnicodeToken('-', 19),
				UnicodeToken('>', 20),
				UnicodeToken(' ', 21),
				{TokenString, "hello", 22},
				EOFToken(29),
			},
		},
		{
			"12345. 2231",
			[]Token{
				{TokenNumber, "12345", 0},
				UnicodeToken('.', 5),
				UnicodeToken(' ', 6),
				{TokenNumber, "2231", 7},
				EOFToken(11),
			},
			[]Token{
				{TokenNumber, "12345", 0},
				UnicodeToken('.', 5),
				{TokenNumber, "2231", 7},
				EOFToken(11),
			},
			[]Token{
				{TokenNumber, "12345", 0},
				UnicodeToken('.', 5),
				UnicodeToken(' ', 6),
				{TokenNumber, "2231", 7},
				EOFToken(11),
			},
		},
		{
			"person.age = 0x18",
			[]Token{
				{TokenIdentifier, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdentifier, "age", 7},
				UnicodeToken(' ', 10),
				UnicodeToken('=', 11),
				UnicodeToken(' ', 12),
				{TokenHexNumber, "0x18", 13},
				EOFToken(17),
			},
			[]Token{
				{TokenIdentifier, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdentifier, "age", 7},
				UnicodeToken('=', 11),
				{TokenHexNumber, "0x18", 13},
				EOFToken(17),
			},
			[]Token{
				{TokenIdentifier, "person", 0},
				UnicodeToken('.', 6),
				{-11, "age", 7},
				UnicodeToken(' ', 10),
				UnicodeToken('=', 11),
				UnicodeToken(' ', 12),
				{TokenHexNumber, "0x18", 13},
				EOFToken(17),
			},
		},
		{
			"person.mark = -923",
			[]Token{
				{TokenIdentifier, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdentifier, "mark", 7},
				UnicodeToken(' ', 11),
				UnicodeToken('=', 12),
				UnicodeToken(' ', 13),
				{TokenNumber, "-923", 14},
				EOFToken(18),
			},
			[]Token{
				{TokenIdentifier, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdentifier, "mark", 7},
				UnicodeToken('=', 12),
				{TokenNumber, "-923", 14},
				EOFToken(18),
			},
			[]Token{
				{TokenIdentifier, "person", 0},
				UnicodeToken('.', 6),
				{-12, "mark", 7},
				UnicodeToken(' ', 11),
				UnicodeToken('=', 12),
				UnicodeToken(' ', 13),
				{TokenNumber, "-923", 14},
				EOFToken(18),
			},
		},
	}

	t.Run("Standard Lexer", func(t *testing.T) {
		for _, test := range tests {
			lex := lexer{0, []rune(test.input), &parseConfig{eatSpaces: false, keywords: nil}}
			assert.Equal(t, test.standardOutput, lex.tokens())
		}
	})

	t.Run("No Spaces Lexer", func(t *testing.T) {
		for _, test := range tests {
			lex := lexer{0, []rune(test.input), &parseConfig{eatSpaces: true, keywords: nil}}
			assert.Equal(t, test.noSpaceOutput, lex.tokens())
		}
	})

	t.Run("Custom Keyword Lexer", func(t *testing.T) {
		for _, test := range tests {
			lex := lexer{0, []rune(test.input), &parseConfig{eatSpaces: false, keywords: customKeywords}}
			assert.Equal(t, test.customOutput, lex.tokens())
		}
	})
}

func TestTokenKind_String(t *testing.T) {
	tests := []struct {
		token  TokenKind
		output string
	}{
		{TokenKind('5'), "<unicode:'5'>"},
		{TokenKind('&'), "<unicode:'&'>"},
		{TokenKind(-10), "<custom:-10>"},
		{TokenEoF, "<eof>"},
		{TokenNumber, "<num>"},
		{TokenIdentifier, "<ident>"},
		{TokenString, "<str>"},
		{TokenHexNumber, "<hex>"},
	}

	for _, test := range tests {
		str := test.token.String()
		assert.Equal(t, test.output, str)
	}
}
