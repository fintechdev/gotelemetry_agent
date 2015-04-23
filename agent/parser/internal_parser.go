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
	yys  int
	cmds []command
	cmd  command
	ex   expression
	exl  map[string]expression
	exi  parseArgument
	t    token
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
const T_ASSIGN = 57358
const T_COLON = 57359
const T_OPEN_PARENS = 57360
const T_CLOSE_PARENS = 57361
const T_OPEN_BRACKET = 57362
const T_CLOSE_BRACKET = 57363
const T_OPEN_BRACE = 57364
const T_CLOSE_BRACE = 57365
const T_OR = 57366
const T_AND = 57367
const T_EQUAL = 57368
const T_NOT_EQUAL = 57369
const T_NEGATE = 57370
const T_GREATER_THAN = 57371
const T_LESS_THAN = 57372
const T_GREATER_THAN_OR_EQUAL = 57373
const T_LESS_THAN_OR_EQUAL = 57374
const T_IF = 57375
const T_ELSE = 57376
const T_TERMINATOR = 57377
const T_COMMENT = 57378
const T_FUNCTION_CALL = 57379
const T_UMINUS = 57380
const T_UPLUS = 57381

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
	"T_ASSIGN",
	"T_COLON",
	"T_OPEN_PARENS",
	"T_CLOSE_PARENS",
	"T_OPEN_BRACKET",
	"T_CLOSE_BRACKET",
	"T_OPEN_BRACE",
	"T_CLOSE_BRACE",
	"T_OR",
	"T_AND",
	"T_EQUAL",
	"T_NOT_EQUAL",
	"T_NEGATE",
	"T_GREATER_THAN",
	"T_LESS_THAN",
	"T_GREATER_THAN_OR_EQUAL",
	"T_LESS_THAN_OR_EQUAL",
	"T_IF",
	"T_ELSE",
	"T_TERMINATOR",
	"T_COMMENT",
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
	-1, 22,
	18, 41,
	-2, 23,
	-1, 59,
	26, 0,
	27, 0,
	-2, 28,
	-1, 60,
	26, 0,
	27, 0,
	-2, 29,
}

const parserNprod = 50
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 288

