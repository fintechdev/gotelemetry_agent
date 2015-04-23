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
const T_COLON = 57358
const T_OPEN_PARENS = 57359
const T_CLOSE_PARENS = 57360
const T_OPEN_BRACKET = 57361
const T_CLOSE_BRACKET = 57362
const T_OPEN_BRACE = 57363
const T_CLOSE_BRACE = 57364
const T_OR = 57365
const T_AND = 57366
const T_EQUAL = 57367
const T_NOT_EQUAL = 57368
const T_NEGATE = 57369
const T_IF = 57370
const T_ELSE = 57371
const T_FUNCTION_CALL = 57372
const T_UMINUS = 57373
const T_UPLUS = 57374

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
	"T_OPEN_BRACE",
	"T_CLOSE_BRACE",
	"T_OR",
	"T_AND",
	"T_EQUAL",
	"T_NOT_EQUAL",
	"T_NEGATE",
	"T_IF",
	"T_ELSE",
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
	17, 35,
	-2, 21,
	-1, 52,
	25, 0,
	26, 0,
	-2, 26,
	-1, 53,
	25, 0,
	26, 0,
	-2, 27,
}

const parserNprod = 44
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 163

var parserAct = []int{

	5, 4, 63, 12, 11, 19, 13, 14, 15, 22,
	21, 38, 3, 66, 9, 25, 10, 26, 27, 47,
	9, 43, 44, 45, 46, 69, 23, 20, 41, 48,
	49, 50, 51, 52, 53, 54, 55, 68, 42, 20,
	36, 67, 59, 60, 56, 61, 30, 31, 65, 36,
	25, 12, 11, 40, 39, 14, 15, 22, 21, 62,
	28, 29, 30, 31, 10, 36, 16, 18, 70, 65,
	72, 71, 24, 2, 23, 12, 11, 64, 39, 14,
	15, 22, 21, 37, 17, 8, 7, 6, 10, 28,
	29, 30, 31, 1, 36, 0, 26, 27, 23, 0,
	9, 0, 35, 34, 32, 33, 28, 29, 30, 31,
	0, 36, 57, 0, 58, 0, 0, 0, 20, 35,
	34, 32, 33, 28, 29, 30, 31, 0, 36, 28,
	29, 30, 31, 0, 36, 0, 35, 34, 32, 33,
	0, 0, 0, 34, 32, 33, 28, 29, 30, 31,
	0, 36, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 32, 33,
}
var parserPact = []int{

	-1, -1000, 11, -1000, -1000, 113, -1000, -1000, -1000, -1,
	47, -1000, -1000, 12, -1000, -1000, -1000, -1000, -1000, 22,
	47, 47, 47, 47, 2, -1000, 22, 12, 47, 47,
	47, 47, 47, 47, 47, 47, 38, 90, 96, -1000,
	-1000, 47, 47, 79, 25, 25, 25, 71, 34, 34,
	25, 25, 50, 50, 136, 119, -1000, -1000, -1000, 113,
	113, -16, 23, -1000, 9, 113, -7, -1000, 71, 47,
	-1000, -1000, 113,
}
var parserPgo = []int{

	0, 73, 93, 1, 12, 87, 86, 85, 0, 84,
	72, 67, 66, 59, 2,
}
var parserR1 = []int{

	0, 2, 2, 1, 1, 1, 1, 3, 4, 4,
	4, 5, 6, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 9, 10, 10, 11, 13, 13, 13,
	14, 14, 7, 7,
}
var parserR2 = []int{

	0, 0, 1, 2, 1, 1, 1, 3, 1, 1,
	1, 3, 3, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 4, 1, 1, 3, 3, 1, 0,
	3, 1, 3, 5,
}
var parserChk = []int{

	-1000, -2, -1, -4, -3, -8, -5, -6, -7, 21,
	17, 5, 4, 7, 8, 9, -12, -9, -11, 6,
	28, 11, 10, 27, -10, -4, 6, 7, 10, 11,
	12, 13, 25, 26, 24, 23, 15, -1, -8, 7,
	6, 16, 16, -8, -8, -8, -8, 17, -8, -8,
	-8, -8, -8, -8, -8, -8, 6, 22, 18, -8,
	-8, -3, -13, -14, 6, -8, 29, 18, 14, 16,
	-3, -14, -8,
}
var parserDef = []int{

	1, -2, 2, 4, 5, 6, 8, 9, 10, 0,
	0, 14, 15, 16, 17, 18, 19, 20, -2, 34,
	0, 0, 0, 0, 0, 3, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 16,
	34, 0, 0, 0, 30, 31, 32, 39, 22, 23,
	24, 25, -2, -2, 28, 29, 36, 7, 13, 12,
	11, 42, 0, 38, 34, 41, 0, 33, 0, 0,
	43, 37, 40,
}
var parserTok1 = []int{

	1,
}
var parserTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32,
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
		//line internal_parser.y:48
		{
			parserVAL.cmds = []command{}
		}
	case 2:
		//line internal_parser.y:50
		{
			for _, cmd := range parserVAL.cmds {
				parserlex.(*aslLexer).AddCommand(cmd)
			}

			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 3:
		//line internal_parser.y:60
		{
			parserVAL.cmds = append(parserVAL.cmds, parserS[parserpt-0].cmd)
		}
	case 4:
		//line internal_parser.y:62
		{
			parserVAL.cmds = []command{parserS[parserpt-0].cmd}
		}
	case 5:
		//line internal_parser.y:64
		{
			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 6:
		//line internal_parser.y:66
		{
			parserVAL.cmds = []command{newEvaluateCommand(parserS[parserpt-0].ex)}
		}
	case 7:
		//line internal_parser.y:70
		{
			parserVAL.cmds = parserS[parserpt-1].cmds
		}
	case 8:
		//line internal_parser.y:74
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 9:
		//line internal_parser.y:76
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 10:
		//line internal_parser.y:78
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 11:
		//line internal_parser.y:82
		{
			parserVAL.cmd = newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 12:
		//line internal_parser.y:86
		{
			parserVAL.cmd = newAssignCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 13:
		//line internal_parser.y:90
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 14:
		//line internal_parser.y:92
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 15:
		//line internal_parser.y:94
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 16:
		//line internal_parser.y:96
		{
			parserVAL.ex = newVariableExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 17:
		//line internal_parser.y:98
		{
			parserVAL.ex = newBooleanExpression(true, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 18:
		//line internal_parser.y:100
		{
			parserVAL.ex = newBooleanExpression(false, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 19:
		//line internal_parser.y:102
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 20:
		//line internal_parser.y:104
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 21:
		//line internal_parser.y:106
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 22:
		//line internal_parser.y:110
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 23:
		//line internal_parser.y:112
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 24:
		//line internal_parser.y:114
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 25:
		//line internal_parser.y:116
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 26:
		//line internal_parser.y:118
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 27:
		//line internal_parser.y:120
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 28:
		//line internal_parser.y:122
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 29:
		//line internal_parser.y:124
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 30:
		//line internal_parser.y:126
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 31:
		//line internal_parser.y:128
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 32:
		//line internal_parser.y:130
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-0].ex, booleanExpressionZero, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 33:
		//line internal_parser.y:134
		{
			parserVAL.ex = newFunctionCallExpression(parserS[parserpt-3].ex, parserS[parserpt-1].exl, parserS[parserpt-3].ex.line(), parserS[parserpt-3].ex.position())
		}
	case 34:
		//line internal_parser.y:138
		{
			parserVAL.ex = newPropertyExpression(newGlobalExpression(parserS[parserpt-0].t.line, parserS[parserpt-0].t.start), parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 35:
		//line internal_parser.y:140
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 36:
		//line internal_parser.y:144
		{
			parserVAL.ex = newPropertyExpression(parserS[parserpt-2].ex, parserS[parserpt-0].t.source, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 37:
		//line internal_parser.y:148
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 38:
		//line internal_parser.y:150
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 39:
		//line internal_parser.y:152
		{
			parserVAL.exl = map[string]expression{}
		}
	case 40:
		//line internal_parser.y:156
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 41:
		//line internal_parser.y:158
		{
			parserVAL.exi = parseArgument{singleUnnamedArgument, parserS[parserpt-0].ex}
		}
	case 42:
		//line internal_parser.y:162
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-1].ex, parserS[parserpt-0].cmds, []command{})
		}
	case 43:
		//line internal_parser.y:164
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-3].ex, parserS[parserpt-2].cmds, parserS[parserpt-0].cmds)
		}
	}
	goto parserstack /* stack new state and value */
}
