package model

import "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	// ID uint `gorm:"primary_key" json:"id"`
	// ExpiresAt int64  `json:"exp,omitempty"`
	// ID        string `json:"jti,omitempty"`
	// Userid uint `gorm:"default:null" json:"user_id"`
	// Token      string `gorm:"unique; not null" json:"token"`
	// Type       string `gorm:"not null" json:"type"`
	// Role     string `json:"role"`
	// Username string `json:"username"`
	// Expiredate int64  `gorm:"not null" json:"expiredate"`
	*jwt.StandardClaims
}

type Responsetoken struct {
	Token     string
	Message   string
	Role      interface{} `json:"role"`
	Username  interface{} `json:"username"`
	ExpiresAt interface{} `json:"exp,omitempty"`
	Id        interface{} `json:"jti,omitempty"`
}
