// Welcome to my own humble minimalistic TOML tokenizer
// it doesn't handle TOML cases the following and would give an error:
// - Quoted keys

package config

import (
	"bufio"
	"fmt"
	"internal/runtime/maps"
	"io"
	"os"
	"regexp"
)

type tokenType byte

type token struct {
	type_    tokenType
	value    string
	pos      struct {row uint // aka line
					 col uint}
}

type tokenStack []token

const (
	EOFToken tokenType = iota

	bareKeyToken // Simply an identifier but TOML calls them like that
	keyDotToken // e.g physical.color = "orange"
		        //             ^
	assignToken

	// Basic types
	floatToken
	integerToken
	boolToken

	// Strings
	multilineStringStartToken
	multilineStringEndToken
	stringStartToken
	stringEndToken
	
	stringEscToken // Such as \n or \t

	// Comments
	commentToken
	multilineCommentStartToken
	multilineCommentEndToken

	// Composite tokens
	arrayStartToken
	arrayEndToken
	tableStartToken
	tableEndToken
	arrayTableStartToken
	arrayTableEndToken
	inlineTableStartToken
	inlineTableEndToken
)

var tokensExpressionsTable = map[tokenType]regexp.Regexp{
	// NOTE: there can't be two tokens with the same position in the stack
	// therefore there must not be a token with multiple types in the same time
	// in other words there must not be a token hirachy.

	// FIXME: add tests for those regexes *(cries in regex)*

	bareKeyToken:               *regexp.MustCompile(`^[ \t]*[A-Za-z0-9_-]+([ \t]|\.[A-Za-z0-9_-]+)*=`),
	keyDotToken:                *regexp.MustCompile(``),

	assignToken:                *regexp.MustCompile(`^=`),

	floatToken:                 *regexp.MustCompile(``),
	integerToken:               *regexp.MustCompile(``),
	boolToken:                  *regexp.MustCompile(``),

	multilineStringStartToken:  *regexp.MustCompile(``),
	multilineStringEndToken:    *regexp.MustCompile(``),
	stringStartToken:           *regexp.MustCompile(``),
	stringEndToken:             *regexp.MustCompile(``),
	stringEscToken:             *regexp.MustCompile(``), 

	commentToken:               *regexp.MustCompile(`#.*$`),
	multilineCommentStartToken: *regexp.MustCompile(``),
	multilineCommentEndToken:   *regexp.MustCompile(``),

	arrayStartToken:            *regexp.MustCompile(`[^\[]\[[^\[]`),
	arrayEndToken:              *regexp.MustCompile(`[^\]]\][^\]]`),
	tableStartToken:            *regexp.MustCompile(`\[\[`),
	tableEndToken:              *regexp.MustCompile(``),
	arrayTableStartToken:       *regexp.MustCompile(``),
	arrayTableEndToken:         *regexp.MustCompile(``),
	inlineTableStartToken:      *regexp.MustCompile(``),
	inlineTableEndToken:        *regexp.MustCompile(``),
}


