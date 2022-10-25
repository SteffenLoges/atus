package user

import (
	"atus/backend/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func GetByToken(token string) (*User, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(config.Base.Auth.JWTSecret), nil
	}

	c := jwt.MapClaims{}

	if _, err := jwt.ParseWithClaims(token, &c, keyFunc); err != nil {
		return nil, err
	}

	if err := c.Valid(); err != nil {
		return nil, err
	}

	return &User{
		UID:  c["uid"].(string),
		Name: c["name"].(string),
	}, nil
}

func (u *User) GenToken() (string, error) {

	atClaims := jwt.MapClaims{
		"uid":  u.UID,
		"name": u.Name,
		"exp":  time.Now().Add(time.Minute * config.Base.Auth.JWTTTL).Unix(),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	return at.SignedString([]byte(config.Base.Auth.JWTSecret))
}
