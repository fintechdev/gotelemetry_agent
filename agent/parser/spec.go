package parser

// Parser specifications
//
// A script is made up of one or more commands
//
// Commands can span multiple lines
//
// Commands can take the following form:
//
// ident : exp									// Assign the result of `exp` to the output property `ident`
//
// var : exp										// Assign the result of `exp` to the variable `var`
//
// exp													// evaluate `exp` and discard its output
//
// Entities take the form /[A-Za-z_][a-zA-Z0-9_]*/
//
// Entities that start with a `$` are considered local variables
//
// Expressions can take the form:
//
// val													// A value
//
// ident(expr)									// A function call
//
// expr1.ident									// Execute method `ident` of `expr1`
//
// val1 op val2									// Execute the operation `op` between `val1` and `val2`
//
// Values can take the following forms:
//
// /[+\-.]?[0-9.]+[0-9]*/				// A number (always represented internally as a float64)
//
// /"[\"]*"/										// A string (usual escapes are allowed)
//
// /true|false/									// A Boolean value
//
// /\[expr*\]/									// An array
//
// /\{(ident:expr)*\}/					// A hash of key/value pairs
