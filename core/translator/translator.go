package translator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

var (
	// ENTranslator english translator
	ENTranslator ut.Translator
)

// InitTranslator init translator
func InitTranslator() {
	translator := en.New()
	uni := ut.New(translator, translator)

	ENTranslator, _ = uni.GetTranslator("en")
}
