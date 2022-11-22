package symbolizer

import "unicode"

// LexerOption represents an option to modify the Lexer behaviour.
// It must be provided with the constructor for Lexer.
type LexerOption func(config *lexerConfig) error

// Keywords returns a LexerOption that can be used to provide the Lexer
// with a set of special keywords mapped to some custom TokenKind value.
// If the Lexer encounters identifiers that match any of the given keywords,
// It returns a Token with the given kind and the actual literal encountered.
//
// Note: Use TokenKind values less than -10 for custom Token classes.
// -10 to -1 are reserved for standard token classes while 0 and above correspond the unicode code points.
func Keywords(keywords map[string]TokenKind) LexerOption {
	return func(config *lexerConfig) error {
		config.keywords = keywords

		return nil
	}
}

// IgnoreWhitespaces returns a LexerOption that specifies the Lexer to ignore unicode characters with the
// whitespace property (' ', '\t', '\n', '\r', etc). They are consumed instead of generating Tokens for them.
func IgnoreWhitespaces() LexerOption {
	return func(config *lexerConfig) error {
		config.eatSpaces = true

		return nil
	}
}

// Lexer is a lexical analyser that can tokenize a given string input into its unicode
// characters while also generating tokens for identifiers, strings and numerics symbols.
type Lexer struct {
	cursor  int
	symbols []rune
	config  *lexerConfig
}

// lexerConfig is an internal configuration object for the
// Lexer that are modified using LexerOption functions
type lexerConfig struct {
	eatSpaces bool
	keywords  map[string]TokenKind
}

// NewLexer generates a new Lexer for a given string input and some LexerOption functions.
// Use the Keywords option to specify custom symbols to treat as keywords/identifiers.
func NewLexer(input string, opts ...LexerOption) (lexer *Lexer, err error) {
	// Create a new lexerConfig and apply all the given options on it
	config := new(lexerConfig)
	for _, option := range opts {
		_ = option(config)
	}

	// Convert the input string into a slice of runes (unicode codepoints)
	return &Lexer{
		config:  config,
		symbols: []rune(input),
	}, nil
}

// Symbol returns the unicode symbols that is currently under the Lexer's cursor.
// If the Lexer tape is exhausted, an EoF rune is returned.
func (lexer *Lexer) Symbol() rune {
	if lexer.Done() {
		return rune(TokenEoF)
	}

	return lexer.symbols[lexer.cursor]
}

// PeekSymbol returns the unicode symbol that is ahead of the Lexer's cursor.
// This look ahead is performed without moving the Lexer's cursor.
// If the Lexer tap is exhausted, an EoF rune is returned.
func (lexer *Lexer) PeekSymbol() rune {
	if lexer.Done() {
		return rune(TokenEoF)
	}

	return lexer.symbols[lexer.cursor+1]
}

// Tokens returns all the remaining Tokens in the Lexer, by parsing
// through the rest of the input until it encounters an EoF.
// Note that if the Lexer has already ingested some symbols, the
// returned Tokens do not represent all the Tokens for a given input.
func (lexer *Lexer) Tokens() (tokens []Token) {
	for {
		token := lexer.Next()
		tokens = append(tokens, token)

		if token.Kind == TokenEoF {
			break
		}
	}

	return tokens
}

// Done returns whether the Lexer tape is exhausted i.e., EoF has been reached
func (lexer *Lexer) Done() bool {
	return lexer.cursor >= len(lexer.symbols)
}

// Next advances the Lexer's cursor and returns the encountered Token.
// Whitespaces are consumed by default
func (lexer *Lexer) Next() Token {
	// If lexer configuration specifies to ignore whitespaces, consume them
	if lexer.config.eatSpaces {
		lexer.consumeSpaces()
	}

	var token Token

	// Get the current symbol of the Lexer and check conditions
	switch symbol := lexer.Symbol(); {
	// End of File
	case symbol == rune(TokenEoF):
		token = EOFToken()

	// Quotes -> Scan for String
	case symbol == '"':
		token = lexer.scanString()

	// Letter -> Scan for Identifier or Keyword
	case unicode.IsLetter(symbol):
		return lexer.scanIdentOrKeyword()

	// Digit -> Scan for Numeric (Integer/Float)
	case unicode.IsDigit(symbol):
		return lexer.scanNumeric()

	default:
		// Generate a token for the Unicode symbol
		token = UnicodeToken(symbol)
	}

	// Push the cursor for the next iteration
	lexer.advanceCursor()
	// Return the generated token
	return token
}

