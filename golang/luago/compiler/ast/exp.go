package ast


/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	 prefixexp | tableconstructor | exp binop exp | unop exp

prefixexp ::= var | functioncall | ‘(’ exp ‘)’

var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name

functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
*/

type Exp interface {}

type NilExp struct{ Line int }    // nil
type TrueExp struct{ Line int }   // true
type FalseExp struct{ Line int }  // false
type VarargExp struct{ Line int } // ...


// Numeral
type IntegerExp struct {
	Line int
	Val  int64
}
type FloatExp struct {
	Line int
	Val  float64
}

// LiteralString
type StringExp struct {
	Line int
	Str  string
}

// 一元运算符表达式
type UnopExp struct {
	Line int
	Op int
	Exp Exp
}

type BinopExp struct {
	Line int
	Op int
	Exp1 Exp
	Exp2 Exp
}

type ConcatExp struct {
	Line int
	Exps []Exp
}

type TableConstructorExp struct {
	Line int		// line of {
	LastLine int	// line of }
	KeyExps []Exp
	ValExps []Exp
}

// functiondef ::= function funcbody
// funcbody ::= ‘(’ [parlist] ‘)’ block end
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
// namelist ::= Name {‘,’ Name}
type FuncDefExp struct {
	Line int
	LastLine int 	// line of end
	ParList []string	// params
	IsVararg bool
	Block *Block
}

/*
prefixexp ::= Name |
              ‘(’ exp ‘)’ |
              prefixexp ‘[’ exp ‘]’ |
              prefixexp ‘.’ Name |
              prefixexp ‘:’ Name args |
              prefixexp args
*/

type NameExp struct {
	Line int
	Name string
}

type ParentsExp struct {
	Exp Exp
}

type TableAccessExp struct {
	LastLine int	// line of ]
	PrefixExp Exp
	KeyExp Exp
}

type FuncCallExp struct {
	Line int
	LastLine int
	PrefixExp Exp
	NameExp *StringExp
	Args []Exp
}