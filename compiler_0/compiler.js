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

    var globalVarP = 0;

    var stackVar = [];
    var stackVarP = 0;

    var current = 0;

    var funcTable = [];
    var varTable = [];

    var currentScope = 1;

    function emitInst(op, p1, p2, p3) {
        allInst.push([op, p1, p2, p3]);
        instP += 4;
    }

    function emitInstLI(r, v) {
        v = ('00000000' + v.toString(16)).slice(-8);
        r = ('00' + r.toString(16)).slice(-2);
        emitInst('lliu', r, v.slice(6, 8), v.slice(4, 6));
        emitInst('lui', r, v.slice(2, 4), v.slice(0, 2));
    }

    function emitInitInst() {
        // Return Value: 0x2000_0000
        emitInstLI(1, 0x2000_0000);

        // Global Variables: 0x2001_0000
        emitInstLI(2, 0x2001_0000);

        // Local Variables: 0x3fff_ffff
        emitInstLI(3, 0x3fff_ffff);

        // Jump to main function
        emitInstLI(8, 0);

        emitInst('jal', '00', '00', '08');
    }

    function advance() {
        var t = tokens[current];
        current++;
        return t;
    }

    function consume(T) {
        if (peek().T === T) {
            advance();
        } else {
            console.log('Error while consuming token', peek());
        };
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
        return { ident: '', isLValue: false, loc: 0, scope: 0, varType: varTypeCreate(), };
    }

    function funcInfoCreate() {
        return { ident: '', loc: 0, parm: [], retVarType: varTypeCreate(), };
    }
    function getVarType() {
        var varType = varTypeCreate();

        if (peek().T == 'MUL') {
            varType.isPointer = true;
            varType.pointerLevel = 1;
            consume('MUL');
            while (peek().T === 'MUL') {
                varType.pointerLevel++;
                consume('MUL');
            }
        } else if (peek().T == 'LBRACK') {
            varType.isArray = true;
            consume('LBRACK');
            if (peek().T !== 'RBRACK') {
                varType.arraySize = parseIntLiteral(advance().S);
            }
            consume('RBRACK');
        }

        varType.varSize = 4;

        consume('IDENT');

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
            consume('FUNC');
            var funcInfo = funcInfoCreate();

            funcInfo.ident = advance().S;
            funcInfo.loc = instP;
            consume('LPAREN');

            while (peek().T !== 'RPAREN') {
                funcInfo.parm.push(getVarInfo());
                if (peek().T === 'COMMA') {
                    consume('COMMA');
                }
            }

            consume('RPAREN');

            if (peek().T !== 'LBRACE') {
                funcInfo.retVarType = getVarType();
            }

            funcTable.push(funcInfo);

            consume('LBRACE');
            consume('NEW_LINE');

            currentScope++;
        },
        'VAR': function () {
            consume('VAR');
            var varInfo = getVarInfo();
            varInfo.scope = currentScope;
            varInfo.isLValue = true;

            if (currentScope !== 1) {
            }

            varTable.push(varInfo);
            consume('NEW_LINE');
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