var parserAct = []int{

	8, 4, 74, 3, 29, 28, 27, 17, 16, 49,
	48, 18, 19, 24, 23, 47, 50, 77, 11, 54,
	46, 15, 79, 80, 51, 52, 53, 78, 45, 43,
	67, 25, 55, 56, 57, 58, 59, 60, 61, 62,
	63, 64, 65, 66, 73, 2, 69, 70, 27, 71,
	20, 31, 32, 33, 34, 76, 43, 44, 22, 17,
	16, 12, 13, 18, 19, 24, 23, 35, 36, 31,
	32, 33, 34, 15, 43, 26, 21, 7, 68, 81,
	76, 83, 82, 25, 33, 34, 6, 43, 14, 5,
	9, 10, 17, 16, 12, 13, 18, 19, 24, 23,
	1, 0, 0, 0, 0, 0, 15, 0, 0, 0,
	11, 0, 31, 32, 33, 34, 25, 43, 0, 0,
	0, 14, 0, 9, 10, 0, 38, 37, 35, 36,
	0, 39, 41, 40, 42, 0, 0, 30, 17, 16,
	12, 13, 18, 19, 24, 23, 0, 0, 0, 0,
	0, 0, 15, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 25, 31, 32, 33, 34, 14, 43, 9,
	10, 0, 72, 0, 0, 0, 0, 38, 37, 35,
	36, 0, 39, 41, 40, 42, 31, 32, 33, 34,
	0, 43, 0, 0, 0, 0, 0, 0, 11, 0,
	38, 37, 35, 36, 0, 39, 41, 40, 42, 31,
	32, 33, 34, 0, 43, 0, 0, 31, 32, 33,
	34, 0, 43, 38, 37, 35, 36, 0, 39, 41,
	40, 42, 37, 35, 36, 0, 39, 41, 40, 42,
	31, 32, 33, 34, 0, 43, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 35, 36, 0, 39,
	41, 40, 42, 17, 16, 75, 48, 18, 19, 24,
	23, 0, 0, 0, 0, 0, 0, 15, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 25,
}
var parserPact = []int{

	88, -1000, 134, -1000, -1000, -30, -31, -1000, 102, -1000,
	-1000, 88, 12, 4, 3, 3, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 3, 3, 3, 1, -1000, -1000, -1000,
	-1000, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 24, 55, 3, 3, 176, -1000, -1000,
	153, 14, 14, 14, 259, 72, 72, 14, 14, 59,
	59, 230, 207, 41, 41, 41, 41, -1000, -1000, 199,
	199, -17, -1000, 8, -1000, 6, 199, -4, -1000, 259,
	3, -1000, -1000, 199,
}
var parserPgo = []int{

	0, 45, 100, 1, 3, 89, 86, 77, 0, 76,
	75, 58, 50, 44, 2,
}
var parserR1 = []int{

	0, 2, 2, 1, 1, 1, 3, 4, 4, 4,
	4, 4, 4, 5, 6, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 9,
	10, 10, 11, 13, 13, 13, 14, 14, 7, 7,
}
var parserR2 = []int{

	0, 0, 1, 2, 1, 1, 3, 2, 2, 1,
	2, 1, 1, 3, 3, 3, 1, 1, 1, 1,
	1, 1, 1, 1, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 2, 2, 2, 4,
	1, 1, 3, 3, 1, 0, 3, 1, 3, 5,
}
var parserChk = []int{

	-1000, -2, -1, -4, -3, -5, -6, -7, -8, 35,
	36, 22, 6, 7, 33, 18, 5, 4, 8, 9,
	-12, -9, -11, 11, 10, 28, -10, -4, 35, 35,
	35, 10, 11, 12, 13, 26, 27, 25, 24, 29,
	31, 30, 32, 15, -1, 16, 16, -8, 7, 6,
	-8, -8, -8, -8, 18, -8, -8, -8, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, 6, 23, -8,
	-8, -3, 19, -13, -14, 6, -8, 34, 19, 14,
	17, -3, -14, -8,
}
var parserDef = []int{

	1, -2, 2, 4, 5, 0, 0, 9, 0, 11,
	12, 0, 40, 18, 0, 0, 16, 17, 19, 20,
	21, 22, -2, 0, 0, 0, 0, 3, 7, 8,
	10, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 18, 40,
	0, 36, 37, 38, 45, 24, 25, 26, 27, -2,
	-2, 30, 31, 32, 33, 34, 35, 42, 6, 13,
	14, 48, 15, 0, 44, 40, 47, 0, 39, 0,
	0, 49, 43, 46,
}
var parserTok1 = []int{

	1,
}
var parserTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39,
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

	case 1:
		//line internal_parser.y:51
		{
			parserVAL.cmds = []command{}
		}
	case 2:
		//line internal_parser.y:53
		{
			for _, cmd := range parserVAL.cmds {
				parserlex.(*aslLexer).AddCommand(cmd)
			}

			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 3:
		//line internal_parser.y:63
		{
			parserVAL.cmds = append(parserVAL.cmds, parserS[parserpt-0].cmd)
		}
	case 4:
		//line internal_parser.y:65
		{
			parserVAL.cmds = []command{parserS[parserpt-0].cmd}
		}
	case 5:
		//line internal_parser.y:67
		{
			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 6:
		//line internal_parser.y:71
		{
			parserVAL.cmds = parserS[parserpt-1].cmds
		}
	case 7:
		//line internal_parser.y:75
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 8:
		//line internal_parser.y:77
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 9:
		//line internal_parser.y:79
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 10:
		//line internal_parser.y:81
		{
			parserVAL.cmd = newEvaluateCommand(parserS[parserpt-1].ex)
		}
	case 11:
		//line internal_parser.y:83
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 12:
		//line internal_parser.y:85
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 13:
		//line internal_parser.y:89
		{
			parserVAL.cmd = newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 14:
		//line internal_parser.y:93
		{
			parserVAL.cmd = newAssignCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 15:
		//line internal_parser.y:97
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 16:
		//line internal_parser.y:99
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 17:
		//line internal_parser.y:101
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 18:
		//line internal_parser.y:103
		{
			parserVAL.ex = newVariableExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 19:
		//line internal_parser.y:105
		{
			parserVAL.ex = newBooleanExpression(true, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 20:
		//line internal_parser.y:107
		{
			parserVAL.ex = newBooleanExpression(false, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 21:
		//line internal_parser.y:109
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 22:
		//line internal_parser.y:111
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 23:
		//line internal_parser.y:113
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 24:
		//line internal_parser.y:117
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 25:
		//line internal_parser.y:119
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 26:
		//line internal_parser.y:121
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 27:
		//line internal_parser.y:123
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 28:
		//line internal_parser.y:125
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 29:
		//line internal_parser.y:127
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 30:
		//line internal_parser.y:129
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 31:
		//line internal_parser.y:131
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 32:
		//line internal_parser.y:133
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 33:
		//line internal_parser.y:135
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 34:
		//line internal_parser.y:137
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 35:
		//line internal_parser.y:139
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 36:
		//line internal_parser.y:141
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 37:
		//line internal_parser.y:143
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 38:
		//line internal_parser.y:145
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-0].ex, booleanExpressionZero, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 39:
		//line internal_parser.y:149
		{
			parserVAL.ex = newFunctionCallExpression(parserS[parserpt-3].ex, parserS[parserpt-1].exl, parserS[parserpt-3].ex.line(), parserS[parserpt-3].ex.position())
		}
	case 40:
		//line internal_parser.y:153
		{
			parserVAL.ex = newPropertyExpression(newGlobalExpression(parserS[parserpt-0].t.line, parserS[parserpt-0].t.start), parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 41:
		//line internal_parser.y:155
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 42:
		//line internal_parser.y:159
		{
			parserVAL.ex = newPropertyExpression(parserS[parserpt-2].ex, parserS[parserpt-0].t.source, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 43:
		//line internal_parser.y:163
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 44:
		//line internal_parser.y:165
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 45:
		//line internal_parser.y:167
		{
			parserVAL.exl = map[string]expression{}
		}
	case 46:
		//line internal_parser.y:171
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 47:
		//line internal_parser.y:173
		{
			parserVAL.exi = parseArgument{singleUnnamedArgument, parserS[parserpt-0].ex}
		}
	case 48:
		//line internal_parser.y:177
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-1].ex, parserS[parserpt-0].cmds, []command{})
		}
	case 49:
		//line internal_parser.y:179
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-3].ex, parserS[parserpt-2].cmds, parserS[parserpt-0].cmds)
		}
	}
	goto parserstack /* stack new state and value */
}
