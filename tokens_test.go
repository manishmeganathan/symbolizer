package symbolizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToken_Value(t *testing.T) {
	tests := []struct {
		token Token
		value any
		err   string
	}{
		{Token{Kind: TokenString, Literal: "hello"}, "hello", ""},
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
