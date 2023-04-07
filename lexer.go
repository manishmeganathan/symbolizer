package symbolizer

import "unicode"

// lexer is a lexical analyser that can tokenize a given string input into its unicode
// characters while also generating tokens for identifiers, strings and numerics symbols.
type lexer struct {
	cursor  int
	symbols []rune
	config  *parseConfig
}

// char returns the unicode symbols that is currently under the Lexer's cursor.
// If the Lexer tape is exhausted, an EoF rune is returned.
func (lexer *lexer) char() rune {
	if lexer.done() {
		return rune(TokenEoF)
	}

	return lexer.symbols[lexer.cursor]
}

// peek returns the unicode symbol that is ahead of the Lexer's cursor.
// This look ahead is performed without moving the Lexer's cursor.
// If the Lexer tap is exhausted, an EoF rune is returned.
func (lexer *lexer) peek() rune {
	if lexer.done() {
		return rune(TokenEoF)
	}

	return lexer.symbols[lexer.cursor+1]
}

// tokens returns all the remaining Tokens in the lexer, by parsing
// through the rest of the input until it encounters an EoF.
// Note that if the lexer has already ingested some symbols, the
// returned Tokens do not represent all the Tokens for a given input.
func (lexer *lexer) tokens() (tokens []Token) {
	for {
		token := lexer.next()
		tokens = append(tokens, token)

		if token.Kind == TokenEoF {
			break
		}
	}

	return tokens
}

// done returns whether the Lexer tape is exhausted i.e., EoF has been reached
func (lexer *lexer) done() bool {
	return lexer.cursor >= len(lexer.symbols)
}

// next advances the Lexer's cursor and returns the encountered Token.
func (lexer *lexer) next() (token Token) {
	// If lexer configuration specifies to ignore whitespaces, consume them
	if lexer.config.eatSpaces {
		lexer.consumeSpaces()
	}

	// Get the current symbol of the Lexer and check conditions
	switch symbol := lexer.char(); {
	// End of File
	case symbol == rune(TokenEoF):
		token = EOFToken(lexer.cursor)

	// Quotes -> Scan for String
	case symbol == '"':
		token = lexer.scanString()

	// Hex Prefix
	case symbol == '0':
		if lexer.peek() == 'x' {
			return lexer.scanHexadecimal()
		}

		fallthrough

	// Digit -> Scan for Numeric (Integer/Float)
	case unicode.IsDigit(symbol):
		return lexer.scanNumeric()

	// Letter -> Scan for Identifier or Keyword
	case unicode.IsLetter(symbol):
		return lexer.scanIdentOrKeyword()

	// Negative Sign -> Scan for Numeric
	case symbol == '-':
		if isDecChar(lexer.peek()) {
			return lexer.scanNumeric()
		}

		fallthrough

	default:
		// Generate a token for the Unicode symbol
		token = UnicodeToken(symbol, lexer.cursor)
	}

	// Push the cursor for the next iteration
	lexer.advanceCursor()
	// Return the generated token
	return token
}

// advanceCursor increments the Lexer's cursor
func (lexer *lexer) advanceCursor() { lexer.cursor++ }

// collectBetween collects all symbols from a specified index
// until the current cursor position and return it as a string
func (lexer *lexer) collectBetween(start, stop int) string {
	return string(lexer.symbols[start:stop])
}

// consumeSpaces moves its cursor to the next character by skips all unicode whitespaces in between.
func (lexer *lexer) consumeSpaces() {
	// Iterate until the read character is a whitespace
	for unicode.IsSpace(lexer.char()) {
		lexer.advanceCursor()
	}
}

// lookupKeyword returns the TokenKind for a given identifier literal.
// If there exists rule entry for the identifier, then the TokenKind
// in the rule is returned, otherwise the literal is treated as a
// regular identifier and TokenIdentifier is returned.
func (lexer *lexer) lookupKeyword(ident string) TokenKind {
	// If no keywords available, immediately return TokenIdentifier
	if lexer.config.keywords == nil {
		return TokenIdent
	}

	// Retrieve the token kind for the ident from the keyword registry and return if it exists
	if tok, ok := lexer.config.keywords[ident]; ok {
		// Return the user defined identifier
		return tok
	}

	return TokenIdent
}

// scanIdentOrKeyword scans for an Identifier token, If the literal has a special
// TokenKind in the keyword registry, the returned Token has the appropriate TokenKind.
func (lexer *lexer) scanIdentOrKeyword() Token {
	// Retrieve the starting position of the identifier
	start := lexer.cursor

	// Iterate over the input until characters are letters
	for unicode.IsLetter(lexer.char()) || unicode.IsDigit(lexer.char()) || lexer.char() == '_' {
		lexer.advanceCursor()
	}

	// Extract the identifier from the input with the start and current position
	identifier := lexer.collectBetween(start, lexer.cursor)

	return Token{
		Kind:     lexer.lookupKeyword(identifier),
		Literal:  identifier,
		Position: start,
	}
}

// scanString scans for a String token by collecting characters until another '"' is encountered.
func (lexer *lexer) scanString() Token {
	// Retrieve the starting position
	start := lexer.cursor

	// Iterate over the input until an " or eof is encountered
	for {
		lexer.advanceCursor()

		if lexer.char() == '"' {
			break
		}

		// If EoF encountered prematurely, return malformed token
		if lexer.char() == rune(TokenEoF) {
			token := Token{
				Kind:     TokenMalformed,
				Literal:  lexer.collectBetween(start, lexer.cursor),
				Position: start,
			}

			lexer.cursor--
			return token
		}
	}

	// Extract the string from input and set as text token literal
	// Includes the quote characters as well
	return Token{
		Kind:     TokenString,
		Literal:  lexer.collectBetween(start, lexer.cursor+1),
		Position: start,
	}
}

// scanNumeric scans for a Numeric token (decimal or hexadecimal).
// If it encounters '0x', it will attempt to read the rest of the
// character as hexadecimal using scanHexadecimal
func (lexer *lexer) scanNumeric() Token {
	// Retrieve the starting position of the number
	start := lexer.cursor

	if lexer.char() == '-' {
		lexer.advanceCursor()
	}

	// Iterate over the input until characters are decimal characters
	for isDecChar(lexer.char()) {
		lexer.advanceCursor()
	}

	// Extract the number from input and set as number token literal
	return Token{
		Kind:     TokenNumber,
		Literal:  lexer.collectBetween(start, lexer.cursor),
		Position: start,
	}
}

// scanHexadecimal scans for a Hex Numeric Token. It must be invoked after
// encountering a '0x' and attempts to read hex characters A-F, a-f, 0-9.
func (lexer *lexer) scanHexadecimal() Token {
	// Retrieve the starting position of the identifier
	start := lexer.cursor

	lexer.advanceCursor()
	lexer.advanceCursor()

	// Iterate over the input until characters are hex characters
	for isHexChar(lexer.char()) {
		lexer.advanceCursor()
	}

	// Extract the number from input and set as digits token literal
	return Token{
		Kind:     TokenHexNumber,
		Literal:  lexer.collectBetween(start, lexer.cursor),
		Position: start,
	}
}

// isDecChar returns true if ch is a decimal character
func isDecChar(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isHexChar returns true if ch is a hexadecimal character
func isHexChar(ch rune) bool {
	return 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F' || isDecChar(ch)
}
