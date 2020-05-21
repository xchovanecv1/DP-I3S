package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

// Access represents an authenticated application user
type DoorAccess struct {
	ID int `json:"id"`

	EnterAt time.Time `json:"enter_at"`
	ExitAt  time.Time `json:"exit_at"`

	MAC string `json:"mac"`

	CardCode string `json:"card_code"`

	CardID int   `json:"-"`
	Card   *Card `json:"card"`

	GatewayID int      `json:"-"`
	Gateway   *Gateway `json:"gateway"`
}

// BeforeInsert hook executed before database insert operation.
func (a *DoorAccess) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if a.EnterAt.IsZero() {
		a.EnterAt = now
	}
	return a.Validate()
}

/*
// BeforeUpdate hook executed before database update operation.
func (a *Access) BeforeUpdate(db orm.DB) error {
	a.UpdatedAt = time.Now()
	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *Access) BeforeDelete(db orm.DB) error {
	return nil
}
*/
// Validate validates Access struct and returns validation errors.
func (a *DoorAccess) Validate() error {
	// TODO remove if not needed
	validation.ValidateStruct(a,
		validation.Field(&a.MAC, validation.Required, is.MAC),
	)
	return nil /* */
}

/*
// CanLogin returns true if user is allowed to login.
func (a *Access) CanLogin() bool {
	return a.Active
}
// Claims returns the Access's claims to be signed
func (a *Access) Claims() jwtauth.Claims {
	return jwtauth.Claims{
		"id":    a.ID,
		"sub":   a.Name,
		"roles": a.Roles,
	}
}

*/
