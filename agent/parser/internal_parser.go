//line internal_parser.y:2
package parser

import __yyfmt__ "fmt"

//line internal_parser.y:3
//line internal_parser.y:7
type parserSymType struct {
	yys int
	cmd command
	ex  expression
	t   token
}

const T_STRING = 57346
const T_NUMBER = 57347
const T_IDENTIFIER = 57348
const T_PLUS = 57349
const T_MINUS = 57350
const T_MULTIPLY = 57351
const T_DIVIDE = 57352
const T_COMMA = 57353
const T_DOT = 57354
const T_COLON = 57355
const T_OPEN_PARENS = 57356
const T_CLOSE_PARENS = 57357
const T_OPEN_BRACKET = 57358
const T_CLOSE_BRACKET = 57359
const T_UMINUS = 57360
const T_UPLUS = 57361

var parserToknames = []string{
	"T_STRING",
	"T_NUMBER",
	"T_IDENTIFIER",
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
}

const parserNprod = 14
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 30

var parserAct = []int{

	6, 5, 1, 12, 13, 14, 15, 4, 16, 17,
	18, 23, 3, 19, 20, 21, 22, 11, 10, 2,
	9, 8, 12, 13, 14, 15, 0, 7, 14, 15,
}
var parserPact = []int{

	-1000, 1, -1000, -1000, -12, 13, 15, 13, 13, 13,
	-1000, -1000, 13, 13, 13, 13, -4, -1000, -1000, 19,
	19, -1000, -1000, -1000,
}
var parserPgo = []int{

	0, 19, 12, 0, 2,
}
var parserR1 = []int{

	0, 4, 4, 1, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3,
}
var parserR2 = []int{

	0, 0, 2, 1, 3, 3, 3, 3, 3, 3,
	2, 2, 1, 1,
}
var parserChk = []int{

	-1000, -4, -1, -2, 6, 13, -3, 14, 8, 7,
	5, 4, 7, 8, 9, 10, -3, -3, -3, -3,
	-3, -3, -3, 15,
}
var parserDef = []int{

	1, -2, 2, 3, 0, 0, 4, 0, 0, 0,
	12, 13, 0, 0, 0, 0, 0, 10, 11, 6,
	7, 8, 9, 5,
}
var parserTok1 = []int{

	1,
}
var parserTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19,
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
		//line internal_parser.y:35
		{
			parserlex.(*aslLexer).AddCommand(newOutputCommand(parserS[parserpt-2].t, parserS[parserpt-0].ex))
		}
	case 5:
		//line internal_parser.y:41
		{
			parserVAL.ex = parserS[parserpt-1].ex
		}
	case 6:
		//line internal_parser.y:43
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t.terminal, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 7:
		//line internal_parser.y:45
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t.terminal, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 8:
		//line internal_parser.y:47
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t.terminal, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 9:
		//line internal_parser.y:49
		{
			parserVAL.ex = newArithmeticExpression(parserS[parserpt-2].ex, parserS[parserpt-0].ex, parserS[parserpt-1].t.terminal, parserS[parserpt-2].ex.line(), parserS[parserpt-2].ex.position())
		}
	case 10:
		//line internal_parser.y:51
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t.terminal, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 11:
		//line internal_parser.y:53
		{
			parserVAL.ex = newArithmeticExpression(numericExpressionZero, parserS[parserpt-0].ex, parserS[parserpt-1].t.terminal, parserS[parserpt-1].t.line, parserS[parserpt-1].t.start)
		}
	case 12:
		//line internal_parser.y:55
		{
			parserVAL.ex = newNumericExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	case 13:
		//line internal_parser.y:57
		{
			parserVAL.ex = newStringExpression(parserS[parserpt-0].t.source, parserS[parserpt-0].t.line, parserS[parserpt-0].t.start)
		}
	}
	goto parserstack /* stack new state and value */
}