// advanceCursor increments the Lexer's cursor
func (lexer *Lexer) advanceCursor() { lexer.cursor++ }

// collectFrom collects all symbols from a specified index
// until the current cursor position and return it as a string
func (lexer *Lexer) collectFrom(start int) string {
	return string(lexer.symbols[start:lexer.cursor])
}

// consumeSpaces moves its cursor to the next character by skips all unicode whitespaces in between.
func (lexer *Lexer) consumeSpaces() {
	// Iterate until the read character is a whitespace
	for unicode.IsSpace(lexer.Symbol()) {
		lexer.advanceCursor()
	}
}

// lookupKeyword returns the TokenKind for a given identifier literal.
// If there exists rule entry for the identifier, then the TokenKind
// in the rule is returned, otherwise the literal is treated as a
// regular identifier and TokenIdentifier is returned.
func (lexer *Lexer) lookupKeyword(ident string) TokenKind {
	// If no keywords available, immediately return TokenIdentifier
	if lexer.config.keywords == nil {
		return TokenIdentifier
	}

	// Retrieve the token kind for the ident from the keyword registry and return if it exists
	if tok, ok := lexer.config.keywords[ident]; ok {
		// Return the user defined identifier
		return tok
	}

	return TokenIdentifier
}

// scanIdentOrKeyword scans for an Identifier token, If the literal has a special
// TokenKind in the keyword registry, the returned Token has the appropriate TokenKind.
func (lexer *Lexer) scanIdentOrKeyword() Token {
	// Retrieve the starting position of the identifier
	start := lexer.cursor

	// Iterate over the input until characters are letters
	for unicode.IsLetter(lexer.Symbol()) || unicode.IsDigit(lexer.Symbol()) || lexer.Symbol() == '_' {
		lexer.advanceCursor()
	}

	// Extract the identifier from the input with the start and current position
	identifier := lexer.collectFrom(start)

	return Token{Kind: lexer.lookupKeyword(identifier), Value: identifier}
}

// scanString scans for a String token by collecting characters until another '"' is encountered.
func (lexer *Lexer) scanString() Token {
	// Retrieve the starting position of the number (after the ")
	start := lexer.cursor + 1

	// Iterate over the input until an " or eof is encountered
	for {
		lexer.advanceCursor()

		if lexer.Symbol() == '"' || lexer.Symbol() == rune(TokenEoF) {
			break
		}
	}

	// Extract the string from input and set as text token literal
	return Token{Kind: TokenString, Value: lexer.collectFrom(start)}
}

// scanNumeric scans for a Numeric token (decimal or hexadecimal).
// If it encounters '0x', it will attempt to read the rest of the
// character as hexadecimal using scanHexadecimal
func (lexer *Lexer) scanNumeric() Token {
	if lexer.Symbol() == '0' && lexer.PeekSymbol() == 'x' {
		lexer.advanceCursor()
		lexer.advanceCursor()

		return lexer.scanHexadecimal()
	}

	// Retrieve the starting position of the number
	start := lexer.cursor

	// Iterate over the input until characters are decimal characters
	for isDecChar(lexer.Symbol()) {
		lexer.advanceCursor()
	}

	// Extract the number from input and set as number token literal
	return Token{Kind: TokenNumber, Value: lexer.collectFrom(start)}
}

// scanHexadecimal scans for a Hex Numeric Token. It must be invoked after
// encountering a '0x' and attempts to read hex characters A-F, a-f, 0-9.
func (lexer *Lexer) scanHexadecimal() Token {
	// Retrieve the starting position of the identifier
	start := lexer.cursor

	// Iterate over the input until characters are hex characters
	for isHexChar(lexer.Symbol()) {
		lexer.advanceCursor()
	}

	// Extract the number from input and set as digits token literal
	return Token{Kind: TokenHexNumber, Value: lexer.collectFrom(start)}
}

// isDecChar returns true if ch is a decimal character
func isDecChar(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isHexChar returns true if ch is a hexadecimal character
func isHexChar(ch rune) bool {
	return 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F' || isDecChar(ch)
}
