package symbolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Split(t *testing.T) {
	tests := []struct {
		inputs  string
		options []ParserOption
		delim   TokenKind
		output  []string
	}{
		{
			"hello-goodbye-aurevoir", nil,
			'-', []string{"hello", "goodbye", "aurevoir"},
		},
		{
			"github.com/manishmeganathan/symbolizer", nil,
			'/', []string{"github.com", "manishmeganathan", "symbolizer"},
		},
		{
			"aaron delim grayl", []ParserOption{Keywords(map[string]TokenKind{"delim": -10}), IgnoreWhitespaces()},
			-10, []string{"aaron", "grayl"},
		},
	}

	for _, test := range tests {
		parser := NewParser(test.inputs, test.options...)

		splits := parser.Split(test.delim)
		assert.Equal(t, test.output, splits)
	}
}
