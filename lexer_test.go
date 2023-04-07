package symbolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	customKeywords := map[string]TokenKind{
		"classes": -10,
		"age":     -11,
		"mark":    -12,
		"True":    TokenBoolean,
	}

	tests := []struct {
		input          string
		standardOutput []Token
		noSpaceOutput  []Token
		customOutput   []Token
	}{
		{
			`hello123,^`,
			[]Token{
				{TokenIdent, "hello123", 0},
				{TokenKind(','), ",", 8},
				{TokenKind('^'), "^", 9},
				EOFToken(10),
			},
			[]Token{
				{TokenIdent, "hello123", 0},
				{TokenKind(','), ",", 8},
				{TokenKind('^'), "^", 9},
				EOFToken(10),
			},
			[]Token{
				{TokenIdent, "hello123", 0},
				{TokenKind(','), ",", 8},
				{TokenKind('^'), "^", 9},
				EOFToken(10),
			},
		},
		{
			`true = True`,
			[]Token{
				{TokenBoolean, "true", 0},
				UnicodeToken(' ', 4),
				{TokenKind('='), "=", 5},
				UnicodeToken(' ', 6),
				{TokenIdent, "True", 7},
				EOFToken(11),
			},
			[]Token{
				{TokenBoolean, "true", 0},
				{TokenKind('='), "=", 5},
				{TokenIdent, "True", 7},
				EOFToken(11),
			},
			[]Token{
				{TokenBoolean, "true", 0},
				UnicodeToken(' ', 4),
				{TokenKind('='), "=", 5},
				UnicodeToken(' ', 6),
				{TokenBoolean, "True", 7},
				EOFToken(11),
			},
		},
		{
			`classes:: MyClass`,
			[]Token{
				{TokenIdent, "classes", 0},
				{TokenKind(':'), ":", 7},
				{TokenKind(':'), ":", 8},
				UnicodeToken(' ', 9),
				{TokenIdent, "MyClass", 10},
				EOFToken(17),
			},
			[]Token{
				{TokenIdent, "classes", 0},
				{TokenKind(':'), ":", 7},
				{TokenKind(':'), ":", 8},
				{TokenIdent, "MyClass", 10},
				EOFToken(17),
			},
			[]Token{
				{-10, "classes", 0},
				{TokenKind(':'), ":", 7},
				{TokenKind(':'), ":", 8},
				UnicodeToken(' ', 9),
				{TokenIdent, "MyClass", 10},
				EOFToken(17),
			},
		},
		{
			`"this is the text" -> "hello"`,
			[]Token{
				{TokenString, `"this is the text"`, 0},
				UnicodeToken(' ', 18),
				UnicodeToken('-', 19),
				UnicodeToken('>', 20),
				UnicodeToken(' ', 21),
				{TokenString, `"hello"`, 22},
				EOFToken(29),
			},
			[]Token{
				{TokenString, `"this is the text"`, 0},
				UnicodeToken('-', 19),
				UnicodeToken('>', 20),
				{TokenString, `"hello"`, 22},
				EOFToken(29),
			},
			[]Token{
				{TokenString, `"this is the text"`, 0},
				UnicodeToken(' ', 18),
				UnicodeToken('-', 19),
				UnicodeToken('>', 20),
				UnicodeToken(' ', 21),
				{TokenString, `"hello"`, 22},
				EOFToken(29),
			},
		},
		{
			`12345. 2231`,
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
				{TokenIdent, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdent, "age", 7},
				UnicodeToken(' ', 10),
				UnicodeToken('=', 11),
				UnicodeToken(' ', 12),
				{TokenHexNumber, "0x18", 13},
				EOFToken(17),
			},
			[]Token{
				{TokenIdent, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdent, "age", 7},
				UnicodeToken('=', 11),
				{TokenHexNumber, "0x18", 13},
				EOFToken(17),
			},
			[]Token{
				{TokenIdent, "person", 0},
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
			`person.mark = -923`,
			[]Token{
				{TokenIdent, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdent, "mark", 7},
				UnicodeToken(' ', 11),
				UnicodeToken('=', 12),
				UnicodeToken(' ', 13),
				{TokenNumber, "-923", 14},
				EOFToken(18),
			},
			[]Token{
				{TokenIdent, "person", 0},
				UnicodeToken('.', 6),
				{TokenIdent, "mark", 7},
				UnicodeToken('=', 12),
				{TokenNumber, "-923", 14},
				EOFToken(18),
			},
			[]Token{
				{TokenIdent, "person", 0},
				UnicodeToken('.', 6),
				{-12, "mark", 7},
				UnicodeToken(' ', 11),
				UnicodeToken('=', 12),
				UnicodeToken(' ', 13),
				{TokenNumber, "-923", 14},
				EOFToken(18),
			},
		},

		{
			`"abcdefg`,
			[]Token{
				{TokenMalformed, `"abcdefg`, 0},
				EOFToken(8),
			},
			[]Token{
				{TokenMalformed, `"abcdefg`, 0},
				EOFToken(8),
			},
			[]Token{
				{TokenMalformed, `"abcdefg`, 0},
				EOFToken(8),
			},
		},
	}

	t.Run("Standard Lexer", func(t *testing.T) {
		for _, test := range tests {
			lex := lexer{0, []rune(test.input), newParseConfig()}
			assert.Equal(t, test.standardOutput, lex.tokens())
		}
	})

	t.Run("No Spaces Lexer", func(t *testing.T) {
		for _, test := range tests {
			lex := lexer{0, []rune(test.input), newParseConfig(IgnoreWhitespaces())}
			assert.Equal(t, test.noSpaceOutput, lex.tokens())
		}
	})

	t.Run("Custom Keyword Lexer", func(t *testing.T) {
		for _, test := range tests {
			lex := lexer{0, []rune(test.input), newParseConfig(Keywords(customKeywords))}
			assert.Equal(t, test.customOutput, lex.tokens())
		}
	})
}
