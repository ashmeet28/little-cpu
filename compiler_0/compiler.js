var Compiler = {};

Compiler.generateTokens = function (source_code) {
    var current;
    var start;

    function isAtEnd() {
        return current === source_code.length;
    }

    function advance() {
        var c = source_code[current];
        current++;
        return c;
    }

    function peek() {
        return source_code[current];
    }

    function isDigit(c) {
        c = c.charCodeAt(0);
        return c >= '0'.charCodeAt(0) && c <= '9'.charCodeAt(0);
    }

    function isAplabet(c) {
        c = c.charCodeAt(0);
        return (c >= 'a'.charCodeAt(0) && c <= 'z'.charCodeAt(0)) || (c >= 'A'.charCodeAt(0) && c <= 'Z'.charCodeAt(0)) || (c === '_'.charCodeAt(0));
    }

    function isKeyword(s) {
        var keywords = [
            'WHILE',
            'BREAK',
            'CONTINUE',
            'IF',
            'ELSE',
            'FUNC',
            'RETURN',
            'VAR',
        ];

        return keywords.includes(s.toUpperCase());
    }
    function scanToken() {
        start = current
        var token = {
            T: '',
            S: '',
        }

        if (isAtEnd()) {
            token.T = 'EOF';
            return token;
        }

        var c = advance()
        if (c === '+') {
            token.T = 'ADD';
        } else if (c === '-') {
            token.T = 'SUB'
        } else if (c === '*') {
            token.T = 'MUL'
        } else if (c === '/') {
            if (peek() === '/') {
                advance()
                while (peek() != '\n') {
                    advance()
                }
                token.T = 'COMMENT'
            } else {
                token.T = 'QUO'
            }
        } else if (c === '%') {
            token.T = 'REM'
        } else if (c === '&') {
            if (peek() === '&') {
                advance()
                token.T = 'LAND'
            } else {
                token.T = 'AND'
            }
        } else if (c === '|') {
            if (peek() === '|') {
                advance()
                token.T = 'LOR'
            } else {
                token.T = 'OR'
            }
        } else if (c === '^') {
            token.T = 'XOR'
        } else if (c === '<') {
            if (peek() === '<') {
                advance()
                token.T = 'SHL'
            } else if (peek() === '=') {
                advance()
                token.T = 'LEQ'
            } else {
                token.T = 'LSS'
            }
        } else if (c === '>') {
            if (peek() === '>') {
                advance()
                token.T = 'SHR'
            } else if (peek() === '=') {
                advance()
                token.T = 'GEQ'
            } else {
                token.T = 'GTR'
            }
        } else if (c === '=') {
            if (peek() === '=') {
                advance()
                token.T = 'EQL'
            } else {
                token.T = 'ASSIGN'
            }
        } else if (c === '!') {
            if (peek() === '=') {
                advance()
                token.T = 'NEQ'
            } else {
                token.T = 'NOT'
            }
        } else if (c === '(') {
            token.T = 'LPAREN'
        } else if (c === '[') {
            token.T = 'LBRACK'
        } else if (c === '{') {
            token.T = 'LBRACE'
        } else if (c === ',') {
            token.T = 'COMMA'
        } else if (c === '.') {
            token.T = 'PERIOD'
        } else if (c === ')') {
            token.T = 'RPAREN'
        } else if (c === ']') {
            token.T = 'RBRACK'
        } else if (c === '}') {
            token.T = 'RBRACE'
        } else if (c === ';') {
            token.T = 'SEMICOLON'
        } else if (c === ':') {
            token.T = 'COLON'
        } else if (c === ' ') {
            while (peek() === ' ') {
                advance()
            }
            token.T = 'SPACE'
        } else if (c === '\n') {
            while (peek() === '\n') {
                advance()
            }
            token.T = 'NEW_LINE'
        } else if (isDigit(c)) {
            while (isAplabet(peek()) || isDigit(peek())) {
                advance()
            }
            token.T = 'INT';
            token.S = source_code.slice(start, current)
        } else if (isAplabet(c)) {
            while (isAplabet(peek()) || isDigit(peek())) {
                advance()
            }
            token.T = 'IDENT';
            var s = source_code.slice(start, current)
            if (isKeyword(s)) {
                token.T = s.toUpperCase();
            } else {
                token.S = s;
            }
        }

        return token;
    }

    function scanTokens() {
        var tokens = [];
        current = 0;
        start = 0;
        var token;
        while (true) {
            token = scanToken();
            if (token.T === 'SPACE' || token.T === 'COMMENT') {
                continue;
            }
            tokens.push(token);
            if (token.T === 'EOF') {
                break;
            }
        }
        return tokens
    }

    return scanTokens();
};

