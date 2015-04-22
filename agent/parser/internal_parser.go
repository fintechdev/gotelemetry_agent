//line internal_parser.y:2
package parser

import __yyfmt__ "fmt"

//line internal_parser.y:3
type parseArgument struct {
	key   string
	value expression
}

const singleUnnamedArgument = "----"

//line internal_parser.y:14
type parserSymType struct {
	yys int
	cmd command
	ex  expression
	exl map[string]expression
	exi parseArgument
	t   token
}

const T_STRING = 57346
const T_NUMBER = 57347
const T_IDENTIFIER = 57348
const T_VARIABLE = 57349
const T_TRUE = 57350
const T_FALSE = 57351
const T_PLUS = 57352
const T_MINUS = 57353
const T_MULTIPLY = 57354
const T_DIVIDE = 57355
const T_COMMA = 57356
const T_DOT = 57357
const T_COLON = 57358
const T_OPEN_PARENS = 57359
const T_CLOSE_PARENS = 57360
const T_OPEN_BRACKET = 57361
const T_CLOSE_BRACKET = 57362
const T_OR = 57363
const T_AND = 57364
const T_EQUAL = 57365
const T_NOT_EQUAL = 57366
const T_NEGATE = 57367
const T_FUNCTION_CALL = 57368
const T_UMINUS = 57369
const T_UPLUS = 57370

var parserToknames = []string{
	"T_STRING",
	"T_NUMBER",
	"T_IDENTIFIER",
	"T_VARIABLE",
	"T_TRUE",
	"T_FALSE",
	"T_PLUS",
	"T_MINUS",
	"T_MULTIPLY",
	"T_DIVIDE",
	"T_COMMA",
	"T_DOT",
	"T_COLON",
	"T_OPEN_PARENS",
	"T_CLOSE_PARENS",
	"T_OPEN_BRACKET",
	"T_CLOSE_BRACKET",
	"T_OR",
	"T_AND",
	"T_EQUAL",
	"T_NOT_EQUAL",
	"T_NEGATE",
	"T_FUNCTION_CALL",
	"T_UMINUS",
	"T_UPLUS",
}
var parserStatenames = []string{}

const parserEofCode = 1
const parserErrCode = 2
const parserMaxDepth = 200

//line yacctab:1
var parserExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 18,
	17, 29,
	-2, 15,
	-1, 43,
	23, 0,
	24, 0,
	-2, 20,
	-1, 44,
	23, 0,
	24, 0,
	-2, 21,
}

const parserNprod = 36
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 123

var parserAct = []int{

	52, 54, 50, 27, 28, 53, 33, 38, 9, 24,
	55, 34, 25, 26, 27, 28, 8, 33, 47, 7,
	35, 36, 37, 33, 5, 6, 39, 40, 41, 42,
	43, 44, 45, 46, 12, 11, 23, 13, 14, 15,
	20, 19, 1, 49, 16, 18, 22, 10, 17, 4,
	3, 2, 0, 0, 0, 21, 57, 56, 12, 11,
	51, 13, 14, 15, 20, 19, 0, 0, 0, 0,
	0, 10, 25, 26, 27, 28, 0, 33, 0, 21,
	48, 0, 0, 32, 31, 29, 30, 25, 26, 27,
	28, 0, 33, 0, 0, 0, 0, 0, 32, 31,
	29, 30, 25, 26, 27, 28, 0, 33, 25, 26,
	27, 28, 0, 33, 31, 29, 30, 0, 0, 0,
	0, 29, 30,
}
var parserPact = []int{

	-1000, 18, -1000, -1000, -1000, 3, 0, 30, 30, 77,
	30, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 30,
	30, 30, -10, -1000, 77, 30, 30, 30, 30, 30,
	30, 30, 30, 12, 62, 8, 8, 8, 54, -9,
	-9, 8, 8, 2, 2, 98, 92, -1000, -1000, -13,
	-1000, -6, 77, -1000, 54, 30, -1000, 77,
}
var parserPgo = []int{

	0, 51, 50, 49, 0, 48, 46, 45, 44, 43,
	2, 42,
}
var parserR1 = []int{

	0, 11, 11, 1, 1, 2, 3, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 5, 6, 6,
	7, 9, 9, 9, 10, 10,
}
var parserR2 = []int{

	0, 0, 2, 1, 1, 3, 3, 3, 1, 1,
	1, 1, 1, 1, 1, 1, 3, 3, 3, 3,
	3, 3, 3, 3, 2, 2, 2, 4, 1, 1,
	3, 3, 1, 0, 3, 1,
}
var parserChk = []int{

	-1000, -11, -1, -2, -3, 6, 7, 16, 16, -4,
	17, 5, 4, 7, 8, 9, -8, -5, -7, 11,
	10, 25, -6, 6, -4, 10, 11, 12, 13, 23,
	24, 22, 21, 15, -4, -4, -4, -4, 17, -4,
	-4, -4, -4, -4, -4, -4, -4, 6, 18, -9,
	-10, 6, -4, 18, 14, 16, -10, -4,
}
var parserDef = []int{

	1, -2, 2, 3, 4, 0, 0, 0, 0, 5,
	0, 8, 9, 10, 11, 12, 13, 14, -2, 0,
	0, 0, 0, 28, 6, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 24, 25, 26, 33, 16,
	17, 18, 19, -2, -2, 22, 23, 30, 7, 0,
	32, 28, 35, 27, 0, 0, 31, 34,
}
var parserTok1 = []int{

	1,
}
var parserTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28,
}
var parserTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var parserDebug = 0

