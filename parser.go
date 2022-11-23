package symbolizer

import "fmt"

// Parser is a symbol parser that parse a given string input and handle
// operations like unwrapping enclosed data or splitting by a given delimiter
type Parser struct {
	// scanner represents the token scanner
	scanner *lexer
	// // trail represents a trailing cursor indicating the beginning of the current Token
	// trail int
	// curr and next represent the current and next Token values
	curr, next Token
}

// NewParser generates a new Parser for a given input string and some options that
// modify the parser behaviour such as ignoring whitespaces or using custom keywords
func NewParser(input string, opts ...ParserOption) *Parser {
	// Create a new parseConfig and apply all the given options on it
	config := new(parseConfig)
	for _, option := range opts {
		option(config)
	}

	// Create a token scanning lexer for the input
	scanner := &lexer{config: config, symbols: []rune(input)}
	parser := &Parser{scanner: scanner}

	// Advance the parser twice to initialize
	// the curr and next Tokens of the parser
	parser.Advance()
	parser.Advance()

	return parser
}

// Peek looks ahead and returns the next Token without advancing the parser
func (parser *Parser) Peek() Token { return parser.next }

// Cursor returns the current Token
func (parser *Parser) Cursor() Token { return parser.curr }

// Unparsed returns the remaining unparsed data in the parser as a string
func (parser *Parser) Unparsed() string {
	return string(parser.scanner.symbols)[parser.curr.Position:]
}

// Advance moves the parser's cursor and peek tokens
func (parser *Parser) Advance() {
	parser.curr = parser.next
	parser.next = parser.scanner.next()
}

// IsPeek checks if the next token is of the specified TokenKind.
// This look ahead is performed without moving the parser's cursor
func (parser *Parser) IsPeek(t TokenKind) bool {
	return parser.next.Kind == t
}

// IsCursor checks if the current token is of the specified TokenKind.
func (parser *Parser) IsCursor(t TokenKind) bool {
	return parser.curr.Kind == t
}

// ExpectPeek advances the cursor if the next token is of the specified TokenKind.
// If it is not the same type, the parser does not advance.
// The returned boolean indicates if the parser was advanced.
func (parser *Parser) ExpectPeek(t TokenKind) bool {
	// Check if peek token matches
	if !parser.IsPeek(t) {
		return false
	}

	// Advance the parse cursor
	parser.Advance()

	return true
}

// Split attempts to split the remaining contents of the parser
// into a set of strings separated by the given delimiting TokenKind.
// This process exhausts the parser consuming all the tokens within it.
func (parser *Parser) Split(delimiter TokenKind) (splits []string) {
	var accumulator string

Loop:
	for {
		switch parser.Cursor().Kind {
		case delimiter:
			// Append the accumulated characters and reset the accumulator
			splits = append(splits, accumulator)
			accumulator = ""

		case TokenEoF:
			// Append accumulated characters
			splits = append(splits, accumulator)
			// Break from loop (end of symbol)
			break Loop

		default:
			// Accumulate character
			accumulator += parser.curr.Literal
		}

		parser.Advance()
	}

	return splits
}

// Unwrap attempts to unravel a substring enclosed between to characters described with an Enclosure.
// When calling Unwrap, the parse cursor must be the opening character of the given Enclosure. Returns
// an error if the opening character is not found or if the symbol terminates before the closing character.
//
// Note: Unwrap will resolve nested enclosures attempting to match one
// opening character with one closing character until it fully resolves.
func (parser *Parser) Unwrap(enc Enclosure) (string, error) {
	// Require the current token of the parser to be the enclosure opening token
	if !parser.IsCursor(TokenKind(enc.start)) {
		return "", fmt.Errorf("missing start of enclosure: '%v'", string(enc.start))
	}

	// Record the start of the enclosed data (1 position after enclose opener)
	start := parser.curr.Position + 1
	// First enclose opener sets the nesting level to 1.
	// This nesting level needs to be resolved for the enclosure to "end"
	nesting := 1

	// Advance the cursor into the enclosed data.
	parser.Advance()

	for {
		switch parser.Cursor().Kind {
		case TokenKind(enc.start):
			// Increase nesting level, if new enclosure start is encountered
			nesting++
		case TokenKind(enc.stop):
			// Reduce nesting level, if new enclosure end is encountered
			nesting--

		case TokenEoF:
			// premature end of symbol
			return "", fmt.Errorf("missing end of enclosure: '%v'", string(enc.stop))
		}

		if nesting == 0 {
			// If nesting is resolved, slice input and return
			return parser.scanner.collectBetween(start, parser.curr.Position), nil
		}

		parser.Advance()
	}
}
