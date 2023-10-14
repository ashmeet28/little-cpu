package main

import (
	"fmt"
	"log"
	"os"
)

const (
	// Token Types

	TT_ILLEGAL int = iota
	TT_EOF
	TT_NEW_LINE
	TT_SPACE

	TT_IDENT  // main
	TT_INT    // 12345
	TT_CHAR   // 'a'
	TT_STRING // "abc"

	TT_ADD    // +
	TT_SUB    // -
	TT_MUL    // *
	TT_QUO    // /
	TT_REM    // %
	TT_AND    // &
	TT_OR     // |
	TT_XOR    // ^
	TT_SHL    // <<
	TT_SHR    // >>
	TT_LAND   // &&
	TT_LOR    // ||
	TT_EQL    // ==
	TT_LSS    // <
	TT_GTR    // >
	TT_ASSIGN // =
	TT_NOT    // !
	TT_NEQ    // !=
	TT_LEQ    // <=
	TT_GEQ    // >=

	TT_LPAREN // (
	TT_LBRACK // [
	TT_LBRACE // {
	TT_RPAREN // )
	TT_RBRACK // ]
	TT_RBRACE // }
	TT_COMMA  // ,
	TT_PERIOD // .
	SEMICOLON // ;
	COLON     // :

	TT_WHILE
	TT_BREAK
	TT_CONTINUE
	TT_IF
	TT_ELSE
	TT_FUNC
	TT_RETURN
	TT_VAR
)

type TokenInfo struct {
	TokenType   int
	TokenString string
}

func GetToken(source []byte) (TokenInfo, int) {
	bytesConsumed := 0
	tok := TokenInfo{TokenType: TT_ILLEGAL, TokenString: ""}

	if len(source) == 0 {

		tok.TokenType = TT_EOF
		return tok, bytesConsumed

	} else if source[0] == 0x0a {

		tok.TokenType = TT_NEW_LINE
		bytesConsumed = 1
		return tok, bytesConsumed

	} else if source[0] == 0x20 {

		tok.TokenType = TT_SPACE
		bytesConsumed = 1
		return tok, bytesConsumed

	}

	srcStr := ""

	for i, c := range source {
		if c == 0x0a {
			srcStr = string(source[:i])
			break
		}
	}

	if srcStr == "" {
		srcStr = string(source)
	}

	tok.TokenType = TT_ILLEGAL

	TokensStrings := map[int]string{
		TT_ADD:    "+",
		TT_SUB:    "-",
		TT_MUL:    "*",
		TT_QUO:    "/",
		TT_REM:    "%",
		TT_AND:    "&",
		TT_OR:     "|",
		TT_XOR:    "^",
		TT_SHL:    "<<",
		TT_SHR:    ">>",
		TT_LAND:   "&&",
		TT_LOR:    "||",
		TT_EQL:    "==",
		TT_LSS:    "<",
		TT_GTR:    ">",
		TT_ASSIGN: "=",
		TT_NOT:    "!",
		TT_NEQ:    "!=",
		TT_LEQ:    "<=",
		TT_GEQ:    ">=",

		TT_LPAREN: "(",
		TT_LBRACK: "[",
		TT_LBRACE: "{",
		TT_RPAREN: ")",
		TT_RBRACK: "]",
		TT_RBRACE: "}",
		TT_COMMA:  ",",
		TT_PERIOD: ".",
		SEMICOLON: ";",
		COLON:     ":",

		TT_WHILE:    "while",
		TT_BREAK:    "break",
		TT_CONTINUE: "continue",
		TT_IF:       "if",
		TT_ELSE:     "else",
		TT_FUNC:     "func",
		TT_RETURN:   "return",
		TT_VAR:      "var",
	}

	for tokType, tokStr := range TokensStrings {
		if len(srcStr) >= len(tokStr) && srcStr[:len(tokStr)] == tokStr {
			if tok.TokenType == TT_ILLEGAL || len(tok.TokenString) < len(tokStr) {
				tok.TokenType = tokType
				tok.TokenString = tokStr
				bytesConsumed = len(tokStr)
			}
		}
	}

	if tok.TokenType != TT_ILLEGAL {
		return tok, bytesConsumed
	}

	isDigit := func(c byte) bool {
		return c >= 0x30 && c <= 0x39
	}

	isAplabet := func(c byte) bool {
		return (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) || (c == 0x5f)
	}

	i := 0

	if isAplabet(srcStr[i]) {

		tok.TokenType = TT_IDENT
		for (i < len(srcStr)) && (isAplabet(srcStr[i]) || isDigit(srcStr[i])) {
			i++
		}
		tok.TokenString = srcStr[:i]
		bytesConsumed = len(srcStr[:i])

	} else if isDigit(srcStr[i]) {

		tok.TokenType = TT_INT
		for (i < len(srcStr)) && (isAplabet(srcStr[i]) || isDigit(srcStr[i])) {
			i++
		}
		tok.TokenString = srcStr[:i]
		bytesConsumed = len(srcStr[:i])

	}

	return tok, bytesConsumed
}

func GetTokens(s []byte) []TokenInfo {
	allToks := []TokenInfo{}
	tok := TokenInfo{TokenType: TT_ILLEGAL, TokenString: ""}
	bytesConsumed := 0

	for tok.TokenType != TT_EOF {
		tok, bytesConsumed = GetToken(s)
		s = s[bytesConsumed:]
		if tok.TokenType != TT_SPACE {
			allToks = append(allToks, tok)
		}
	}

	return allToks
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for _, tok := range GetTokens(data) {
		fmt.Println(tok)
	}
}
