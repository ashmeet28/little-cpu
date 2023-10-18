package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	tokType int
	tokStr  string
}

func GenerateToken(src []byte) (TokenInfo, int) {
	var bytesConsumed int
	currTok := TokenInfo{tokType: TT_ILLEGAL, tokStr: ""}

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

	i := 0

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
	toks := []TokenInfo{}
	currTok := TokenInfo{tokType: TT_ILLEGAL, tokStr: ""}
	bytesConsumed := 0

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

func GenerateInstructions(toks []TokenInfo) []string {
	type VarInfo struct {
		ident   string
		varType int
		arrLen  int
		funcSig []VarInfo
		scope   int
		addr    int
	}

	const (
		VT_ILLEGAL int = iota
		VT_VOID
		VT_FUNC
		VT_INT
		VT_ARRAY
	)

	symbolTable := []VarInfo{}

	peek := func() TokenInfo {
		return toks[0]
	}

	advance := func() TokenInfo {
		tok := toks[0]
		toks = toks[1:]
		return tok
	}

	consume := func(tokType int) TokenInfo {
		tok := toks[0]
		if tok.tokType != tokType {
			fmt.Println("Error while consuming token")
			os.Exit(1)
		}
		toks = toks[1:]
		return tok
	}

	currScope := 1
	// freeVarMemAddr := 0

	var instructions []string

	emitInst := func(op string, p1 string, p2 string, p3 string) {
		op = (op + "    ")[:4]
		p1 = "00" + p1
		p1 = p1[len(p1)-2:]
		p2 = "00" + p2
		p2 = p2[len(p2)-2:]
		p3 = "00" + p3
		p3 = p3[len(p3)-2:]
		instructions = append(instructions, op+" "+p1+" "+p2+" "+p3)
	}

	getNextInstAddr := func() int {
		return len(instructions) * 4
	}

	emitInstNOP := func() {
		emitInst("ADD", "00", "00", "00")
	}

	emitInstLoadImm := func(reg int, v int) {
		emitInst("LLIU", strconv.FormatInt(int64(reg), 16), strconv.FormatInt(int64(v&0xff), 16), strconv.FormatInt(int64((v>>8)&0xff), 16))
		emitInst("LUI", strconv.FormatInt(int64(reg), 16), strconv.FormatInt(int64((v>>16)&0xff), 16), strconv.FormatInt(int64((v>>24)&0xff), 16))
	}

	emitInstInit := func() {
		// Global Variables Base Address
		emitInstLoadImm(1, 0x2000_0000)

		// Frame Pointer
		emitInstLoadImm(2, 0x3fff_ffff)

		// Stack Pointer
		emitInstLoadImm(3, 0x3fff_ffff)

		// Jump to main function
		emitInstLoadImm(8, 0)
		emitInstNOP()
	}

	// allocVarMem := func(s int) int {
	// 	addr := freeVarMemAddr
	// 	freeVarMemAddr = freeVarMemAddr + s
	// 	return addr
	// }

	emitInstInit()

	for peek().tokType != TT_EOF {
		switch peek().tokType {
		case TT_FUNC:
			consume(TT_FUNC)

			currVarInfo := VarInfo{}
			currVarInfo.varType = TT_FUNC
			currVarInfo.ident = consume(TT_IDENT).tokStr
			currVarInfo.scope = currScope
			currVarInfo.addr = getNextInstAddr()
			symbolTable = append(symbolTable, currVarInfo)

			consume(TT_LPAREN)
			consume(TT_RPAREN)
			consume(TT_LBRACE)
			consume(TT_NEW_LINE)
			consume(TT_RBRACE)
			emitInstNOP()
		default:
			advance()
		}
	}

	for _, v := range instructions {
		fmt.Println(v)
	}

	return instructions
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	toks := GenerateTokens(data)

	GenerateInstructions(toks)
}
