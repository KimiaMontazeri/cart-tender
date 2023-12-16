package request

import (
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	CartCreateRequest struct {
		Data string `json:"data"`
	}

	CartUpdateRequest struct {
		ID    int64  `json:"id"`
		Data  string `json:"data"`
		State string `json:"state"`
	}

	CartDeleteRequest struct {
		ID int64 `json:"id"`
	}
)

func (r CartCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Data, validation.Required),
	)
}

func (r CartUpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Data, validation.Required),
		validation.Field(&r.State, validation.Required, validation.In(model.PENDING, model.COMPLETED)),
	)
}

func (r CartDeleteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
	)
}
