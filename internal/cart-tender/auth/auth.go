package auth

import (
	"errors"
	"time"

	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/model"
	"github.com/KimiaMontazeri/cart-tender/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

type (
	User struct {
		Username string `json:"username"`
	}

	UserClaims struct {
		Username string `json:"username"`
		Iat      int64  `json:"iat"`
		Exp      int64  `json:"exp"`
	}

	Authenticator struct {
		cfg config.JsonWebToken
	}

	Auth interface {
		ValidateUserToken(tkn string) (User, error)
		GenerateUserJWT(user model.User) (string, error)
	}
)

func (c User) String() string {
	return c.Username
}

// Valid checks claims issuer.
func (c UserClaims) Valid() error {
	if c.Exp != 0 && c.Exp < time.Now().Unix() {
		return errors.New("token has been expired")
	}

	return nil
}

func NewAuthenticator(cfg config.JsonWebToken) *Authenticator {
	return &Authenticator{cfg: cfg}
}

// ValidateUserToken given token with given secret and if it is valid it returns client information from its claims.
func (a *Authenticator) ValidateUserToken(tkn string) (User, error) {
	// Validating and parsing the tokenString
	token, err := jwt.ParseWithClaims(tkn, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validating if algorithm used for signing is same as the algorithm in token
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(a.cfg.Secret), nil
	})
	if err != nil {
		return User{}, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return User{}, jwt.ErrInvalidKey
	}

	return User{
		Username: claims.Username,
	}, nil
}

func (a *Authenticator) GenerateUserJWT(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(a.cfg.Expiration).Unix()

	return token.SignedString([]byte(a.cfg.Secret))
}
