package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// "github.com/jinzhu/gorm"

//User struct declaration
type User struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Profile   string `json:"profile"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Msg       string `json:"msg"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{})
	return db
}

type Tokens struct {
	Role       string `json:"role,omitempty"`
	Username   string `json:"username,omitempty"`
	ExpiresAt  int64  `json:"exp,omitempty"`
	Id         string `json:"jti,omitempty"`
	Is_revoked bool   `json:"is_revoked"`
	jwt.StandardClaims
	// jwt.StandardClaims
}

type Changepassword struct {
	Currentpassword string `json:"currentpassword"`
	Newpassword     string `json:"newpassword"`
	Retypepassword  string `json:"retypepassword"`
}

type Confirmpassword struct {
	Password       string `json:"password"`
	Repeatpassword string `json:"repeatpassword"`
}
