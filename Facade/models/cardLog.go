package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"
)

// CardLog represents an authenticated application user
type CardLog struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`

	Code      string `json:"code"`
	Type      string `json:"type"`
	GatewayID int    `json:"gateway_id"`
}

// BeforeInsert hook executed before database insert operation.
func (c *CardLog) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}

	return c.Validate()
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
func (c *CardLog) Validate() error {

	return validation.ValidateStruct(c,
		validation.Field(&c.Code, validation.Required),
		validation.Field(&c.Type, validation.Required),
	)
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
