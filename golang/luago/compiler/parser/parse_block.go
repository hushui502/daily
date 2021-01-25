package parser

import (
	. "luago/compiler/ast"
	. "luago/compiler/lexer"
)

// block ::= {stat} [retstat]
func parseBlock(lexer *Lexer) *Block {
	return &Block{
		Stats:    parseStats(lexer),
		RetExps:  parseRetExps(lexer),
		LastLine: lexer.Line(),
	}
}

func parseStats(lexer *Lexer) []Stat {
	stats := make([]Stat, 0, 8)
	for !_isReturnOrBlockEnd(lexer.LookAhead()) {
		stat := parseStat(lexer)
		if _, ok := stat.(*EmptyStat); !ok {
			stats = append(stats, stat)
		}
	}

	return stats
}

// do block end
// while exp do block end
// repeat block until exp
// if exp then block ... end
// for ... end
// function ... end
// local function ... end
func _isReturnOrBlockEnd(tokenKind int) bool {
	switch tokenKind {
	case TOKEN_KW_RETURN, TOKEN_EOF, TOKEN_KW_END, TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return true
	}
	return false
}

func parseRetExps(lexer *Lexer) []Exp {
	if lexer.LookAhead() != TOKEN_KW_RETURN {
		return nil
	}
	lexer.NextToken()	// return
	
	switch lexer.LookAhead() {	// return 之后的值
	case TOKEN_EOF, TOKEN_KW_END, TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return []Exp{}
	case TOKEN_SEP_SEMI:	// return;
		lexer.NextToken()
		return []Exp{}
	default:			// return 1
		exps := parseExpList(lexer)
		if lexer.LookAhead() == TOKEN_SEP_SEMI {
			lexer.NextToken()
		}
		return exps
	}
}



