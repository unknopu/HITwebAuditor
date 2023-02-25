package sql

import "auditor/handlers/common"

type SQLiBased int

const (
	UnkownBased SQLiBased = iota
	LengthValidation
	ErrorSQLiBased
	BlindSQLiBased
	BetweenSQLiBased
)

type BaseForm struct {
	common.PageQuery
	URL   string  `json:"url"`
	Param string  `json:"param"`
	Level *string `json:"level"`
}
