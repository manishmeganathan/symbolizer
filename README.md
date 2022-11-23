# Symbolizer 🔣

[godoclink]: https://godoc.org/github.com/manishmeganathan/symbolizer
[![go docs](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godoclink]
![go version](https://img.shields.io/github/go-mod/go-version/manishmeganathan/symbolizer?style=flat-square)
![latest tag](https://img.shields.io/github/v/tag/manishmeganathan/symbolizer?color=brightgreen&label=latest%20tag&sort=semver&style=flat-square)
![license](https://img.shields.io/github/license/manishmeganathan/symbolizer?color=g&style=flat-square)
![test status](https://img.shields.io/github/workflow/status/manishmeganathan/symbolizer/Go%20Tests?label=tests&style=flat-square)
![issue count](https://img.shields.io/github/issues/manishmeganathan/symbolizer?style=flat-square&color=yellow)

A Go Package for Parsing Simple Symbols

### Overview
This package is designed for parsing very simple symbols and **not** large files or multi-file directories. It exposes a type ``Parser`` constructable
with ``NewParser`` using the string input that needs to be parsed and optional ``ParserOption`` functions to modify its behaviour.

### Installation
```
go get github.com/manishmeganathan/symbolizer
```

### Token Model
The ``Token`` type in this package contains the ``TokenKind``, the literal value as string as well the start position of the token.
``TokenKind`` are pseudo-runes that represents unicode code points for values above 0. Special tokens for literals and control are
represented as negative variants with values extending below 0. It can extended with custom variants but be mindful of collsions.

```go
// TokenKind is an enum for representing token grouping/values.
// For unicode tokens, the TokenKind is equal to its code point value.
// For literal such identifiers and numerics, the TokenKind values descend from 0.
// 
// Note: Custom TokenKind values can be used by external packages for keyword detection
// for special literals, but these values should be below -10 to prevent collisions
type TokenKind int32

const (
	TokenEoF TokenKind = -(iota + 1)
	TokenIdentifier
	TokenNumber
	TokenString
	TokenHexNumber
)

// Token represents a lexical Token.
// It may be either a lone unicode character or some literal value
type Token struct {
	Kind     TokenKind
	Literal  string
	Position int
}
```

### Usage Examples

```go
// TokenIterator describes a routine that leverages the Token inspection methods of Parser
// to view the current Token and check if the current Token is of a specific TokenKind.
func TokenIterator() {
    symbol := "map[string]string"
    parser := symbolizer.NewParser(symbol)
  
    // Check if the cursor has reach the end of the symbol
    for !parser.IsCursor(symbolizer.TokenEoF) { 
        // Print the current token
        fmt.Println(parser.Cursor()) 
        // Advance the cursor and ingest the next token
        parser.Advance()
    }

    // Output:
    // {<ident> map 0}
    // {<unicode:'['> [ 3}
    // {<ident> string 4}
    // {<unicode:']'> ] 10}
    // {<ident> string 11}
}
```

```go
// TokenLookAhead describes a routine that leverages the Token look ahead methods of Parser
// to view the next Token without ingesting it. ExpectPeek can be used to move the parser
// if the next token is of a specific kind.
func TokenLookAhead() {
    symbol := "[32]string"
    parser := symbolizer.NewParser(symbol)

    // Print the current and next token
    fmt.Println(parser.Cursor())
    fmt.Println(parser.Peek())
    
    // Check if the next token is a numeric
    // If it is, then the parse cursor is moved forward
    if parser.ExpectPeek(symbolizer.TokenNumber) {
        fmt.Println("numeric encountered")
    }
    
    // Print the current and next token after the peek expectation
    fmt.Println(parser.Cursor())
    fmt.Println(parser.Peek())

    // Output:
    // {<unicode:'['> [ 0}
    // {<num> 32 1}
    // numeric encountered
    // {<num> 32 1}
    // {<unicode:']'> ] 3}
}
```

```go
// CustomSymbols describes a routine that injects some custom keywords and token kinds
// into the parser, which can then be used to inspect tokens just as would regular TokenKind variants.
func CustomSymbols() {
    symbol := "map[string]string"
	
    // Define a custom token kind enum
    type MyTokenKind int32
    const Datatype = -10
    
    // Defines a mapping of identifier to custom token kinds
    keywords := map[string]symbolizer.TokenKind{"map": Datatype, "string": Datatype}
    // Create a Parser with a ParserOption that injects the custom keywords
    parser := symbolizer.NewParser(symbol, symbolizer.Keywords(keywords))
    
    // Check if the cursor has reach the end of the symbol
    for !parser.IsCursor(symbolizer.TokenEoF) {
        // Print the current token
        fmt.Println(parser.Cursor())
        // Advance the cursor and ingest the next token
        parser.Advance()
    }

    // Output:
    // 	{<custom:-10> map 0}
    // 	{<unicode:'['> [ 3}
    // 	{<custom:-10> string 4}
    // 	{<unicode:']'> ] 10}
    //  {<custom:-10> string 11}
}
```

```go
// SymbolSplit describes a routine that splits a Symbols into sub-strings (sub symbols) based on 
// some delimiter rune. Use a Whitespace Ignorant Parser if they should be ignored while splitting.
func SymbolSplit() {
    symbol := "23, 56, 8902342"
	
	// Create a Parser with a ParseOption to ignore whitespaces
    parser := NewParser(symbol, IgnoreWhitespaces())
    // Split the parser contents with the comma delimiter
    components := parser.Split(',')
	
	// Print all the components
    for _, component := range components {
        fmt.Println(component)
    }
	
    // Output: 
    // 23
    // 56
    // 8902342
}
```

```go
// SymbolUnwrap describes a routine that unwraps an inner sub-string (inner symbol) from within some 
// Enclosure, which is defined as a pair of unicode characters (cannot be the same). Unwrap can also
// handle nested contents of the same enclosure and works to resolve each opening with a closing.
func SymbolUnwrap() {
    symbol := "(outer[inner])"
    parser := NewParser(symbol)

	// Unwrap the symbol from within a set of parenthesis
    unwrapped, err := parser.Unwrap(EnclosureParens())
    if err != nil {
        panic(err)
    }

	// Print the unwrapped symbol
    fmt.Println(unwrapped)

    // Output: 
    // outer[inner]
}
```

### Notes:
This package is still a work in progress and can be heavily extended for a lot of different use cases.
If you are using this package and need some new functionality, please open an issue or a pull request.
