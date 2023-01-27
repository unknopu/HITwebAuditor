package me

import (
	"github.com/golang-jwt/jwt"
)

type Request struct {
	ProductName string `json:"productName"`
	ProductID   string `json:"productid"`
}
type (
	JwtClaims struct {
		CHash          string `json:"c_hash"`
		Email          string `json:"email"`
		EmailVerified  string `json:"email_verified"`
		AuthTime       int    `json:"auth_time"`
		NonceSupported bool   `json:"nonce_supported"`
		jwt.StandardClaims
	}

	JwtHeader struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}

	JwtKeys struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		Use string `json:"use"`
		Alg string `json:"alg"`
		N   string `json:"n"`
		E   string `json:"e"`
	}
)
