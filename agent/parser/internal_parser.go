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
const T_WHILE = 57379
const T_NULL = 57380
const T_FUNCTION_CALL = 57381
const T_UMINUS = 57382
const T_UPLUS = 57383

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
	"T_WHILE",
	"T_NULL",
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
	-1, 26,
	18, 52,
	-2, 28,
	-1, 72,
	26, 0,
	27, 0,
	-2, 39,
	-1, 73,
	26, 0,
	27, 0,
	-2, 40,
}

const parserNprod = 62
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 410

var parserAct = []int{

	9, 4, 94, 51, 83, 34, 3, 33, 100, 31,
	36, 37, 38, 39, 97, 48, 55, 58, 59, 62,
	88, 53, 84, 104, 92, 82, 40, 41, 63, 64,
	65, 91, 66, 38, 39, 84, 48, 68, 69, 70,
	71, 72, 73, 74, 75, 76, 77, 78, 79, 54,
	53, 48, 67, 103, 85, 86, 31, 87, 102, 80,
	89, 36, 37, 38, 39, 61, 48, 96, 36, 37,
	38, 39, 50, 48, 93, 43, 42, 40, 41, 24,
	44, 46, 45, 47, 26, 99, 35, 98, 49, 2,
	30, 25, 8, 101, 20, 19, 13, 14, 21, 22,
	28, 27, 105, 7, 96, 107, 106, 6, 17, 5,
	18, 1, 32, 81, 0, 0, 0, 0, 29, 0,
	0, 0, 0, 15, 0, 10, 11, 16, 23, 20,
	19, 13, 14, 21, 22, 28, 27, 0, 0, 0,
	0, 0, 0, 17, 0, 18, 0, 12, 0, 0,
	0, 0, 0, 29, 0, 0, 0, 0, 15, 0,
	10, 11, 16, 23, 20, 19, 52, 14, 21, 22,
	28, 27, 0, 0, 0, 0, 0, 0, 17, 0,
	18, 0, 12, 0, 0, 0, 0, 0, 29, 0,
	0, 0, 0, 15, 0, 10, 11, 16, 23, 20,
	19, 13, 14, 21, 22, 28, 27, 0, 0, 0,
	0, 0, 0, 17, 0, 18, 0, 32, 0, 0,
	0, 0, 0, 29, 0, 0, 0, 0, 15, 0,
	10, 11, 16, 23, 20, 19, 57, 56, 21, 22,
	28, 27, 0, 0, 0, 0, 0, 0, 17, 0,
	18, 60, 32, 0, 0, 0, 0, 0, 29, 20,
	19, 57, 56, 21, 22, 28, 27, 0, 23, 0,
	0, 0, 0, 17, 0, 18, 0, 32, 0, 0,
	0, 0, 0, 29, 20, 19, 95, 56, 21, 22,
	28, 27, 0, 23, 0, 0, 0, 0, 17, 0,
	18, 0, 32, 0, 0, 0, 0, 0, 29, 0,
	36, 37, 38, 39, 0, 48, 0, 0, 23, 90,
	0, 0, 0, 0, 43, 42, 40, 41, 0, 44,
	46, 45, 47, 36, 37, 38, 39, 0, 48, 0,
	0, 0, 0, 0, 0, 88, 0, 43, 42, 40,
	41, 0, 44, 46, 45, 47, 36, 37, 38, 39,
	0, 48, 0, 0, 36, 37, 38, 39, 0, 48,
	43, 42, 40, 41, 0, 44, 46, 45, 47, 42,
	40, 41, 0, 44, 46, 45, 47, 36, 37, 38,
	39, 0, 48, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 40, 41, 0, 44, 46, 45, 47,
}
var parserPact = []int{

	125, -1000, 195, -1000, -1000, -28, -30, -1000, -1000, 51,
	-1000, -1000, 160, 34, 33, 255, 255, 255, 230, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 255, 255, 255,
	14, -1000, 46, -1000, -1000, -1000, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 53, 90,
	-10, -1000, 5, 255, 255, 323, -1000, -1000, 323, 300,
	-1000, 10, 346, 36, 36, 36, 280, 18, 21, 21,
	36, 36, 58, 58, 377, 354, 0, 0, 0, 0,
	-1000, -1000, -9, 46, 255, 346, 346, -26, 125, -1000,
	-1000, -1000, 255, 39, -1000, 6, 346, -1000, -1000, 346,
	-2, 346, -1000, 280, 255, -1000, -1000, 346,
}
var parserPgo = []int{

	0, 88, 111, 1, 6, 109, 107, 103, 92, 0,
	91, 90, 84, 79, 74, 72, 2, 3, 65,
}
var parserR1 = []int{

	0, 2, 2, 1, 1, 1, 3, 4, 4, 4,
	4, 4, 4, 4, 5, 6, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 18,
	18, 15, 15, 15, 17, 13, 13, 13, 13, 13,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	10, 11, 11, 12, 14, 14, 14, 16, 16, 7,
	7, 8,
}
var parserR2 = []int{

	0, 0, 1, 2, 1, 1, 3, 2, 2, 1,
	1, 2, 1, 1, 3, 3, 3, 2, 3, 4,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	1, 3, 1, 0, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2, 2,
	4, 1, 1, 3, 3, 1, 0, 3, 1, 3,
	5, 3,
}
var parserChk = []int{

	-1000, -2, -1, -4, -3, -5, -6, -7, -8, -9,
	35, 36, 22, 6, 7, 33, 37, 18, 20, 5,
	4, 8, 9, 38, -13, -10, -12, 11, 10, 28,
	-11, -4, 22, 35, 35, 35, 10, 11, 12, 13,
	26, 27, 25, 24, 29, 31, 30, 32, 15, -1,
	-15, -17, 6, 16, 16, -9, 7, 6, -9, -9,
	21, -18, -9, -9, -9, -9, 18, 6, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	6, 23, 35, 14, 17, -9, -9, -3, 22, -3,
	19, 21, 14, -14, -16, 6, -9, 23, -17, -9,
	34, -9, 19, 14, 17, -3, -16, -9,
}
var parserDef = []int{

	1, -2, 2, 4, 5, 0, 0, 9, 10, 0,
	12, 13, 33, 51, 22, 0, 0, 0, 0, 20,
	21, 23, 24, 25, 26, 27, -2, 0, 0, 0,
	0, 3, 33, 7, 8, 11, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 32, 51, 0, 0, 0, 22, 51, 0, 0,
	17, 0, 30, 47, 48, 49, 56, 0, 35, 36,
	37, 38, -2, -2, 41, 42, 43, 44, 45, 46,
	53, 6, 0, 0, 0, 14, 15, 59, 0, 61,
	16, 18, 0, 0, 55, 51, 58, 19, 31, 34,
	0, 29, 50, 0, 0, 60, 54, 57,
}
var parserTok1 = []int{

	1,
}
var parserTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
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
		//line internal_parser.y:55
		{
			parserVAL.cmds = []command{}
		}
	case 2:
		//line internal_parser.y:57
		{
			for _, cmd := range parserVAL.cmds {
				parserlex.(*aslLexer).AddCommand(cmd)
			}

			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 3:
		//line internal_parser.y:67
		{
			if parserS[parserpt-0].cmd != nil {
				parserVAL.cmds = append(parserVAL.cmds, parserS[parserpt-0].cmd)
			}
		}
	case 4:
		//line internal_parser.y:69
		{
			if parserS[parserpt-0].cmd != nil {
				parserVAL.cmds = []command{parserS[parserpt-0].cmd}
			} else {
				parserVAL.cmds = []command{}
			}
		}
	case 5:
		//line internal_parser.y:71
		{
			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 6:
		//line internal_parser.y:75
		{
			parserVAL.cmds = parserS[parserpt-1].cmds
		}
	case 7:
		//line internal_parser.y:79
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 8:
		//line internal_parser.y:81
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 9:
		//line internal_parser.y:83
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 10:
		//line internal_parser.y:85
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 11:
		//line internal_parser.y:87
		{
			parserVAL.cmd = newEvaluateCommand(parserS[parserpt-1].ex)
		}
	case 12:
		//line internal_parser.y:89
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 13:
		//line internal_parser.y:91
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 14:
		//line internal_parser.y:95
		{
			parserVAL.cmd = newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 15:
		//line internal_parser.y:99
		{
			parserVAL.cmd = newAssignCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 16:
		//line internal_parser.y:103
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 17:
		//line internal_parser.y:105
		{
			parserVAL.ex = newArrayExpression([]interface{}{}, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 18:
		//line internal_parser.y:107
		{
			parserVAL.ex = newArrayExpression(parserS[parserpt-1].exa, parserS[parserpt-2].t.line, parserS[parserpt-2].t.start)
		}
	case 19:
		//line internal_parser.y:109
		{
			parserVAL.ex = newMapExpression(parserS[parserpt-2].exl, parserS[parserpt-3].t.line, parserS[parserpt-3].t.start)
		}
	case 20:
		//line internal_parser.y:111
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 21:
		//line internal_parser.y:113
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 22:
		//line internal_parser.y:115
		{
			parserVAL.ex = newVariableExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 23:
		//line internal_parser.y:117
		{
			parserVAL.ex = newBooleanExpression(true, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 24:
		//line internal_parser.y:119
		{
			parserVAL.ex = newBooleanExpression(false, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 25:
		//line internal_parser.y:121
		{
			parserVAL.ex = newNullExpression(nil, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 26:
		//line internal_parser.y:123
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 27:
		//line internal_parser.y:125
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 28:
		//line internal_parser.y:127
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 29:
		//line internal_parser.y:131
		{
			parserVAL.exa = append(parserVAL.exa, parserS[parserpt-0].ex)
		}
	case 30:
		//line internal_parser.y:133
		{
			parserVAL.exa = []expression{parserS[parserpt-0].ex}
		}
	case 31:
		//line internal_parser.y:137
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 32:
		//line internal_parser.y:139
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 33:
		//line internal_parser.y:141
		{
			parserVAL.exl = map[string]expression{}
		}
	case 34:
		//line internal_parser.y:145
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 35:
		//line internal_parser.y:150
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 36:
		//line internal_parser.y:152
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 37:
		//line internal_parser.y:154
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 38:
		//line internal_parser.y:156
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 39:
		//line internal_parser.y:158
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 40:
		//line internal_parser.y:160
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 41:
		//line internal_parser.y:162
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 42:
		//line internal_parser.y:164
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 43:
		//line internal_parser.y:166
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 44:
		//line internal_parser.y:168
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 45:
		//line internal_parser.y:170
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 46:
		//line internal_parser.y:172
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 47:
		//line internal_parser.y:174
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 48:
		//line internal_parser.y:176
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 49:
		//line internal_parser.y:178
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-0].ex, booleanExpressionZero, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 50:
		//line internal_parser.y:182
		{
			parserVAL.ex = newFunctionCallExpression(parserS[parserpt-3].ex, parserS[parserpt-1].exl, parserS[parserpt-3].ex.line(), parserS[parserpt-3].ex.position())
		}
	case 51:
		//line internal_parser.y:186
		{
			parserVAL.ex = newPropertyExpression(newGlobalExpression(parserS[parserpt-0].t.line, parserS[parserpt-0].t.start), parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 52:
		//line internal_parser.y:188
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 53:
		//line internal_parser.y:192
		{
			parserVAL.ex = newPropertyExpression(parserS[parserpt-2].ex, parserS[parserpt-0].t.source, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 54:
		//line internal_parser.y:196
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 55:
		//line internal_parser.y:198
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 56:
		//line internal_parser.y:200
		{
			parserVAL.exl = map[string]expression{}
		}
	case 57:
		//line internal_parser.y:204
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 58:
		//line internal_parser.y:206
		{
			parserVAL.exi = parseArgument{singleUnnamedArgument, parserS[parserpt-0].ex}
		}
	case 59:
		//line internal_parser.y:210
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-1].ex, parserS[parserpt-0].cmds, []command{})
		}
	case 60:
		//line internal_parser.y:212
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-3].ex, parserS[parserpt-2].cmds, parserS[parserpt-0].cmds)
		}
	case 61:
		//line internal_parser.y:216
		{
			parserVAL.cmd = newWhileLoopCommand(parserS[parserpt-1].ex, parserS[parserpt-0].cmds)
		}
	}
	goto parserstack /* stack new state and value */
}
