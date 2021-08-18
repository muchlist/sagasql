package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate input
func (p ProductReq) Validate() error {
	if err := validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Price, validation.Min(0)),
	); err != nil {
		return err
	}

	return nil
}
