package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type TokenType int

const (
	// Token Types

	TT_ILLEGAL TokenType = iota
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

	TT_LPAREN    // (
	TT_LBRACK    // [
	TT_LBRACE    // {
	TT_RPAREN    // )
	TT_RBRACK    // ]
	TT_RBRACE    // }
	TT_COMMA     // ,
	TT_PERIOD    // .
	TT_SEMICOLON // ;
	TT_COLON     // :

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
	tokType TokenType
	tokStr  string
}

var (
	OP_NOP   byte = 1
	OP_ECALL byte = 2

	OP_ADD byte = 4
	OP_SUB byte = 5
	OP_XOR byte = 6
	OP_OR  byte = 7
	OP_AND byte = 8
	OP_SR  byte = 9
	OP_SL  byte = 10

	OP_PUSH_LITERAL byte = 16
	OP_PUSH_LOCAL   byte = 17
	OP_PUSH_GLOBAL  byte = 18

	OP_POP_LITERAL byte = 20
	OP_POP_LOCAL   byte = 21
	OP_POP_GLOBAL  byte = 22

	OP_EQ byte = 24
	OP_NE byte = 25
	OP_LT byte = 26
	OP_GE byte = 27

	OP_JUMP   byte = 28
	OP_CALL   byte = 29
	OP_RETURN byte = 30
)

func GenerateToken(src []byte) (TokenInfo, int) {
	var bytesConsumed int

	var currTok TokenInfo
	currTok.tokType = TT_ILLEGAL

	if len(src) == 0 {
		currTok.tokType = TT_EOF
		bytesConsumed = 0
		return currTok, bytesConsumed
	} else if src[0] == 0x0a {
		currTok.tokType = TT_NEW_LINE
		bytesConsumed = 1
		return currTok, bytesConsumed
	} else if src[0] == 0x20 {
		currTok.tokType = TT_SPACE
		bytesConsumed = 1
		return currTok, bytesConsumed
	}

	var srcStr string

	for i, c := range src {
		if c == 0x0a {
			srcStr = string(src[:i])
			break
		}
	}

	if len(srcStr) == 0 {
		srcStr = string(src)
	}

	TokensStrings := map[TokenType]string{
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

		TT_LPAREN:    "(",
		TT_LBRACK:    "[",
		TT_LBRACE:    "{",
		TT_RPAREN:    ")",
		TT_RBRACK:    "]",
		TT_RBRACE:    "}",
		TT_COMMA:     ",",
		TT_PERIOD:    ".",
		TT_SEMICOLON: ";",
		TT_COLON:     ":",

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
			if currTok.tokType == TT_ILLEGAL || len(currTok.tokStr) < len(tokStr) {
				currTok.tokType = tokType
				currTok.tokStr = tokStr
				bytesConsumed = len(tokStr)
			}
		}
	}

	if currTok.tokType != TT_ILLEGAL {
		return currTok, bytesConsumed
	}

	isDigit := func(c byte) bool {
		return c >= 0x30 && c <= 0x39
	}

	isAplabet := func(c byte) bool {
		return (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) || (c == 0x5f)
	}

	var i int = 0

	if isAplabet(srcStr[i]) {

		currTok.tokType = TT_IDENT
		for (i < len(srcStr)) && (isAplabet(srcStr[i]) || isDigit(srcStr[i])) {
			i++
		}
		currTok.tokStr = srcStr[:i]
		bytesConsumed = len(srcStr[:i])

	} else if isDigit(srcStr[i]) {

		currTok.tokType = TT_INT
		for (i < len(srcStr)) && (isAplabet(srcStr[i]) || isDigit(srcStr[i])) {
			i++
		}
		currTok.tokStr = srcStr[:i]
		bytesConsumed = len(srcStr[:i])

	}

	return currTok, bytesConsumed
}

func GenerateTokens(src []byte) []TokenInfo {
	var toks []TokenInfo
	var currTok TokenInfo
	currTok.tokType = TT_ILLEGAL
	var bytesConsumed int = 0

	for currTok.tokType != TT_EOF {
		currTok, bytesConsumed = GenerateToken(src)
		if currTok.tokType == TT_ILLEGAL {
			fmt.Println("Error while tokenizing")
			os.Exit(1)
		}
		src = src[bytesConsumed:]
		if currTok.tokType != TT_SPACE {
			toks = append(toks, currTok)
		}
	}

	return toks
}

