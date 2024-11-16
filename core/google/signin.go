package google

import (
	"encoding/json"
	"net/http"
)

// Signin sign in model
type Signin struct {
	baseURL string
}

// SigninInterface google signin interface
type SigninInterface interface {
	Get(token string) (*Response, error)
}

// New signin service
func New() *Signin {
	return &Signin{
		baseURL: "https://oauth2.googleapis.com/tokeninfo",
	}
}

// Response response
type Response struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Jti           string `json:"jti"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	Typ           string `json:"typ"`
}

// Get get data from token
func (s *Signin) Get(token string) (*Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", s.baseURL+"?id_token="+token, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
