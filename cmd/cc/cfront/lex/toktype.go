package lex

//go:generate go run golang.org/x/tools/cmd/stringer -type=Toktype
type Toktype int

const (
	EOF Toktype = iota
	ERR

	//
	IDENT

	//
	INT_CONST

	// keywords
	AUTO
	BREAK
	CASE
	CHAR
	CONST
	CONTINUE
	DEFAULT
	DO
	DOUBLE
	ELSE
	ENUM
	EXTERN
	FLOAT
	FOR
	GOTO
	IF
	INT
	LONG
	REGISTER
	RETURN
	SHORT
	SIGNED
	SIZEOF
	STATIC
	STRUCT
	SWITCH
	TYPEDEF
	UNION
	UNSIGNED
	VAR
	VOID
	VOLATILE
	WHILE

	//
	LPAREN     // (
	RPAREN     // )
	LBRACE     // {
	RBRACE     // }
	LBRACKET   // [
	RBRACKET   // ]
	LSHIFT     // <<
	RSHIFT     // >>
	GT         // >
	LT         // <
	GEQ        // >=
	LEQ        // <=
	EQ         // ==
	NEQ        // !=
	BAND       // &
	BXOR       // ^
	BOR        // |
	BCOMP      // ~
	AND        // &&
	OR         // ||
	QMARK      // ?
	INC        // ++
	DEC        // --
	DOT        // .
	ARROW      // ->
	MUL        // *
	DIV        // /
	MOD        // %
	ADD        // +
	SUB        // -
	COMMA      // ,
	ASSIGN     // =
	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	DIV_ASSIGN // /=
	MOD_ASSIGN // %=
	RS_ASSIGN  // >>=
	LS_ASSIGN  // <<=
	BA_ASSIGN  // &=
	XO_ASSIGN  // ^=
	BO_ASSIGN  // |=
	SCOLON     // ;
	NOT        // !
	ELLIP      // ...
	COLON      // :
)
