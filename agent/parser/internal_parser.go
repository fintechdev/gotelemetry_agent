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
const T_PLUS = 57350
const T_MINUS = 57351
const T_MULTIPLY = 57352
const T_DIVIDE = 57353
const T_COMMA = 57354
const T_DOT = 57355
const T_COLON = 57356
const T_OPEN_PARENS = 57357
const T_CLOSE_PARENS = 57358
const T_OPEN_BRACKET = 57359
const T_CLOSE_BRACKET = 57360
const T_FUNCTION_CALL = 57361
const T_UMINUS = 57362
const T_UPLUS = 57363

var parserToknames = []string{
	"T_STRING",
	"T_NUMBER",
	"T_IDENTIFIER",
	"T_VARIABLE",
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
	-1, 17,
	15, 21,
	-2, 18,
}

const parserNprod = 28
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 59

var parserAct = []int{

	39, 41, 37, 23, 24, 40, 25, 29, 9, 20,
	42, 26, 27, 28, 8, 14, 13, 19, 15, 12,
	11, 7, 30, 31, 32, 33, 10, 14, 13, 38,
	15, 12, 11, 25, 21, 22, 23, 24, 10, 25,
	5, 6, 35, 44, 43, 21, 22, 23, 24, 34,
	25, 1, 36, 17, 18, 16, 4, 3, 2,
}
var parserPact = []int{

	-1000, 34, -1000, -1000, -1000, 7, 0, 11, 11, 37,
	11, 11, 11, -1000, -1000, -1000, -1000, -1000, -8, -1000,
	37, 11, 11, 11, 11, 43, 26, 20, 20, 23,
	-7, -7, 20, 20, -1000, -1000, -11, -1000, -4, 37,
	-1000, 23, 11, -1000, 37,
}
var parserPgo = []int{

	0, 58, 57, 56, 0, 55, 54, 53, 52, 2,
	51,
}
var parserR1 = []int{

	0, 10, 10, 1, 1, 2, 3, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 5,
	6, 6, 7, 8, 8, 8, 9, 9,
}
var parserR2 = []int{

	0, 0, 2, 1, 1, 3, 3, 3, 3, 3,
	3, 3, 2, 2, 1, 1, 1, 1, 1, 4,
	1, 1, 3, 3, 1, 0, 3, 1,
}
var parserChk = []int{

	-1000, -10, -1, -2, -3, 6, 7, 14, 14, -4,
	15, 9, 8, 5, 4, 7, -5, -7, -6, 6,
	-4, 8, 9, 10, 11, 13, -4, -4, -4, 15,
	-4, -4, -4, -4, 6, 16, -8, -9, 6, -4,
	16, 12, 14, -9, -4,
}
var parserDef = []int{

	1, -2, 2, 3, 4, 0, 0, 0, 0, 5,
	0, 0, 0, 14, 15, 16, 17, -2, 0, 20,
	6, 0, 0, 0, 0, 0, 0, 12, 13, 25,
	8, 9, 10, 11, 22, 7, 0, 24, 20, 27,
	19, 0, 0, 23, 26,
}
var parserTok1 = []int{

	1,
}
var parserTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
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
		//line internal_parser.y:50
		{
			parserlex.(*aslLexer).AddCommand(newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex))
		}
	case 6:
		//line internal_parser.y:57
		{
			parserlex.(*aslLexer).AddCommand(newAssignCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex))
		}
	case 7:
		//line internal_parser.y:63
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 8:
		//line internal_parser.y:65
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 9:
		//line internal_parser.y:67
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 10:
		//line internal_parser.y:69
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 11:
		//line internal_parser.y:71
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 12:
		//line internal_parser.y:73
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 13:
		//line internal_parser.y:75
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 14:
		//line internal_parser.y:77
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 15:
		//line internal_parser.y:79
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 16:
		//line internal_parser.y:81
		{
			parserVAL.ex = newVariableExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 17:
		//line internal_parser.y:83
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 18:
		//line internal_parser.y:85
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 19:
		//line internal_parser.y:89
		{
			parserVAL.ex = newFunctionCallExpression(parserS[parserpt-3].ex, parserS[parserpt-1].exl, parserS[parserpt-3].ex.line(), parserS[parserpt-3].ex.position())
		}
	case 20:
		//line internal_parser.y:93
		{
			parserVAL.ex = newPropertyExpression(newGlobalExpression(parserS[parserpt-0].t.line, parserS[parserpt-0].t.start), parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 21:
		//line internal_parser.y:96
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 22:
		//line internal_parser.y:100
		{
			parserVAL.ex = newPropertyExpression(parserS[parserpt-2].ex, parserS[parserpt-0].t.source, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 23:
		//line internal_parser.y:104
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 24:
		//line internal_parser.y:106
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 25:
		//line internal_parser.y:108
		{
			parserVAL.exl = map[string]expression{}
		}
	case 26:
		//line internal_parser.y:112
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 27:
		//line internal_parser.y:115
		{
			parserVAL.exi = parseArgument{singleUnnamedArgument, parserS[parserpt-0].ex}
		}
	}
	goto parserstack /* stack new state and value */
}
