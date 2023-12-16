package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	MaxUsernameLen = 20
	MaxPasswordLen = 30
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r UserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.Length(4, MaxUsernameLen)),
		validation.Field(&r.Password, validation.Required, validation.Length(4, MaxPasswordLen)),
	)
}
