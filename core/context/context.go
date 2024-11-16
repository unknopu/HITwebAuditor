package context

import (
	"auditor/core/fileutil"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	uuid "github.com/satori/go.uuid"
)

const pathKey = "path"

// Context custom echo context
type Context struct {
	echo.Context
	Parameters interface{}
	SessionID  *string
	Mobile     string
	Locale     string
	Platform   string
	Localizer  *i18n.Localizer
}

// BindAndValidate bind and validate form
func (c *Context) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	c.parsePathParams(i)
	if err := c.Validate(i); err != nil {
		return err
	}
	c.Parameters = i
	return nil
}

// IsWeb is web
func (c *Context) IsWeb() bool {
	return c.Platform == "web"
}

func (c *Context) parsePathParams(form interface{}) {
	formValue := reflect.ValueOf(form)
	if formValue.Kind() == reflect.Ptr {
		formValue = formValue.Elem()
	}
	t := reflect.TypeOf(formValue.Interface())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(pathKey)
		if tag != "" {
			fieldName := t.Field(i).Name
			paramValue := formValue.FieldByName(fieldName)
			if paramValue.IsValid() {
				paramValue.Set(reflect.ValueOf(c.Param(tag)))
			}
		}
	}
}

// Claims jwt claims
type Claims struct {
	jwt.StandardClaims
	Role         int    `json:"role,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	SessionID    string `json:"session_id,omitempty"`
}

// GetUserSession get user session
func (c *Context) GetUserSession() *UserContext {
	token := c.Get("user")
	uc := &UserContext{
		ClientIP:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		Platform:  c.Request().Header.Get("platform"),
	}
	if token != nil {
		user := token.(*jwt.Token)
		cc := user.Claims.(*Claims)
		if cc != nil {
			uc.UserID = cc.Subject
			uc.Role = cc.Role
			uc.SessionID = &cc.SessionID
			uc.MobileNumber = cc.MobileNumber
		}
		return uc
	}
	if c.SessionID != nil {
		uc.SessionID = c.SessionID
		uc.GuestMobileNumber = c.Mobile
	}
	return uc
}

// IsWeb is web
func (c *UserContext) IsWeb() bool {
	return c.Platform == "web"
}

// ToUserSession convert claims to user session
func (c *Claims) ToUserSession() *UserContext {
	return &UserContext{
		UserID: c.Subject,
		Role:   c.Role,
	}
}

// UserContext user context
type UserContext struct {
	ClientIP          string
	UserAgent         string
	UserID            string
	Role              int
	SessionID         *string
	Platform          string
	MobileNumber      string `json:"mobile_number,omitempty"`
	GuestMobileNumber string `json:"guest_mobile_number,omitempty"`
}

// GetThaiPhoneNumber get thai phone number
func (i *UserContext) GetThaiPhoneNumber() string {
	return strings.ReplaceAll(i.MobileNumber, "+66", "0")
}

// GetThaiPhoneNumber get thai phone number
func (i *UserContext) GetGuestThaiPhoneNumber() string {
	return strings.ReplaceAll(i.GuestMobileNumber, "+66", "0")
}

// StreamResponse custom response
func (c *Context) StreamResponse(code int, fi *fileutil.File) error {
	uu := uuid.NewV4().String()
	f, err := os.Open(fi.Path())
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	c.Response().
		Header().
		Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s%s", uu, fi.Ext()))
	return c.Stream(code, fi.ContentType(), f)
}
