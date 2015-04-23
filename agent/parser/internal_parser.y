%{

	package parser

	type parseArgument struct {
		key string
		value expression
	}

	const singleUnnamedArgument = "----"
	
%}

%union {
	cmds []command
	cmd command
	ex expression
	exl map[string]expression
	exi parseArgument
	t token
}

%type <cmds> commands, command_list, command_block
%type <cmd> command, set_property, assign_to_var, if_then_else
%type <ex> expr, function_call, callable_expr, property, operation
%type <exl> expr_list
%type <exi> expr_item
%token <t> T_STRING T_NUMBER T_IDENTIFIER T_VARIABLE T_TRUE T_FALSE
%token <t> T_PLUS T_MINUS T_MULTIPLY T_DIVIDE
%token <t> T_COMMA T_DOT T_COLON
%token <t> T_OPEN_PARENS T_CLOSE_PARENS T_OPEN_BRACKET T_CLOSE_BRACKET T_OPEN_BRACE T_CLOSE_BRACE
%token <t> T_OR T_AND T_EQUAL T_NOT_EQUAL T_NEGATE T_GREATER_THAN T_LESS_THAN T_GREATER_THAN_OR_EQUAL T_LESS_THAN_OR_EQUAL
%token <t> T_IF T_ELSE
%token <t> T_TERMINATOR

%left T_FUNCTION_CALL
%left T_OR
%left T_AND
%left T_GREATER_THAN T_LESS_THAN T_GREATER_THAN_OR_EQUAL T_LESS_THAN_OR_EQUAL
%nonassoc T_EQUAL T_NOT_EQUAL
%left T_PLUS T_MINUS
%left T_MULTIPLY T_DIVIDE
%right T_NEGATE 
%right T_UMINUS T_UPLUS
%left T_DOT

%%

command_list		: 
									{ $$ = []command{} }
								|	commands
									{ 
										for _, cmd := range $$ {
											parserlex.(*aslLexer).AddCommand(cmd)
										}

										$$ = $1
									}
								;

commands				: commands command 
									{ $$ = append($$, $2) }
								| command
									{ $$ = []command{$1} }
								| command_block
									{ $$ = $1 }
								;

command_block		: T_OPEN_BRACE commands T_CLOSE_BRACE
									{ $$ = $2 }


command					: set_property T_TERMINATOR
									{ $$ = $1 }
								| assign_to_var T_TERMINATOR
									{ $$ = $1 }
								| if_then_else
									{ $$ = $1 }
								| expr T_TERMINATOR
									{ $$ = newEvaluateCommand($1) }
								| T_TERMINATOR
									{ $$ = $$ }
								;

set_property    : T_IDENTIFIER T_COLON expr
									{ $$ = newOutputCommand($1, $3) }
								;

assign_to_var		: T_VARIABLE T_COLON expr
									{	$$ = newAssignCommand($1, $3) }
								;

expr						: T_OPEN_PARENS expr T_CLOSE_PARENS
										{ $$ = $2 }
								| T_NUMBER
										{ $$ = newNumericExpression($1.source, $1.line, $1.start) }
								| T_STRING
										{ $$ = newStringExpression($1.source, $1.line, $1.start) }
								| T_VARIABLE
										{ $$ = newVariableExpression($1.source, $1.line, $1.start) }
								|	T_TRUE
										{ $$ = newBooleanExpression(true, $1.line, $1.start) }
								|	T_FALSE
										{ $$ = newBooleanExpression(false, $1.line, $1.start) }
								| operation
										{ $$ = $1 }
								| function_call
										{ $$ = $1 }
								| property
										{ $$ = $1 }
								; 

operation 		  : expr T_PLUS expr
										{ $$ = newArithmeticExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_MINUS expr
										{ $$ = newArithmeticExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_MULTIPLY expr
										{ $$ = newArithmeticExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_DIVIDE expr
										{ $$ = newArithmeticExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_EQUAL expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_NOT_EQUAL expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_AND expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_OR expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_GREATER_THAN expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_GREATER_THAN_OR_EQUAL expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_LESS_THAN expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								| expr T_LESS_THAN_OR_EQUAL expr
										{ $$ = newLogicalExpression($1, $3, $2, $1.line(), $1.position()) }
								|	T_MINUS expr 			%prec T_UMINUS
										{ $$ = newArithmeticExpression(numericExpressionZero, $2, $1, $1.line, $1.start) }
								|	T_PLUS expr 			%prec T_UPLUS
										{ $$ = newArithmeticExpression(numericExpressionZero, $2, $1, $1.line, $1.start) }
								|	T_NEGATE expr 			%prec T_UPLUS
										{ $$ = newLogicalExpression($2, booleanExpressionZero, $1, $1.line, $1.start) }
								;

function_call 	: callable_expr T_OPEN_PARENS expr_list T_CLOSE_PARENS		%prec T_FUNCTION_CALL
										{ $$ = newFunctionCallExpression($1, $3, $1.line(), $1.position()) }
								;

callable_expr		: T_IDENTIFIER
										{ $$ = newPropertyExpression(newGlobalExpression($1.line, $1.start), $1.source, $1.line, $1.start) }
								| property
										{ $$ = $1 }
								;

property 				: expr T_DOT T_IDENTIFIER
										{ $$ = newPropertyExpression($1, $3.source, $1.line(), $1.position()) }
								;

expr_list				: expr_list T_COMMA expr_item
										{ $$[$3.key] = $3.value }
								| expr_item
										{ $$ = map[string]expression{$1.key: $1.value }}
								| /* empty */
										{ $$ = map[string]expression{} }
								;

expr_item				: T_IDENTIFIER T_COLON expr
										{ $$ = parseArgument{$1.source, $3} }
								| expr
										{ $$ = parseArgument{singleUnnamedArgument, $1} }
								;

if_then_else		: T_IF expr command_block
										{ $$ = newIfThenElseCommand($2, $3, []command{}) }
								| T_IF expr command_block T_ELSE command_block
										{ $$ = newIfThenElseCommand($2, $3, $5) }
								;