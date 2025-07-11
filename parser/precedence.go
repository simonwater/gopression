package parser

const (
	PREC_NONE       = 0
	PREC_ASSIGNMENT = 1  // =
	PREC_OR         = 2  // or
	PREC_AND        = 3  // and
	PREC_EQUALITY   = 4  // == !=
	PREC_COMPARISON = 5  // < > <= >=
	PREC_TERM       = 6  // + -
	PREC_FACTOR     = 7  // * / %
	PREC_POWER      = 8  // **
	PREC_UNARY      = 9  // ! -
	PREC_CALL       = 10 // . ()
	PREC_PRIMARY    = 11 // number, string, id
)