type parserLexer interface {
	Lex(lval *parserSymType) int
	Error(s string)
}

const parserFlag = -1000

func parserTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(parserToknames) {
		if parserToknames[c-4] != "" {
			return parserToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func parserStatname(s int) string {
	if s >= 0 && s < len(parserStatenames) {
		if parserStatenames[s] != "" {
			return parserStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func parserlex1(lex parserLexer, lval *parserSymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = parserTok1[0]
		goto out
	}
	if char < len(parserTok1) {
		c = parserTok1[char]
		goto out
	}
	if char >= parserPrivate {
		if char < parserPrivate+len(parserTok2) {
			c = parserTok2[char-parserPrivate]
			goto out
		}
	}
	for i := 0; i < len(parserTok3); i += 2 {
		c = parserTok3[i+0]
		if c == char {
			c = parserTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = parserTok2[1] /* unknown char */
	}
	if parserDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", parserTokname(c), uint(char))
	}
	return c
}

func parserParse(parserlex parserLexer) int {
	var parsern int
	var parserlval parserSymType
	var parserVAL parserSymType
	parserS := make([]parserSymType, parserMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	parserstate := 0
	parserchar := -1
	parserp := -1
	goto parserstack

ret0:
	return 0

ret1:
	return 1

parserstack:
	/* put a state and value onto the stack */
	if parserDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", parserTokname(parserchar), parserStatname(parserstate))
	}

	parserp++
	if parserp >= len(parserS) {
		nyys := make([]parserSymType, len(parserS)*2)
		copy(nyys, parserS)
		parserS = nyys
	}
	parserS[parserp] = parserVAL
	parserS[parserp].yys = parserstate

parsernewstate:
	parsern = parserPact[parserstate]
	if parsern <= parserFlag {
		goto parserdefault /* simple state */
	}
	if parserchar < 0 {
		parserchar = parserlex1(parserlex, &parserlval)
	}
	parsern += parserchar
	if parsern < 0 || parsern >= parserLast {
		goto parserdefault
	}
	parsern = parserAct[parsern]
	if parserChk[parsern] == parserchar { /* valid shift */
		parserchar = -1
		parserVAL = parserlval
		parserstate = parsern
		if Errflag > 0 {
			Errflag--
		}
		goto parserstack
	}

parserdefault:
	/* default state action */
	parsern = parserDef[parserstate]
	if parsern == -2 {
		if parserchar < 0 {
			parserchar = parserlex1(parserlex, &parserlval)
		}

		/* look through exception table */
		xi := 0
		for {
			if parserExca[xi+0] == -1 && parserExca[xi+1] == parserstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			parsern = parserExca[xi+0]
			if parsern < 0 || parsern == parserchar {
				break
			}
		}
		parsern = parserExca[xi+1]
		if parsern < 0 {
			goto ret0
		}
	}
	if parsern == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			parserlex.Error("syntax error")
			Nerrs++
			if parserDebug >= 1 {
				__yyfmt__.Printf("%s", parserStatname(parserstate))
				__yyfmt__.Printf(" saw %s\n", parserTokname(parserchar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for parserp >= 0 {
				parsern = parserPact[parserS[parserp].yys] + parserErrCode
				if parsern >= 0 && parsern < parserLast {
					parserstate = parserAct[parsern] /* simulate a shift of "error" */
					if parserChk[parserstate] == parserErrCode {
						goto parserstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if parserDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", parserS[parserp].yys)
				}
				parserp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if parserDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", parserTokname(parserchar))
			}
			if parserchar == parserEofCode {
				goto ret1
			}
			parserchar = -1
			goto parsernewstate /* try again in the same state */
		}
	}

	/* reduction by production parsern */
	if parserDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", parsern, parserStatname(parserstate))
	}

	parsernt := parsern
	parserpt := parserp
	_ = parserpt // guard against "declared and not used"

	parserp -= parserR2[parsern]
	parserVAL = parserS[parserp+1]

	/* consult goto table to find next state */
	parsern = parserR1[parsern]
	parserg := parserPgo[parsern]
	parserj := parserg + parserS[parserp].yys + 1

	if parserj >= parserLast {
		parserstate = parserAct[parserg]
	} else {
		parserstate = parserAct[parserj]
		if parserChk[parserstate] != -parsern {
			parserstate = parserAct[parserg]
		}
	}
	// dummy call; replaced with literal code
	switch parsernt {

	case 3:
		parserVAL.cmd = parserS[parserpt-0].cmd
	case 4:
		parserVAL.cmd = parserS[parserpt-0].cmd
	case 5:
		//line internal_parser.y:52
		{
			parserlex.(*aslLexer).AddCommand(newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex))
		}
	case 6:
		//line internal_parser.y:58
		{
			parserlex.(*aslLexer).AddCommand(newAssignCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex))
		}
	case 7:
		//line internal_parser.y:64
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 8:
		//line internal_parser.y:66
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 9:
		//line internal_parser.y:68
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 10:
		//line internal_parser.y:70
		{
			parserVAL.ex = newVariableExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 11:
		//line internal_parser.y:72
		{
			parserVAL.ex = newBooleanExpression(true, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 12:
		//line internal_parser.y:74
		{
			parserVAL.ex = newBooleanExpression(false, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 13:
		//line internal_parser.y:76
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 14:
		//line internal_parser.y:78
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 15:
		//line internal_parser.y:80
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 16:
		//line internal_parser.y:84
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 17:
		//line internal_parser.y:86
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 18:
		//line internal_parser.y:88
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 19:
		//line internal_parser.y:90
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 20:
		//line internal_parser.y:92
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 21:
		//line internal_parser.y:94
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 22:
		//line internal_parser.y:96
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 23:
		//line internal_parser.y:98
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 24:
		//line internal_parser.y:100
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 25:
		//line internal_parser.y:102
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 26:
		//line internal_parser.y:104
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-0].ex, booleanExpressionZero, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 27:
		//line internal_parser.y:108
		{
			parserVAL.ex = newFunctionCallExpression(parserS[parserpt-3].ex, parserS[parserpt-1].exl, parserS[parserpt-3].ex.line(), parserS[parserpt-3].ex.position())
		}
	case 28:
		//line internal_parser.y:112
		{
			parserVAL.ex = newPropertyExpression(newGlobalExpression(parserS[parserpt-0].t.line, parserS[parserpt-0].t.start), parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 29:
		//line internal_parser.y:114
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 30:
		//line internal_parser.y:118
		{
			parserVAL.ex = newPropertyExpression(parserS[parserpt-2].ex, parserS[parserpt-0].t.source, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 31:
		//line internal_parser.y:122
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 32:
		//line internal_parser.y:124
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 33:
		//line internal_parser.y:126
		{
			parserVAL.exl = map[string]expression{}
		}
	case 34:
		//line internal_parser.y:130
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 35:
		//line internal_parser.y:132
		{
			parserVAL.exi = parseArgument{singleUnnamedArgument, parserS[parserpt-0].ex}
		}
	}
	goto parserstack /* stack new state and value */
}