Compiler.generateInst = function (tokens) {
    var allInst = [];
    var instP = 0;
    var instPOffset = 0x10000000;

    var globalVarP = 0;
    var globalVarPOffset = 0x20010000;

    var current = 0;

    var funcTable = [];
    var varTable = [];

    var currentScope = 1;

    function emitInst(op, p1, p2, p3) {
        allInst.push([op, p1, p2, p3]);
        instP += 4;
    }

    function emitInitInst() {
        emitInst('lliu', '01', '00', '00'); // Return Value: 0x2000_0000
        emitInst('lui', '01', '00', '20');

        emitInst('lliu', '02', '00', '00'); // Global Variables: 0x2001_0000
        emitInst('lui', '02', '01', '20');

        emitInst('lliu', '03', 'ff', 'ff'); // Local Variables: 0x3fff_ffff
        emitInst('lui', '03', 'ff', '3f');

        emitInst('lliu', '08', '00', '00'); // Jump to main function
        emitInst('lui', '08', '00', '00');
        emitInst('jal', '00', '00', '08');
    }

    function advance() {
        var t = tokens[current];
        current++;
        return t;
    }

    function peek() {
        return tokens[current];
    }

    function parseIntLiteral(s) {
        return parseInt(s.replace(/_/g, ''));
    }

    function varTypeCreate() {
        return { varSize: 0, isPointer: false, pointerLevel: 0, isArray: false, arraySize: 0, };
    }

    function varInfoCreate() {
        return { ident: '', loc: 0, scope: 0, varType: varTypeCreate(), };
    }

    function funcInfoCreate() {
        return { ident: '', loc: 0, parm: [], retVarType: varInfoCreate(), };
    }
    function getVarType() {
        var varType = varTypeCreate();

        if (peek().T == 'MUL') {
            varType.isPointer = true;
            varType.pointerLevel = 1;
            advance();
            while (peek().T === 'MUL') {
                varType.pointerLevel++;
                advance();
            }
        } else if (peek().T == 'LBRACK') {
            varType.isArray = true;
            advance();
            if (peek().T !== 'RBRACK') {
                varType.arraySize = parseIntLiteral(advance().S);
                advance();
            }
        }

        if (peek().T === 'IDENT' && peek().S === 'int') {
            varType.varSize = 4;
        }

        advance();

        return varType;
    }

    function getVarInfo() {
        var varInfo = varInfoCreate();
        varInfo.ident = advance().S;
        varInfo.varType = getVarType();
        return varInfo;
    }

    var functionsOfCompilingTokens = {
        'FUNC': function () {
            advance();
            var funcInfo = funcInfoCreate();


            funcInfo.ident = advance().S;
            funcInfo.loc = instP + instPOffset;
            advance();

            while (peek().T !== 'RPAREN') {
                funcInfo.parm.push(getVarInfo());
                if (peek().T === 'COMMA') {
                    advance();
                }
            }

            advance();

            if (peek().T !== 'LBRACK') {
                funcInfo.retVarType = getVarType();
            }

            funcTable.push(funcInfo);

            advance();
            advance();

            currentScope++;
        },
        'VAR': function () {
            advance();
            var varInfo = getVarInfo();
            varInfo.scope = currentScope;
            varTable.push(varInfo);
            advance();
        },
    };

    function compileTokens() {
        while (peek().T !== 'EOF') {
            if (functionsOfCompilingTokens[peek().T] !== undefined) {
                functionsOfCompilingTokens[peek().T]();
            } else {
                console.log('Unhandled token', peek());
                advance();
            }
        }
    }

    emitInitInst();
    compileTokens();
    console.log(funcTable);
    console.log(varTable);
    return allInst;
};

Compiler.compile = function (source_code) {
    var inst = Compiler.generateInst(Compiler.generateTokens(source_code));
    return inst;
};