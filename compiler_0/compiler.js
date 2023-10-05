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
        return c >= '0'.charCodeAt(0) && c <= '9'.charCodeAt(0)
    }

    function isAplabet(c) {
        c = c.charCodeAt(0);
        return (c >= 'a'.charCodeAt(0) && c <= 'z'.charCodeAt(0)) || (c >= 'A'.charCodeAt(0) && c <= 'Z'.charCodeAt(0)) || (c === '_'.charCodeAt(0))
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

Compiler.createAST = function (tokens) {
    return tokens;
};

Compiler.compileASTAndGenerateInst = function (ast) { return ast; };

Compiler.generateBytecode = function (inst) { return inst; };

Compiler.compile = function (source_code) {
    var bytecode = Compiler.compileASTAndGenerateInst(Compiler.createAST(Compiler.generateTokens(source_code)));
    console.log(bytecode);
    return bytecode;
};