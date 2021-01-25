package compiler

import (
	. "luago/binchunk"
 	. "luago/compiler/codegen"
 	. "luago/compiler/parser"
)

func Compile(chunk, chunkName string) *Prototype {
	ast := Parse(chunk, chunkName)

	return GenProto(ast)
}
