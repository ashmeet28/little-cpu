package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func GenerateInstructions(toks []TokenInfo) []string {
	type VarInfo struct {
		ident   string
		varType int
		arrLen  int
		funcSig []VarInfo
		scope   int
		addr    int
	}

	type BlockInfo struct {
		blockType int
	}

	const (
		VT_ILLEGAL int = iota
		VT_VOID
		VT_FUNC
		VT_INT
		VT_ARRAY
	)

	const (
		BT_ILLEGAL int = iota
		BT_FUNC
	)

	var varTable []VarInfo
	var blockTable []BlockInfo

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

	var GLOBAL_SCOPE int = 1
	var currScope int = GLOBAL_SCOPE

	var instructions []string

	// var REG_ZERO int = 0
	var REG_INST_BASE_ADDR int = 1
	var REG_GLOBAL_VAR_BASE_ADDR int = 2
	var REG_FRAME_PTR_ int = 3
	var REG_STACK_PTR int = 4
	var REG_A int = 8
	// var REG_B int = 9

	var REG_ZERO_X string = "00"
	var REG_INST_BASE_ADDR_X string = "01"
	var REG_GLOBAL_VAR_BASE_ADDR_X string = "02"
	var REG_FRAME_PTR_X string = "03"
	var REG_STACK_PTR_X string = "04"
	var REG_A_X string = "08"
	var REG_B_X string = "09"

	formatInst := func(op string, p1 string, p2 string, p3 string) string {
		op = (op + "    ")[:4]
		p1 = "00" + p1
		p1 = p1[len(p1)-2:]
		p2 = "00" + p2
		p2 = p2[len(p2)-2:]
		p3 = "00" + p3
		p3 = p3[len(p3)-2:]
		return (op + " " + p1 + " " + p2 + " " + p3)
	}

	emitInst := func(op string, p1 string, p2 string, p3 string) {
		instructions = append(instructions, formatInst(op, p1, p2, p3))
	}

	setInst := func(op string, p1 string, p2 string, p3 string, i int) {
		instructions[i] = formatInst(op, p1, p2, p3)
	}

	emitInstNOP := func() {
		emitInst("ADD", REG_ZERO_X, REG_ZERO_X, REG_ZERO_X)
	}

	setInstLoadImm := func(reg int, v int, i int) {
		setInst("LLIU", strconv.FormatInt(int64(reg), 16), strconv.FormatInt(int64(v&0xff), 16), strconv.FormatInt(int64((v>>8)&0xff), 16), i)
		setInst("LUI", strconv.FormatInt(int64(reg), 16), strconv.FormatInt(int64((v>>16)&0xff), 16), strconv.FormatInt(int64((v>>24)&0xff), 16), i+1)
	}

	emitInstLoadImm := func(reg int, v int) {
		emitInst("LLIU", strconv.FormatInt(int64(reg), 16), strconv.FormatInt(int64(v&0xff), 16), strconv.FormatInt(int64((v>>8)&0xff), 16))
		emitInst("LUI", strconv.FormatInt(int64(reg), 16), strconv.FormatInt(int64((v>>16)&0xff), 16), strconv.FormatInt(int64((v>>24)&0xff), 16))
	}

	emitInstStackPushWord := func() {
		emitInst("SW", REG_A_X, REG_ZERO_X, REG_STACK_PTR_X)
		emitInstLoadImm(REG_A, 4)
		emitInst("SUB", REG_STACK_PTR_X, REG_STACK_PTR_X, REG_A_X)
	}

	emitInstStackPopWord := func() {
		emitInstLoadImm(REG_A, 4)
		emitInst("ADD", REG_STACK_PTR_X, REG_STACK_PTR_X, REG_A_X)
		emitInst("LW", REG_A_X, REG_ZERO_X, REG_STACK_PTR_X)
	}

	emitInstStackStoreLocalWord := func() {
		emitInstStackPopWord()
		emitInst("ADD", REG_B_X, REG_ZERO_X, REG_A_X)
		emitInstStackPopWord()
		emitInst("SW", REG_B_X, REG_FRAME_PTR_X, REG_A_X)
	}

	emitInstStackStoreGlobalWord := func() {
		emitInstStackPopWord()
		emitInst("ADD", REG_B_X, REG_ZERO_X, REG_A_X)
		emitInstStackPopWord()
		emitInst("SW", REG_B_X, REG_GLOBAL_VAR_BASE_ADDR_X, REG_A_X)
	}

	emitInstInit := func() {
		emitInstLoadImm(REG_INST_BASE_ADDR, 0x1000_0000)

		emitInstLoadImm(REG_GLOBAL_VAR_BASE_ADDR, 0x2000_0000)

		emitInstLoadImm(REG_FRAME_PTR_, 0x3fff_ffff)

		emitInstLoadImm(REG_STACK_PTR, 0x3fff_ffff)

		// Jump to main function
		emitInstNOP()
		emitInstNOP()
		emitInstNOP()

		emitInst("ECALL", "00", "00", "00")
	}

	emitInstInit()

	findVar := func(ident string) VarInfo {
		for _, v := range varTable {
			if v.ident == ident {
				return v
			}
		}

		var varInfo VarInfo
		varInfo.varType = VT_ILLEGAL
		return varInfo
	}

	linkMainFunc := func() {
		var varInfo VarInfo = findVar("main")
		var JUMP_TO_MAIN_FUNC_INST_INDEX int = 8
		setInstLoadImm(REG_A, varInfo.addr, JUMP_TO_MAIN_FUNC_INST_INDEX)
		setInst("JALR", REG_A_X, REG_INST_BASE_ADDR_X, REG_A_X, JUMP_TO_MAIN_FUNC_INST_INDEX+2)
	}

	getNextInstAddr := func() int {
		return len(instructions) * 4
	}

	getNextLocalVarAddr := func() int {
		var addr int
		for _, varInfo := range varTable {
			if varInfo.scope != GLOBAL_SCOPE {
				if varInfo.varType == VT_INT {
					addr += 4
				} else if varInfo.varType == VT_ARRAY {
					addr += (varInfo.arrLen * 4)
				}
			}
		}
		return addr
	}

	getNextGlobalVarAddr := func() int {
		var addr int
		for _, varInfo := range varTable {
			if varInfo.scope == GLOBAL_SCOPE {
				if varInfo.varType == VT_INT {
					addr += 4
				} else if varInfo.varType == VT_ARRAY {
					addr += (varInfo.arrLen * 4)
				}
			}
		}
		return addr
	}

	clearLocalVarFromVarTable := func(scope int) {
		for len(varTable) != 0 && varTable[len(varTable)-1].scope > scope {
			varTable = varTable[:len(varTable)-1]
		}
	}

	for peek().tokType != TT_EOF {
		switch peek().tokType {
		case TT_FUNC:
			consume(TT_FUNC)

			var currVarInfo VarInfo
			currVarInfo.ident = consume(TT_IDENT).tokStr
			currVarInfo.varType = TT_FUNC
			currVarInfo.scope = currScope
			currVarInfo.addr = getNextInstAddr()
			varTable = append(varTable, currVarInfo)

			consume(TT_LPAREN)
			consume(TT_RPAREN)
			consume(TT_LBRACE)

			emitInstStackPushWord()
			emitInst("ADD", REG_FRAME_PTR_X, REG_ZERO_X, REG_STACK_PTR_X)

			var currBlockInfo BlockInfo
			currBlockInfo.blockType = BT_FUNC
			blockTable = append(blockTable, currBlockInfo)

			currScope++
		case TT_VAR:
			consume(TT_VAR)

			var currVarInfo VarInfo
			currVarInfo.ident = consume(TT_IDENT).tokStr
			if peek().tokStr == "int" {
				consume(TT_IDENT)
				currVarInfo.varType = VT_INT
			}
			currVarInfo.scope = currScope
			if currScope == GLOBAL_SCOPE {
				currVarInfo.addr = getNextGlobalVarAddr()
			} else {
				emitInstLoadImm(REG_A, 0)
				emitInstStackPushWord()
				currVarInfo.addr = getNextLocalVarAddr()
			}
			varTable = append(varTable, currVarInfo)
		case TT_RBRACE:
			var currBlockInfo BlockInfo = blockTable[len(blockTable)-1]
			blockTable = blockTable[:len(blockTable)-1]
			if currBlockInfo.blockType == BT_FUNC {
				emitInst("ADD", REG_STACK_PTR_X, REG_ZERO_X, REG_FRAME_PTR_X)
				emitInstStackPopWord()
				emitInst("ADD", REG_FRAME_PTR_X, REG_ZERO_X, REG_STACK_PTR_X)
				emitInst("JALR", REG_ZERO_X, REG_INST_BASE_ADDR_X, REG_A_X)
				currScope = GLOBAL_SCOPE
				clearLocalVarFromVarTable(currScope)
			}
			consume(TT_RBRACE)
		case TT_IDENT:
			var varInfo = findVar(consume(TT_IDENT).tokStr)
			emitInstLoadImm(REG_A, varInfo.addr)
			emitInstStackPushWord()

			consume(TT_ASSIGN)

			v, _ := strconv.ParseInt(consume(TT_INT).tokStr, 0, 64)
			emitInstLoadImm(REG_A, int(v))
			emitInstStackPushWord()

			if varInfo.scope == GLOBAL_SCOPE {
				emitInstStackStoreGlobalWord()
			} else {
				emitInstStackStoreLocalWord()
			}
		default:
			advance()
		}
	}

	linkMainFunc()

	for i, v := range instructions {
		fmt.Println(i, v)
	}
	for _, v := range varTable {
		fmt.Println(v)
	}

	return instructions
}