func GenerateBytecode(toks []TokenInfo) []byte {
	type VarType int

	const (
		VT_ILLEGAL VarType = iota
		VT_FUNC
		VT_INT
	)

	type VarInfo struct {
		ident   string
		varType VarType
		scope   int
		addr    int
	}

	type BlockType int

	type BlockInfo struct {
		blockType BlockType
	}

	const (
		BT_ILLEGAL BlockType = iota
		BT_FUNC
	)

	var varTable []VarInfo
	var blockTable []BlockInfo

	var GLOBAL_SCOPE int = 1
	var currScope int = GLOBAL_SCOPE

	var blankLiteralAddr []int

	var allBytes []byte

	emitByte := func(b byte) {
		allBytes = append(allBytes, b)
	}

	emitPushLiteral := func(lit int) {
		for ((len(allBytes) + 1) % 4) != 0 {
			emitByte(OP_NOP)
		}

		emitByte(OP_PUSH_LITERAL)

		allBytes = append(allBytes, byte(lit&0xff))
		allBytes = append(allBytes, byte((lit>>8)&0xff))
		allBytes = append(allBytes, byte((lit>>16)&0xff))
		allBytes = append(allBytes, byte((lit>>24)&0xff))
	}

	emitPushBlankLiteral := func() {
		var lit int = 0

		for ((len(allBytes) + 1) % 4) != 0 {
			emitByte(OP_NOP)
		}

		emitByte(OP_PUSH_LITERAL)

		blankLiteralAddr = append(blankLiteralAddr, len(allBytes))

		allBytes = append(allBytes, byte(lit&0xff))
		allBytes = append(allBytes, byte((lit>>8)&0xff))
		allBytes = append(allBytes, byte((lit>>16)&0xff))
		allBytes = append(allBytes, byte((lit>>24)&0xff))
	}

	setPushBlankLiteral := func(lit int) {
		var i int = blankLiteralAddr[len(blankLiteralAddr)-1]

		blankLiteralAddr = blankLiteralAddr[:len(blankLiteralAddr)-1]

		allBytes[i] = byte(lit & 0xff)
		allBytes[i+1] = byte((lit >> 8) & 0xff)
		allBytes[i+2] = byte((lit >> 16) & 0xff)
		allBytes[i+3] = byte((lit >> 24) & 0xff)
	}

	findVar := func(ident string) VarInfo {
		var varInfo VarInfo
		varInfo.varType = VT_ILLEGAL

		for _, v := range varTable {
			if v.ident == ident {
				varInfo = v
				break
			}
		}

		return varInfo
	}

	peek := func() TokenInfo {
		return toks[0]
	}

	advance := func() TokenInfo {
		tok := toks[0]
		toks = toks[1:]
		return tok
	}

	consume := func(tokType TokenType) TokenInfo {
		tok := toks[0]
		if tok.tokType != tokType {
			fmt.Println("Error while consuming token")
			os.Exit(1)
		}
		toks = toks[1:]
		return tok
	}

	getNextByteAddr := func() int {
		return len(allBytes)
	}

	getNextLocalVarAddr := func() int {
		var addr int = 0
		for _, varInfo := range varTable {
			if varInfo.scope != GLOBAL_SCOPE {
				addr++
			}
		}
		return addr
	}

	getNextGlobalVarAddr := func() int {
		var addr int = 0
		for _, varInfo := range varTable {
			if varInfo.scope == GLOBAL_SCOPE {
				addr++
			}
		}
		return addr
	}

	clearLocalVarFromVarTable := func(scope int) {
		for len(varTable) != 0 && varTable[len(varTable)-1].scope > scope {
			varTable = varTable[:len(varTable)-1]
		}
	}

	emitPushBlankLiteral()
	emitByte(OP_CALL)
	emitByte(OP_ECALL)

	for peek().tokType != TT_EOF {
		switch peek().tokType {
		case TT_FUNC:
			consume(TT_FUNC)

			var currVarInfo VarInfo
			currVarInfo.ident = consume(TT_IDENT).tokStr
			currVarInfo.varType = VT_FUNC
			currVarInfo.scope = currScope
			currVarInfo.addr = getNextByteAddr()

			varTable = append(varTable, currVarInfo)

			var blockInfo BlockInfo
			blockInfo.blockType = BT_FUNC
			blockTable = append(blockTable, blockInfo)

			consume(TT_LPAREN)
			consume(TT_RPAREN)
			consume(TT_LBRACE)
			consume(TT_NEW_LINE)

			currScope++

		case TT_VAR:
			consume(TT_VAR)

			var currVarInfo VarInfo
			currVarInfo.ident = consume(TT_IDENT).tokStr
			currVarInfo.varType = VT_INT
			currVarInfo.scope = currScope

			if currVarInfo.scope == GLOBAL_SCOPE {
				currVarInfo.addr = getNextGlobalVarAddr()
				emitPushLiteral(currVarInfo.addr)
				emitPushLiteral(0)
				emitByte(OP_POP_GLOBAL)
			} else {
				currVarInfo.addr = getNextLocalVarAddr()
				emitPushLiteral(0)
			}

			varTable = append(varTable, currVarInfo)

			consume(TT_NEW_LINE)

		case TT_IDENT:
			var varInfo = findVar(consume(TT_IDENT).tokStr)

			emitPushLiteral(varInfo.addr)

			consume(TT_ASSIGN)

			v, _ := strconv.ParseInt(consume(TT_INT).tokStr, 0, 64)
			emitPushLiteral(int(v))

			if varInfo.scope == GLOBAL_SCOPE {
				emitByte(OP_POP_GLOBAL)
			} else {
				emitByte(OP_POP_LOCAL)
			}

			consume(TT_NEW_LINE)

		case TT_RBRACE:

			var blockInfo BlockInfo
			blockInfo = blockTable[len(blockTable)-1]
			blockTable = blockTable[:len(blockTable)-1]

			if blockInfo.blockType == BT_FUNC {
				currScope = GLOBAL_SCOPE
				clearLocalVarFromVarTable(currScope)
				emitByte(OP_RETURN)
			}

			consume(TT_RBRACE)
			consume(TT_NEW_LINE)

		case TT_NEW_LINE:
			advance()

		default:
			fmt.Println("Error while compiling")
			os.Exit(1)
		}
	}

	setPushBlankLiteral(int(findVar("main").addr))

	return allBytes
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	toks := GenerateTokens(data)

	byteCode := GenerateBytecode(toks)

	fmt.Println(byteCode)

	os.WriteFile(os.Args[2], byteCode, 0666)
}
