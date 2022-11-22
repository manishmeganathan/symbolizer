package symbolizer

// parseConfig is an internal configuration object for the
// lexer/parser that are modified using ParserOption functions
type parseConfig struct {
	eatSpaces bool
	keywords  map[string]TokenKind
}

// ParserOption represents an option to modify the Parser behaviour.
// It must be provided with the constructor for Parser.
type ParserOption func(config *parseConfig)

// Keywords returns a LexerOption that can be used to provide the Lexer
// with a set of special keywords mapped to some custom TokenKind value.
// If the Lexer encounters identifiers that match any of the given keywords,
// It returns a Token with the given kind and the actual literal encountered.
//
// Note: Use TokenKind values less than -10 for custom Token classes.
// -10 to -1 are reserved for standard token classes while 0 and above correspond the unicode code points.
func Keywords(keywords map[string]TokenKind) ParserOption {
	return func(config *parseConfig) {
		config.keywords = keywords
	}
}

// IgnoreWhitespaces returns a LexerOption that specifies the Lexer to ignore unicode characters with the
// whitespace property (' ', '\t', '\n', '\r', etc). They are consumed instead of generating Tokens for them.
func IgnoreWhitespaces() ParserOption {
	return func(config *parseConfig) {
		config.eatSpaces = true
	}
}
