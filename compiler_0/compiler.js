var Compiler = {};

Compiler.compile = function (tokens) {
    var current = 0;

    var assembly = [];

    var tableGlobalFunc = [];
    var tableGlobalVar = [];
    var tableLocalVar = [];

    var computStack = [];

    var computOpStack = [];

    var currentScopeLevel = 1;

    var instPOrigin = 0x10000000;
    var instP = instPOrigin;

    var globalVarPOrigin = 0x20010000;
    var globalVarP = globalVarPOrigin;

    var localVarPOrigin = 0x3fffffff
    var localVarP = localVarPOrigin;

    var token;

    function advance() {
        var token = tokens[current];
        current++;
        return token;
    }

    function peek() {
        return tokens[current];
    }

    function compileFunc() {
        advance();
        var token = advance();
        var funcInfo = {
            Ident: '',
            parameters: [],
            returnValueType: {
                varSize: 0,
            },
            isLocatable: true,
            location: instP,
        };

        funcInfo.Ident = token.S;

        advance();

        while (true) {
            if (peek().T === 'RPAREN') {
                advance();
                break;
            }
            funcInfo.parameters.push(getVarNameAndType());
            advance();
        }

        if (peek().T !== 'LBRACE') {
            funcInfo.returnValueType = getVarType();
        }

        advance();
        advance();

        tableGlobalFunc.push(funcInfo);
        currentScopeLevel = 2;
        localVarP = localVarPOrigin;
    }

    function compileVar() {
        advance();
        var varInfo = getVarNameAndType();

        if (currentScopeLevel === 1) {
            varInfo.location = globalVarP;
            varInfo.isLocatable = true;
            tableGlobalVar.push(varInfo);
            globalVarP += varInfo.varType.varSize;
        } else {
            varInfo.location = localVarP;
            varInfo.isLocatable = true;
            tableLocalVar.push(varInfo);
            localVarP += varInfo.varType.varSize;
        }
    }

    function getVarNameAndType() {
        var token = advance();
        var varInfo = {
            Ident: '',
            varType: {},
            isLocatable: false,
            location: 0,
        };
        varInfo.Ident = token.S;
        varInfo.varType = getVarType();
        return varInfo;
    }

    function parseIntLiteral() {
        var intLiteralInfo = {
            varSize: 0,
            isSigned: false,
            varValue: 0,
            isLocatable: false,
            location: 0,
        };

        var token = advance();
        var s = token.S;
        if (s.slice(-2) === 'u8') {
            s = s.slice(0, -2);
            intLiteralInfo.varSize = 1;
        } else if (s.slice(-3) === 'i32') {
            s = s.slice(0, -3);
            intLiteralInfo.varSize = 4;
            intLiteralInfo.isSigned = true;
        }
        intLiteralInfo.varValue = parseInt(s.replace(/_/g, ''));
        return intLiteralInfo;
    }

    function getVarType() {
        var varType = {
            varSize: 0,
            isSigned: false,
            isPointer: false,
            pointerLevel: 0,
            isArray: false,
            arraySize: 0,
        };
        if (peek().T === 'MUL') {
            varType.isPointer = true;
            varType.pointerLevel = 1;
            while (peek().T === 'MUL') {
                advance();
                varType.pointerLevel++;
            }
        } else if (peek().T === 'LBRACK') {
            varType.isArray = true;
            varType.arraySize = parseIntLiteral().varValue;
            advance();
        }

        var token = advance();

        if (token.S === 'u8') {
            varType.varSize = 1;
        } else if (token.S === 'i32') {
            varType.varSize = 4;
            varType.isSigned = true;
        }
        return varType;
    }

    function compileSatement() {
        function findVar(ident) {
            for (var i = tableLocalVar.length - 1; i >= 0; i--) {
                if (tableLocalVar[i].Ident === ident) {
                    return tableLocalVar[i];
                }
            }
            for (var i = tableGlobalVar.length - 1; i >= 0; i--) {
                if (tableGlobalVar[i].Ident === ident) {
                    return tableGlobalVar[i];
                }
            }
            for (var i = tableGlobalFunc.length - 1; i >= 0; i--) {
                if (tableGlobalFunc[i].Ident === ident) {
                    return tableGlobalFunc[i];
                }
            }
        }

        var token = advance();

        while (peek().T !== 'NEW_LINE') {
            if (token.T === 'IDENT') {
                computStack.push(findVar(token.S));
            } else if (token.T === 'INT') {
                computStack.push(findVar(token.S));
            }
            token = advance();
        }
    }

    function emitInst(inst) {
        assembly.push(inst);
        instP += 4;
    }

    function emitInitInst() {
        // Memory map:
        // 0x0002_0000 Ecall Parameters
        // 0x0003_0000 Ecall Return Value
        // 0x1000_0000 Instructions
        // 0x2000_0000 Return Value
        // 0x2001_0000 Global Variables
        // 0x3fff_ffff Local Variables

        emitInst('lliu 01 00 00'); // Return Value: 0x2000_0000
        emitInst('lui  01 00 20');

        emitInst('lliu 02 00 00'); // Global Variables: 0x2001_0000
        emitInst('lui  02 01 20');

        emitInst('lliu 03 ff ff'); // Local Variables: 0x3fff_ffff
        emitInst('lui  03 ff 3f');

        emitInst('lliu 08 00 00'); // Jump to main function
        emitInst('lui  08 00 00');
        emitInst('jal  00 00 08');
    }

    function runMainLoop() {
        var token;
        while (true) {
            token = peek();
            if (token.T === 'EOF') {
                break;
            } else if (token.T === 'FUNC') {
                compileFunc();
            } else if (token.T === 'VAR') {
                compileVar();
            } else if (token.T === 'IDENT') {
                compileSatement();
            } else {
                advance();
            }
        }
    }

    function logDebugInfo() {
        console.log(tableGlobalFunc);
        console.log(tableGlobalVar);
        console.log(tableLocalVar);
        console.log(instP);
    }

    emitInitInst();
    runMainLoop();
    logDebugInfo();

    return assembly;
};