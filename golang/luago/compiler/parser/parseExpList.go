package parser

import (
	. "luago/compiler/ast"
	. "luago/compiler/lexer"
)

func parseExpList(lexer *Lexer) []Exp {
	exps := make([]Exp, 0, 4)
	exps = append(exps, parseExp(lexer))		// exp
	for lexer.LookAhead() == TOKEN_SEP_COMMA {	// {
		lexer.NextToken()						// ,
		exps = append(exps, parseExp(lexer))	// }
	}

	return exps
}


