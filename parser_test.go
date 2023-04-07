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
		{
			"0x1888, -122452", []ParserOption{IgnoreWhitespaces()},
			',', []string{"0x1888", "-122452"},
		},
	}

	for _, test := range tests {
		parser := NewParser(test.inputs, test.options...)
		splits := parser.Split(test.delim)
		assert.Equal(t, test.output, splits)
	}
}

func TestParser_Unwrap(t *testing.T) {
	// mustEnclose is a helper function
	mustEnclose := func(enclosure Enclosure, err error) Enclosure {
		if err != nil {
			panic(err)
		}

		return enclosure
	}

	tests := []struct {
		input    string
		options  []ParserOption
		enclose  Enclosure
		output   string
		error    string
		unparsed string
	}{
		{
			"[string]", nil, EnclosureSquare(),
			"string", "", "",
		},
		{
			"{map[string]string}", nil, EnclosureCurly(),
			"map[string]string", "", "",
		},
		{
			"< mycroft  holmes >", []ParserOption{IgnoreWhitespaces()}, EnclosureAngle(),
			" mycroft  holmes ", "", "",
		},
		{
			"@sarah[chapman&", nil, mustEnclose(NewEnclosure('@', '&')),
			"sarah[chapman", "", "",
		},
		{
			"( 12345(555))hello123", []ParserOption{IgnoreWhitespaces()}, EnclosureParens(),
			" 12345(555)", "", "hello123",
		},
		{
			"map(sequence[map])", nil, EnclosureParens(),
			"", "missing start of enclosure: '('", "map(sequence[map])",
		},
		{
			"(map(sequence[map]", nil, EnclosureParens(),
			"", "missing end of enclosure: ')'", "",
		},
	}

	for _, test := range tests {
		parser := NewParser(test.input, test.options...)
		unwrapped, err := parser.Unwrap(test.enclose)
		assert.Equal(t, test.output, unwrapped, "Unwrapped Data Check")

		if test.error != "" {
			assert.EqualError(t, err, test.error, "Error Check")
		}

		assert.Equal(t, test.unparsed, parser.Unparsed(), "Unparsed Data Check")
	}
}

func TestParser_Unparsed(t *testing.T) {
	tests := []struct {
		input    string
		advances int
		output   string
	}{
		{"string&&string", 3, "string"},
		{"map[string]string", 4, "string"},
		{"[32]uint64", 3, "uint64"},
		{"[1024]map[string]string", 3, "map[string]string"},
	}

	for _, test := range tests {
		parser := NewParser(test.input)
		for i := 0; i < test.advances; i++ {
			parser.Advance()
		}

		assert.Equal(t, test.output, parser.Unparsed())
	}
}

func TestParser_Peeking(t *testing.T) {
	tests := []struct {
		input     string
		options   []ParserOption
		advances  int
		tryPeek   TokenKind
		tryResult bool
	}{
		{
			"string&&string", nil, 2,
			TokenIdentifier, true,
		},
		{
			"[32]uint64", nil, 1,
			TokenKind('['), false,
		},
	}

	for _, test := range tests {
		parser := NewParser(test.input, test.options...)
		for i := 0; i < test.advances; i++ {
			parser.Advance()
		}

		peek := parser.Peek()
		assert.True(t, parser.IsPeek(peek.Kind))

		tryPeek := parser.ExpectPeek(test.tryPeek)
		assert.Equal(t, test.tryResult, tryPeek)
		if test.tryResult {
			assert.Equal(t, peek, parser.Cursor())
		}
	}
}
