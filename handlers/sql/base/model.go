package base

type SQLi int

const (
	UnkownBased SQLi = iota
	LengthValidation
	ErrorSQLiBased
	BlindSQLiBased
	BetweenSQLiBased
)
