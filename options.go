package symbolizer

// parseConfig is an internal configuration object for the
// lexer/parser that are modified using ParserOption functions
type parseConfig struct {
	eatSpaces bool
	keywords  map[string]TokenKind
}

// newParseConfig generate a new parseConfig with all default params
// and then applies any options provided to modify the config
func newParseConfig(opts ...ParserOption) *parseConfig {
	// Create a new parseConfig and apply all the given options on it
	config := new(parseConfig)
	// Set the default keywords
	config.keywords = map[string]TokenKind{
		"true":  TokenBoolean,
		"false": TokenBoolean,
	}

	for _, option := range opts {
		option(config)
	}

	return config
}

// ParserOption represents an option to modify the Parser behaviour.
// It must be provided with the constructor for Parser.
type ParserOption func(config *parseConfig)

// Keywords returns a ParserOption that can be used to provide the Parser
// with a set of special keywords mapped to some custom TokenKind value.
// If the Parser encounters identifiers that match any of the given keywords,
// it returns a Token with the given kind and the actual literal encountered.
// Any default keywords are overwritten if specified in the custom set.
//
// Note: Use TokenKind values less than -10 for custom Token classes.
// -10 to -1 are reserved for standard token classes while 0 and above correspond the unicode code points.
func Keywords(keywords map[string]TokenKind) ParserOption {
	return func(config *parseConfig) {
		// Add each keyword to the config
		for keyword, kind := range keywords {
			config.keywords[keyword] = kind
		}
	}
}

// IgnoreWhitespaces returns a ParserOption that specifies the Parser to ignore unicode characters with the
// whitespace property (' ', '\t', '\n', '\r', etc). They are consumed instead of generating Tokens for them.
func IgnoreWhitespaces() ParserOption {
	return func(config *parseConfig) {
		config.eatSpaces = true
	}
}
