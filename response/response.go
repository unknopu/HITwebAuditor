package response

import (
	"auditor/core/mongodb"
	"auditor/entities"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Error error
type Error struct {
	Code    int    `json:"code,omitempty" mapstructure:"code"`
	Message string `json:"message,omitempty" mapstructure:"message"`
}

// Results return results
type Results struct {
	DuplicateError             Error `mapstructure:"duplicate_error"`
	InvalidParameters          Error `mapstructure:"invalid_parameters"`
	InvalidID                  Error `mapstructure:"invalid_id"`
	DataNotFound               Error `mapstructure:"data_notfound"`
	InvalidLoginType           Error `mapstructure:"invalid_login_type"`
	InvalidEmail               Error `mapstructure:"invalid_email"`
	InvalidPassword            Error `mapstructure:"invalid_password"`
	InvalidFacebookToken       Error `mapstructure:"invalid_facebook_token"`
	InvalidGoogleToken         Error `mapstructure:"invalid_google_token"`
	BootsServiceUnavailable    Error `mapstructure:"boots_service_unavailable"`
	LoginMemberNotFound        Error `mapstructure:"member_not_found"`
	LoginMemberAlreadyExisting Error `mapstructure:"member_already_existing"`
	RegisterBoots              Error `mapstructure:"register_boots"`
	ErrInvalidCitizen          Error `mapstructure:"invalid_citizen"`
	ErrFlashSaleExceed         Error `mapstructure:"boots_flash_sales_exceed"`
}

func (ec Error) Error() string {
	return ec.Message
}

// ErrorCode get error code
func (ec Error) ErrorCode() int {
	return ec.Code
}

// GetResponse get error response
func (r *Results) GetResponse(err error, l *i18n.Localizer) error {
	if _, ok := err.(*echo.HTTPError); ok {
		echoE := err.(*echo.HTTPError)
		echoE.Message = r.getErrorLocalizeMessage(echoE.Message.(string), l)
		return echoE
	} else if _, ok := err.(Error); ok {
		echoE := err.(Error)
		echoE.Message = r.getErrorLocalizeMessage(echoE.Message, l)
		return err
	}
	switch true {
	case err == mongodb.ErrorNotFound:
		return r.DataNotFound
	case err == mongodb.ErrorInvalidID:
		return r.InvalidID
	case err == mongodb.ErrorDucumentDuplicate:
		return r.DuplicateError
	case err == bcrypt.ErrMismatchedHashAndPassword:
		return r.InvalidPassword
	case err == entities.ErrBootsServiceUnavailable:
		return r.BootsServiceUnavailable
	case err == entities.ErrLoginMemberNotFound:
		return r.LoginMemberNotFound
	case err == entities.ErrLoginMemberAlreadyExisting:
		return r.LoginMemberAlreadyExisting
	case err == entities.ErrRegisterBoots:
		return r.RegisterBoots
	case err == entities.ErrInvalidCitizen:
		return r.ErrInvalidCitizen
	default:
		return Error{
			Code:    0,
			Message: r.getErrorLocalizeMessage(err.Error(), l),
		}
	}
}

// ReadReturnResult read response
func ReadReturnResult(path, filename string) (*Results, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigType("yml")
	v.SetConfigName(filename)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	r := &Results{}
	if err := v.Unmarshal(r); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Results) getErrorLocalizeMessage(message string, l *i18n.Localizer) string {
	m, err := l.LocalizeMessage(&i18n.Message{
		ID: message,
	})
	if err == nil {
		return m
	}
	return message
}
