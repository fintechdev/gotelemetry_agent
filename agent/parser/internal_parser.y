%{

	package parser
	
%}

%union {
	cmd command
	ex expression
	t token
}

%type <cmd> command, set_property
%type <ex> expr
%token <t> T_STRING T_NUMBER T_IDENTIFIER 
%token <t> T_PLUS T_MINUS T_MULTIPLY T_DIVIDE
%token <t> T_COMMA T_DOT T_COLON
%token <t> T_OPEN_PARENS T_CLOSE_PARENS T_OPEN_BRACKET T_CLOSE_BRACKET 

%left T_PLUS T_MINUS
%left T_MULTIPLY T_DIVIDE
%left T_UMINUS T_UPLUS

%%

commands				: /* empty */
								| commands command

command					:
								set_property
								;

set_property    :
								T_IDENTIFIER T_COLON expr
									{
										parserlex.(*aslLexer).AddCommand(newOutputCommand($1, $3))
									}
								;

expr						: T_OPEN_PARENS expr T_CLOSE_PARENS
										{ $$ = $2 }
								|	expr T_PLUS expr
										{ $$ = newArithmeticExpression($1, $3, $2.terminal, $1.line(), $1.position()) }
								|	expr T_MINUS expr
										{ $$ = newArithmeticExpression($1, $3, $2.terminal, $1.line(), $1.position()) }
								|	expr T_MULTIPLY expr
										{ $$ = newArithmeticExpression($1, $3, $2.terminal, $1.line(), $1.position()) }
								|	expr T_DIVIDE expr
										{ $$ = newArithmeticExpression($1, $3, $2.terminal, $1.line(), $1.position()) }
								|	T_MINUS expr 			%prec T_UMINUS
										{ $$ = newArithmeticExpression(numericExpressionZero, $2, $1.terminal, $1.line, $1.start) }
								|	T_PLUS expr 			%prec T_UPLUS
										{ $$ = newArithmeticExpression(numericExpressionZero, $2, $1.terminal, $1.line, $1.start) }
								| T_NUMBER
										{ $$ = newNumericExpression($1.source, $1.line, $1.start) }
								| T_STRING
										{ $$ = newStringExpression($1.source, $1.line, $1.start) }
								;