func GenerateBytecode(instructions []string) string {
	InstHex := map[string]string{
		"ECALL": "38",

		"ADD": "08",
		"SUB": "09",
		"XOR": "0c",
		"OR":  "0e",
		"AND": "0f",
		"SRA": "0a",
		"SRL": "0b",
		"SLL": "0d",

		"LB":  "10",
		"LH":  "11",
		"LW":  "12",
		"LBU": "14",
		"LHU": "15",

		"SB": "18",
		"SH": "19",
		"SW": "1a",

		"LUI":  "21",
		"LLI":  "22",
		"LLIU": "26",

		"BEQ":  "28",
		"BNE":  "29",
		"BLT":  "2c",
		"BGE":  "2d",
		"BLTU": "2e",
		"BGEU": "2f",

		"JAL":  "30",
		"JALR": "31",
	}

	var s string = strings.Join(instructions, " ")
	for opStr, hexStr := range InstHex {
		s = strings.ReplaceAll(s, opStr+" ", hexStr+" ")
	}
	instructions = strings.Split(s, " ")
	s = ""
	for _, inst := range instructions {
		if inst != "" {
			s = s + "0x" + inst + ", "
		}
	}
	s = s[:len(s)-2]

	return s
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	toks := GenerateTokens(data)

	insts := GenerateInstructions(toks)
	fmt.Println(GenerateBytecode(insts))
}
