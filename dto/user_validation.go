package dto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/muchlist/sagasql/config"
	"github.com/muchlist/sagasql/utils/sfunc"
)

// Validate input
func (u UserRegisterReq) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Roles, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	); err != nil {
		return err
	}

	// validate role
	if !sfunc.InSlice(u.Roles, config.GetRolesAvailable()) {
		return fmt.Errorf("role yang dimasukkan tidak tersedia. gunakan %s", config.GetRolesAvailable())
	}

	return nil
}
