package validator

import (
	"auditor/core/translator"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	entranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

// CoreValidator validator
type CoreValidator struct {
	validator *validator.Validate
}

// New new
func New() *CoreValidator {
	v := &CoreValidator{
		validator: validator.New(),
	}
	v.customValidateor()
	v.translator()
	return v
}

func (cv *CoreValidator) customValidateor() {
	_ = cv.validator.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return validPassword(fl.Field().String())
	})
}

func validPassword(s string) bool {
	if len(s) >= 8 {
		return true
	}
	return false
}

func (cv *CoreValidator) translator() {
	if err := entranslations.RegisterDefaultTranslations(cv.validator, translator.ENTranslator); err != nil {
		panic(err)
	}
	_ = cv.validator.RegisterTranslation("required", translator.ENTranslator,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} is a required field", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		},
	)
}

// Validate validator
func (cv *CoreValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Var var
func (cv *CoreValidator) Var(field interface{}, tag string) error {
	return cv.validator.Var(field, tag)
}
