package parser

const (
	PREC_NONE       = 0
	PREC_ASSIGNMENT = 1  // =
	PREC_OR         = 2  // or
	PREC_AND        = 3  // and
	PREC_EQUALITY   = 4  // == !=
	PREC_COMPARISON = 5  // < > <= >=
	PREC_TERM       = 6  // + -
	PREC_MODE       = 7  // %
	PREC_FACTOR     = 8  // * /
	PREC_POWER      = 9  // **
	PREC_UNARY      = 10 // ! -
	PREC_CALL       = 11 // . ()
	PREC_PRIMARY    = 12 // number, string, id
)
