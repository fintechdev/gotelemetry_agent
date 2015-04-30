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
const T_FUNCTION_CALL = 57380
const T_UMINUS = 57381
const T_UPLUS = 57382

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
	-1, 25,
	18, 50,
	-2, 26,
	-1, 70,
	26, 0,
	27, 0,
	-2, 37,
	-1, 71,
	26, 0,
	27, 0,
	-2, 38,
}

const parserNprod = 60
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 365

var parserAct = []int{

	9, 4, 92, 50, 81, 33, 3, 32, 98, 30,
	95, 86, 64, 90, 52, 82, 54, 57, 58, 60,
	89, 101, 102, 82, 53, 80, 100, 61, 62, 63,
	35, 36, 37, 38, 52, 47, 66, 67, 68, 69,
	70, 71, 72, 73, 74, 75, 76, 77, 37, 38,
	47, 47, 65, 83, 84, 30, 85, 78, 59, 87,
	48, 2, 49, 91, 23, 94, 20, 19, 56, 55,
	21, 22, 27, 26, 25, 29, 24, 8, 7, 6,
	17, 5, 18, 97, 31, 96, 1, 0, 0, 0,
	28, 99, 20, 19, 13, 14, 21, 22, 27, 26,
	103, 0, 94, 105, 104, 0, 17, 0, 18, 0,
	31, 79, 0, 0, 0, 0, 28, 0, 0, 0,
	0, 15, 0, 10, 11, 16, 20, 19, 13, 14,
	21, 22, 27, 26, 0, 0, 35, 36, 37, 38,
	17, 47, 18, 0, 12, 0, 0, 0, 0, 0,
	28, 0, 39, 40, 0, 15, 0, 10, 11, 16,
	20, 19, 51, 14, 21, 22, 27, 26, 0, 0,
	0, 0, 0, 0, 17, 0, 18, 0, 12, 0,
	0, 0, 0, 0, 28, 0, 0, 0, 0, 15,
	0, 10, 11, 16, 20, 19, 13, 14, 21, 22,
	27, 26, 0, 0, 0, 0, 0, 0, 17, 0,
	18, 0, 31, 0, 35, 36, 37, 38, 28, 47,
	0, 0, 0, 15, 0, 10, 11, 16, 42, 41,
	39, 40, 0, 43, 45, 44, 46, 0, 0, 34,
	35, 36, 37, 38, 0, 47, 0, 0, 0, 88,
	0, 0, 0, 0, 42, 41, 39, 40, 0, 43,
	45, 44, 46, 35, 36, 37, 38, 0, 47, 0,
	0, 0, 0, 0, 0, 86, 0, 42, 41, 39,
	40, 0, 43, 45, 44, 46, 35, 36, 37, 38,
	0, 47, 0, 0, 35, 36, 37, 38, 0, 47,
	42, 41, 39, 40, 0, 43, 45, 44, 46, 41,
	39, 40, 0, 43, 45, 44, 46, 35, 36, 37,
	38, 0, 47, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 39, 40, 0, 43, 45, 44, 46,
	20, 19, 93, 55, 21, 22, 27, 26, 0, 0,
	0, 0, 0, 0, 17, 0, 18, 0, 31, 0,
	0, 0, 0, 0, 28,
}
var parserPact = []int{

	122, -1000, 190, -1000, -1000, -28, -30, -1000, -1000, 204,
	-1000, -1000, 156, 18, 8, 62, 62, 62, 62, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 62, 62, 62, -6,
	-1000, 46, -1000, -1000, -1000, 62, 62, 62, 62, 62,
	62, 62, 62, 62, 62, 62, 62, 51, 88, -10,
	-1000, -2, 62, 62, 253, -1000, -1000, 253, 230, -1,
	276, 35, 35, 35, 336, 6, 36, 36, 35, 35,
	20, 20, 307, 284, 126, 126, 126, 126, -1000, -1000,
	-13, 46, 62, 276, 276, -26, 122, -1000, -1000, -1000,
	62, 7, -1000, 5, 276, -1000, -1000, 276, -11, 276,
	-1000, 336, 62, -1000, -1000, 276,
}
var parserPgo = []int{

	0, 60, 86, 1, 6, 81, 79, 78, 77, 0,
	76, 75, 74, 64, 63, 62, 2, 3, 58,
}
var parserR1 = []int{

	0, 2, 2, 1, 1, 1, 3, 4, 4, 4,
	4, 4, 4, 4, 5, 6, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 18, 18, 15,
	15, 15, 17, 13, 13, 13, 13, 13, 13, 13,
	13, 13, 13, 13, 13, 13, 13, 13, 10, 11,
	11, 12, 14, 14, 14, 16, 16, 7, 7, 8,
}
var parserR2 = []int{

	0, 0, 1, 2, 1, 1, 3, 2, 2, 1,
	1, 2, 1, 1, 3, 3, 3, 3, 4, 1,
	1, 1, 1, 1, 1, 1, 1, 3, 1, 3,
	1, 0, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 2, 2, 2, 4, 1,
	1, 3, 3, 1, 0, 3, 1, 3, 5, 3,
}
var parserChk = []int{

	-1000, -2, -1, -4, -3, -5, -6, -7, -8, -9,
	35, 36, 22, 6, 7, 33, 37, 18, 20, 5,
	4, 8, 9, -13, -10, -12, 11, 10, 28, -11,
	-4, 22, 35, 35, 35, 10, 11, 12, 13, 26,
	27, 25, 24, 29, 31, 30, 32, 15, -1, -15,
	-17, 6, 16, 16, -9, 7, 6, -9, -9, -18,
	-9, -9, -9, -9, 18, 6, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, 6, 23,
	35, 14, 17, -9, -9, -3, 22, -3, 19, 21,
	14, -14, -16, 6, -9, 23, -17, -9, 34, -9,
	19, 14, 17, -3, -16, -9,
}
var parserDef = []int{

	1, -2, 2, 4, 5, 0, 0, 9, 10, 0,
	12, 13, 31, 49, 21, 0, 0, 0, 0, 19,
	20, 22, 23, 24, 25, -2, 0, 0, 0, 0,
	3, 31, 7, 8, 11, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	30, 49, 0, 0, 0, 21, 49, 0, 0, 0,
	28, 45, 46, 47, 54, 0, 33, 34, 35, 36,
	-2, -2, 39, 40, 41, 42, 43, 44, 51, 6,
	0, 0, 0, 14, 15, 57, 0, 59, 16, 17,
	0, 0, 53, 49, 56, 18, 29, 32, 0, 27,
	48, 0, 0, 58, 52, 55,
}
var parserTok1 = []int{

	1,
}
var parserTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40,
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
		//line internal_parser.y:54
		{
			parserVAL.cmds = []command{}
		}
	case 2:
		//line internal_parser.y:56
		{
			for _, cmd := range parserVAL.cmds {
				parserlex.(*aslLexer).AddCommand(cmd)
			}

			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 3:
		//line internal_parser.y:66
		{
			if parserS[parserpt-0].cmd != nil {
				parserVAL.cmds = append(parserVAL.cmds, parserS[parserpt-0].cmd)
			}
		}
	case 4:
		//line internal_parser.y:68
		{
			if parserS[parserpt-0].cmd != nil {
				parserVAL.cmds = []command{parserS[parserpt-0].cmd}
			} else {
				parserVAL.cmds = []command{}
			}
		}
	case 5:
		//line internal_parser.y:70
		{
			parserVAL.cmds = parserS[parserpt-0].cmds
		}
	case 6:
		//line internal_parser.y:74
		{
			parserVAL.cmds = parserS[parserpt-1].cmds
		}
	case 7:
		//line internal_parser.y:78
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 8:
		//line internal_parser.y:80
		{
			parserVAL.cmd = parserS[parserpt-1].cmd
		}
	case 9:
		//line internal_parser.y:82
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 10:
		//line internal_parser.y:84
		{
			parserVAL.cmd = parserS[parserpt-0].cmd
		}
	case 11:
		//line internal_parser.y:86
		{
			parserVAL.cmd = newEvaluateCommand(parserS[parserpt-1].ex)
		}
	case 12:
		//line internal_parser.y:88
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 13:
		//line internal_parser.y:90
		{
			parserVAL.cmd = parserVAL.cmd
		}
	case 14:
		//line internal_parser.y:94
		{
			parserVAL.cmd = newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 15:
		//line internal_parser.y:98
		{
			parserVAL.cmd = newAssignCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex)
		}
	case 16:
		//line internal_parser.y:102
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 17:
		//line internal_parser.y:104
		{
			parserVAL.ex = newArrayExpression(parserS[parserpt-1].exa, parserS[parserpt-2].t.line, parserS[parserpt-2].t.start)
		}
	case 18:
		//line internal_parser.y:106
		{
			parserVAL.ex = newMapExpression(parserS[parserpt-2].exl, parserS[parserpt-3].t.line, parserS[parserpt-3].t.start)
		}
	case 19:
		//line internal_parser.y:108
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 20:
		//line internal_parser.y:110
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 21:
		//line internal_parser.y:112
		{
			parserVAL.ex = newVariableExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 22:
		//line internal_parser.y:114
		{
			parserVAL.ex = newBooleanExpression(true, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 23:
		//line internal_parser.y:116
		{
			parserVAL.ex = newBooleanExpression(false, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 24:
		//line internal_parser.y:118
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 25:
		//line internal_parser.y:120
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 26:
		//line internal_parser.y:122
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 27:
		//line internal_parser.y:126
		{
			parserVAL.exa = append(parserVAL.exa, parserS[parserpt-0].ex)
		}
	case 28:
		//line internal_parser.y:128
		{
			parserVAL.exa = []expression{parserS[parserpt-0].ex}
		}
	case 29:
		//line internal_parser.y:132
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 30:
		//line internal_parser.y:134
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 31:
		//line internal_parser.y:136
		{
			parserVAL.exl = map[string]expression{}
		}
	case 32:
		//line internal_parser.y:140
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 33:
		//line internal_parser.y:145
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 34:
		//line internal_parser.y:147
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 35:
		//line internal_parser.y:149
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 36:
		//line internal_parser.y:151
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 37:
		//line internal_parser.y:153
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 38:
		//line internal_parser.y:155
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 39:
		//line internal_parser.y:157
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 40:
		//line internal_parser.y:159
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 41:
		//line internal_parser.y:161
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 42:
		//line internal_parser.y:163
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 43:
		//line internal_parser.y:165
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 44:
		//line internal_parser.y:167
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 45:
		//line internal_parser.y:169
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 46:
		//line internal_parser.y:171
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 47:
		//line internal_parser.y:173
		{
			parserVAL.ex = newLogicalExpression(parserS[parserpt-0].ex, booleanExpressionZero, parserS[parserpt-1].t, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 48:
		//line internal_parser.y:177
		{
			parserVAL.ex = newFunctionCallExpression(parserS[parserpt-3].ex, parserS[parserpt-1].exl, parserS[parserpt-3].ex.line(), parserS[parserpt-3].ex.position())
		}
	case 49:
		//line internal_parser.y:181
		{
			parserVAL.ex = newPropertyExpression(newGlobalExpression(parserS[parserpt-0].t.line, parserS[parserpt-0].t.start), parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 50:
		//line internal_parser.y:183
		{
			parserVAL.ex = parserS[parserpt-0].ex
		}
	case 51:
		//line internal_parser.y:187
		{
			parserVAL.ex = newPropertyExpression(parserS[parserpt-2].ex, parserS[parserpt-0].t.source, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 52:
		//line internal_parser.y:191
		{
			parserVAL.exl[parserS[parserpt-0].exi.key] = parserS[parserpt-0].exi.value
		}
	case 53:
		//line internal_parser.y:193
		{
			parserVAL.exl = map[string]expression{parserS[parserpt-0].exi.key: parserS[parserpt-0].exi.value}
		}
	case 54:
		//line internal_parser.y:195
		{
			parserVAL.exl = map[string]expression{}
		}
	case 55:
		//line internal_parser.y:199
		{
			parserVAL.exi = parseArgument{parserS[parserpt-2].t.source, parserS[parserpt-0].ex}
		}
	case 56:
		//line internal_parser.y:201
		{
			parserVAL.exi = parseArgument{singleUnnamedArgument, parserS[parserpt-0].ex}
		}
	case 57:
		//line internal_parser.y:205
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-1].ex, parserS[parserpt-0].cmds, []command{})
		}
	case 58:
		//line internal_parser.y:207
		{
			parserVAL.cmd = newIfThenElseCommand(parserS[parserpt-3].ex, parserS[parserpt-2].cmds, parserS[parserpt-0].cmds)
		}
	case 59:
		//line internal_parser.y:211
		{
			parserVAL.cmd = newWhileLoopCommand(parserS[parserpt-1].ex, parserS[parserpt-0].cmds)
		}
	}
	goto parserstack /* stack new state and value */
}
