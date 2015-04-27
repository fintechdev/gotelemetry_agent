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
	exa  []expression
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
	-1, 23,
	18, 44,
	-2, 24,
	-1, 62,
	26, 0,
	27, 0,
	-2, 31,
	-1, 63,
	26, 0,
	27, 0,
	-2, 32,
}

const parserNprod = 53
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 298

var parserAct = []int{

	8, 4, 79, 3, 30, 29, 28, 18, 17, 50,
	49, 19, 20, 25, 24, 48, 51, 53, 82, 11,
	57, 15, 77, 16, 86, 54, 55, 56, 47, 76,
	46, 26, 44, 58, 59, 60, 61, 62, 63, 64,
	65, 66, 67, 68, 69, 70, 85, 72, 73, 28,
	74, 84, 52, 78, 32, 33, 34, 35, 81, 44,
	2, 18, 17, 12, 13, 19, 20, 25, 24, 21,
	36, 37, 45, 34, 35, 15, 44, 16, 83, 23,
	71, 27, 22, 7, 87, 26, 81, 89, 88, 6,
	14, 5, 9, 10, 18, 17, 12, 13, 19, 20,
	25, 24, 32, 33, 34, 35, 1, 44, 15, 0,
	16, 0, 11, 0, 0, 0, 0, 0, 26, 0,
	0, 0, 0, 14, 0, 9, 10, 18, 17, 12,
	13, 19, 20, 25, 24, 0, 0, 0, 0, 0,
	0, 15, 0, 16, 0, 0, 0, 32, 33, 34,
	35, 26, 44, 0, 0, 0, 14, 0, 9, 10,
	0, 39, 38, 36, 37, 0, 40, 42, 41, 43,
	0, 0, 31, 32, 33, 34, 35, 0, 44, 0,
	0, 0, 75, 0, 0, 0, 0, 39, 38, 36,
	37, 0, 40, 42, 41, 43, 32, 33, 34, 35,
	0, 44, 0, 0, 0, 0, 0, 0, 11, 0,
	39, 38, 36, 37, 0, 40, 42, 41, 43, 32,
	33, 34, 35, 0, 44, 0, 0, 32, 33, 34,
	35, 0, 44, 39, 38, 36, 37, 0, 40, 42,
	41, 43, 38, 36, 37, 0, 40, 42, 41, 43,
	32, 33, 34, 35, 0, 44, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 36, 37, 0, 40,
	42, 41, 43, 18, 17, 80, 49, 19, 20, 25,
	24, 0, 0, 0, 0, 0, 0, 15, 0, 16,
	0, 0, 0, 0, 0, 0, 0, 26,
}
var parserPact = []int{

	90, -1000, 123, -1000, -1000, -30, -31, -1000, 137, -1000,
	-1000, 90, 14, 12, 3, 3, 3, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 3, 3, 3, 2, -1000, -1000,
	-1000, -1000, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 39, 57, 3, 3, 186, -1000,
	-1000, 163, 8, 209, 17, 17, 17, 269, 61, 61,
	17, 17, 92, 92, 240, 217, 44, 44, 44, 44,
	-1000, -1000, 209, 209, -16, -1000, -1000, 3, 32, -1000,
	7, 209, -3, 209, -1000, 269, 3, -1000, -1000, 209,
}
var parserPgo = []int{

	0, 60, 106, 1, 3, 91, 89, 83, 0, 82,
	81, 79, 69, 53, 2, 52,
}
var parserR1 = []int{

	0, 2, 2, 1, 1, 1, 3, 4, 4, 4,
	4, 4, 4, 5, 6, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 15, 15, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 9, 10, 10, 11, 13, 13, 13, 14,
	14, 7, 7,
}
var parserR2 = []int{

	0, 0, 1, 2, 1, 1, 3, 2, 2, 1,
	2, 1, 1, 3, 3, 3, 3, 1, 1, 1,
	1, 1, 1, 1, 1, 3, 1, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 2,
	2, 2, 4, 1, 1, 3, 3, 1, 0, 3,
	1, 3, 5,
}
var parserChk = []int{

	-1000, -2, -1, -4, -3, -5, -6, -7, -8, 35,
	36, 22, 6, 7, 33, 18, 20, 5, 4, 8,
	9, -12, -9, -11, 11, 10, 28, -10, -4, 35,
	35, 35, 10, 11, 12, 13, 26, 27, 25, 24,
	29, 31, 30, 32, 15, -1, 16, 16, -8, 7,
	6, -8, -15, -8, -8, -8, -8, 18, -8, -8,
	-8, -8, -8, -8, -8, -8, -8, -8, -8, -8,
	6, 23, -8, -8, -3, 19, 21, 14, -13, -14,
	6, -8, 34, -8, 19, 14, 17, -3, -14, -8,
}
var parserDef = []int{

	1, -2, 2, 4, 5, 0, 0, 9, 0, 11,
	12, 0, 43, 19, 0, 0, 0, 17, 18, 20,
	21, 22, 23, -2, 0, 0, 0, 0, 3, 7,
	8, 10, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 19,
	43, 0, 0, 26, 39, 40, 41, 48, 27, 28,
	29, 30, -2, -2, 33, 34, 35, 36, 37, 38,
	45, 6, 13, 14, 51, 15, 16, 0, 0, 47,
	43, 50, 0, 25, 42, 0, 0, 52, 46, 49,
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
		//line internal_parser.y:53
		{
			parserVAL.cmds = []command{}
		}
	case 2:
		//line internal_parser.y:55
		{
			for _, cmd := range parserVAL.cmds {
				parserlex.(*aslLexer).AddCommand(cmd)
			}

			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 3:
		//line internal_parser.y:65
		{
			if parserS[parserpt-0].cmd != nil {
				parserVAL.cmds = append(parserVAL.cmds, parserS[parserpt-0].cmd)
			}
		}
	case 4:
		//line internal_parser.y:67
		{
			if parserS[parserpt-0].cmd != nil {
				parserVAL.cmds = []command{parserS[parserpt-0].cmd}
			} else {
				parserVAL.cmds = []command{}
			}
		}
	case 5:
		//line internal_parser.y:69
		{
			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 6:
		//line internal_parser.y:73
		{
			parserVAL.cmds = parserS[parserpt-1].cmds
		}
	case 7:
		//line internal_parser.y:77
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 8:
		//line internal_parser.y:79
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 9:
		//line internal_parser.y:81
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 10:
		//line internal_parser.y:83
		{
			parserVAL.cmd = newEvaluateCommand(parserS[parserpt-1].ex)
		}
	case 11:
		//line internal_parser.y:85
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 12:
		//line internal_parser.y:87
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 13:
		//line internal_parser.y:91
		{
			parserVAL.cmd = newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 14:
		//line internal_parser.y:95
		{
			parserVAL.cmd = newAssignCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 15:
		//line internal_parser.y:99
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 16:
		//line internal_parser.y:101
		{
			parserVAL.ex = newArrayExpression(parserS[parserpt-1].exa, parserS[parserpt-2].t.line, parserS[parserpt-2].t.start)
		}
	case 17:
		//line internal_parser.y:103
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 18:
		//line internal_parser.y:105
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 19:
		//line internal_parser.y:107
		{
			parserVAL.ex = newVariableExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 20:
		//line internal_parser.y:109
		{
			parserVAL.ex = newBooleanExpression(true, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 21:
		//line internal_parser.y:111
		{
			parserVAL.ex = newBooleanExpression(false, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 22:
		//line internal_parser.y:113
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 23:
		//line internal_parser.y:115
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 24:
		//line internal_parser.y:117
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 25:
		//line internal_parser.y:121
		{
			parserVAL.exa = append(parserVAL.exa, parserS[parserpt-0].ex)
		}
	case 26:
		//line internal_parser.y:123
		{
			parserVAL.exa = []expression{parserS[parserpt-0].ex}
		}
	case 27:
		//line internal_parser.y:127
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 28:
		//line internal_parser.y:129
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 29:
		//line internal_parser.y:131
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 30:
		//line internal_parser.y:133
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 31:
		//line internal_parser.y:135
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 32:
		//line internal_parser.y:137
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 33:
		//line internal_parser.y:139
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 34:
		//line internal_parser.y:141
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 35:
		//line internal_parser.y:143
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 36:
		//line internal_parser.y:145
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 37:
		//line internal_parser.y:147
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 38:
		//line internal_parser.y:149
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 39:
		//line internal_parser.y:151
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 40:
		//line internal_parser.y:153
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 41:
		//line internal_parser.y:155
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-0].ex, booleanExpressionZero, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 42:
		//line internal_parser.y:159
		{
			parserVAL.ex = newFunctionCallExpression(parserS[parserpt-3].ex, parserS[parserpt-1].exl, parserS[parserpt-3].ex.line(), parserS[parserpt-3].ex.position())
		}
	case 43:
		//line internal_parser.y:163
		{
			parserVAL.ex = newPropertyExpression(newGlobalExpression(parserS[parserpt-0].t.line, parserS[parserpt-0].t.start), parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 44:
		//line internal_parser.y:165
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 45:
		//line internal_parser.y:169
		{
			parserVAL.ex = newPropertyExpression(parserS[parserpt-2].ex, parserS[parserpt-0].t.source, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 46:
		//line internal_parser.y:173
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 47:
		//line internal_parser.y:175
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 48:
		//line internal_parser.y:177
		{
			parserVAL.exl = map[string]expression{}
		}
	case 49:
		//line internal_parser.y:181
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 50:
		//line internal_parser.y:183
		{
			parserVAL.exi = parseArgument{singleUnnamedArgument, parserS[parserpt-0].ex}
		}
	case 51:
		//line internal_parser.y:187
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-1].ex, parserS[parserpt-0].cmds, []command{})
		}
	case 52:
		//line internal_parser.y:189
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-3].ex, parserS[parserpt-2].cmds, parserS[parserpt-0].cmds)
		}
	}
	goto parserstack /* stack new state and value */
}
