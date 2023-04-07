package symbolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		{TokenIdent, "<ident>"},
		{TokenString, "<str>"},
		{TokenHexNumber, "<hex>"},
		{TokenBoolean, "<bool>"},
		{TokenMalformed, "<malformed>"},
	}

	for _, test := range tests {
		str := test.token.String()
		assert.Equal(t, test.output, str)
	}
}

func TestTokenKind_CanValue(t *testing.T) {
	tests := []struct {
		token  TokenKind
		output bool
	}{
		{TokenKind('5'), false},
		{TokenKind('&'), false},
		{TokenKind(-10), false},
		{TokenEoF, false},
		{TokenNumber, true},
		{TokenIdent, false},
		{TokenString, true},
		{TokenHexNumber, true},
		{TokenBoolean, true},
		{TokenMalformed, false},
	}

	for _, test := range tests {
		str := test.token.CanValue()
		assert.Equal(t, test.output, str)
	}
}
func TestToken_Value(t *testing.T) {
	tests := []struct {
		token Token
		value any
		err   string
	}{
		{Token{Kind: TokenString, Literal: `"hello"`}, "hello", ""},
		{Token{Kind: TokenKind('-'), Literal: "-"}, nil, "cannot generate from value from token of kind '<unicode:'-'>'"},

		{Token{Kind: TokenBoolean, Literal: "true"}, true, ""},
		{Token{Kind: TokenBoolean, Literal: "TRUE"}, true, ""},
		{Token{Kind: TokenBoolean, Literal: "False"}, false, ""},
		{Token{Kind: TokenBoolean, Literal: "Quantum"}, nil, "invalid boolean token: could not parse as boolean"},

		{Token{Kind: TokenHexNumber, Literal: "0x23ab8492"}, []byte{0x23, 0xab, 0x84, 0x92}, ""},
		{Token{Kind: TokenHexNumber, Literal: "23ab8492"}, []byte{0x23, 0xab, 0x84, 0x92}, ""},
		{Token{Kind: TokenHexNumber, Literal: "23ab842"}, nil, "invalid hex token: encoding/hex: odd length hex string"},

		{Token{Kind: TokenNumber, Literal: "9328572352"}, uint64(9328572352), ""},
		{Token{Kind: TokenNumber, Literal: "-9223372036854775807"}, int64(-9223372036854775807), ""},
		{Token{Kind: TokenNumber, Literal: "18446744073709551615"}, uint64(18446744073709551615), ""},
		{Token{Kind: TokenNumber, Literal: "1844674407370955161523123"}, nil, "invalid numeric token: strconv.ParseUint: parsing \"1844674407370955161523123\": value out of range"},
		{Token{Kind: TokenNumber, Literal: "-18446744073709551615"}, nil, "invalid signed numeric token: strconv.ParseInt: parsing \"-18446744073709551615\": value out of range"},
	}

	for _, test := range tests {
		value, err := test.token.Value()

		if test.err == "" {
			require.NoError(t, err)
			require.Equal(t, test.value, value)
		} else {
			require.Nil(t, value)
			require.EqualError(t, err, test.err)
		}
	}
}
