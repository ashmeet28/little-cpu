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

	const (
		VT_ILLEGAL int = iota
		VT_VOID
		VT_FUNC
		VT_INT
		VT_ARRAY
	)

	type BlockInfo struct {
		blockType int
	}

	const (
		BT_ILLEGAL int = iota
		BT_FUNC
	)

	var varTable []VarInfo
	var blockTable []BlockInfo

	var GLOBAL_SCOPE int = 1
	var currScope int = GLOBAL_SCOPE

	var instructions []string

	var REG_ZERO string = "00"
	var REG_INST string = "01"
	var REG_GLOBAL string = "02"
	var REG_FRAME string = "03"
	var REG_STACK string = "04"
	var REG_A string = "08"
	var REG_B string = "09"
	var REG_C string = "10"

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

	getNextInstAddr := func() int {
		return len(instructions) << 2
	}

	getNextLocalVarAddr := func() int {
		var addr int
		for _, varInfo := range varTable {
			if varInfo.scope != GLOBAL_SCOPE {
				if varInfo.varType == VT_INT {
					addr += 4
				} else if varInfo.varType == VT_ARRAY {
					addr += (varInfo.arrLen << 2)
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
					addr += (varInfo.arrLen << 2)
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

	formatInst := func(op string, p1 string, p2 string, p3 string) string {
		op = (op + "      ")[:5]
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
		emitInst("ADD", REG_ZERO, REG_ZERO, REG_ZERO)
	}

	emitInstLoadImm := func(v string) {
		v = "00000000" + v
		v = v[len(v)-8:]
		emitInst("LLIU", REG_A, v[6:8], v[4:6])
		emitInst("LUI", REG_B, v[2:4], v[0:2])
		emitInst("ADD", REG_A, REG_A, REG_B)
	}

	setInstLoadImm := func(v string, i int) {
		v = "00000000" + v
		v = v[len(v)-8:]
		setInst("LLIU", REG_A, v[6:8], v[4:6], i)
		setInst("LUI", REG_B, v[2:4], v[0:2], i+1)
		setInst("ADD", REG_A, REG_A, REG_B, i+2)
	}

	emitInstStackPushWord := func() {
		emitInst("SW", REG_A, REG_ZERO, REG_STACK)
		emitInstLoadImm("4")
		emitInst("ADD", REG_STACK, REG_STACK, REG_A)
	}

	emitInstStackPopWord := func() {
		emitInstLoadImm("4")
		emitInst("SUB", REG_STACK, REG_STACK, REG_A)
		emitInst("LW", REG_A, REG_ZERO, REG_STACK)
	}

	emitInstStackStoreLocalWord := func() {
		emitInstStackPopWord()
		emitInst("ADD", REG_C, REG_ZERO, REG_A)
		emitInstStackPopWord()
		emitInst("SW", REG_C, REG_FRAME, REG_A)
	}

	emitInstStackStoreGlobalWord := func() {
		emitInstStackPopWord()
		emitInst("ADD", REG_C, REG_ZERO, REG_A)
		emitInstStackPopWord()
		emitInst("SW", REG_C, REG_GLOBAL, REG_A)
	}

	emitInstInit := func() {
		emitInstLoadImm(strconv.FormatInt(0x1000_0000, 16))
		emitInst("ADD", REG_INST, REG_ZERO, REG_A)

		emitInstLoadImm(strconv.FormatInt(0x2000_0000, 16))
		emitInst("ADD", REG_GLOBAL, REG_ZERO, REG_A)

		emitInstLoadImm(strconv.FormatInt(0x3000_0000, 16))
		emitInst("ADD", REG_FRAME, REG_ZERO, REG_A)

		emitInstLoadImm(strconv.FormatInt(0x3000_0000, 16))
		emitInst("ADD", REG_STACK, REG_ZERO, REG_A)

		// Jump to main function
		emitInstNOP()
		emitInstNOP()
		emitInstNOP()
		emitInstNOP()

		emitInst("ECALL", "00", "00", "00")
	}

	emitInstLinkMainFunc := func() {
		var varInfo VarInfo = findVar("main")
		var MAIN_FUNC_JUMP_INST_INDEX int = 16
		setInstLoadImm(strconv.FormatInt(int64(varInfo.addr), 16), MAIN_FUNC_JUMP_INST_INDEX)
		setInst("JALR", REG_A, REG_INST, REG_A, MAIN_FUNC_JUMP_INST_INDEX+3)
	}

	emitInstInit()

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
			emitInst("ADD", REG_FRAME, REG_ZERO, REG_STACK)

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
				emitInstLoadImm("0")
				emitInstStackPushWord()
				currVarInfo.addr = getNextLocalVarAddr()
			}
			varTable = append(varTable, currVarInfo)

		case TT_IDENT:

			var varInfo = findVar(consume(TT_IDENT).tokStr)
			emitInstLoadImm(strconv.FormatInt(int64(varInfo.addr), 16))
			emitInstStackPushWord()

			consume(TT_ASSIGN)

			v, _ := strconv.ParseInt(consume(TT_INT).tokStr, 0, 64)

			emitInstLoadImm(strconv.FormatInt(int64(v), 16))

			emitInstStackPushWord()

			if varInfo.scope == GLOBAL_SCOPE {
				emitInstStackStoreGlobalWord()
			} else {
				emitInstStackStoreLocalWord()
			}

		case TT_RBRACE:

			var currBlockInfo BlockInfo = blockTable[len(blockTable)-1]
			blockTable = blockTable[:len(blockTable)-1]

			if currBlockInfo.blockType == BT_FUNC {
				emitInst("ADD", REG_STACK, REG_ZERO, REG_FRAME)
				emitInstStackPopWord()
				emitInst("ADD", REG_FRAME, REG_ZERO, REG_STACK)
				emitInst("JALR", REG_ZERO, REG_ZERO, REG_A)
				currScope = GLOBAL_SCOPE
				clearLocalVarFromVarTable(currScope)
			}

			consume(TT_RBRACE)

		default:
			advance()
		}
	}

	emitInstLinkMainFunc()

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

	var instStr string = strings.Join(instructions, " ")
	for opStr, hexStr := range InstHex {
		instStr = strings.ReplaceAll(instStr, opStr+" ", hexStr+" ")
	}
	instructions = strings.Split(instStr, " ")
	instStr = ""
	for _, inst := range instructions {
		if inst != "" {
			instStr = instStr + "0x" + inst + ", "
		}
	}
	instStr = instStr[:len(instStr)-2]
	instStr = string(append([]byte(instStr), 0x0a))

	return instStr
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	toks := GenerateTokens(data)

	insts := GenerateInstructions(toks)
	byteCode := GenerateBytecode(insts)
	fmt.Println(byteCode)
}
