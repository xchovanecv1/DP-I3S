package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"
)

// PropertyTransit represents an authenticated application user
type PropertyTransit struct {
	tableName struct{} `pg:",discard_unknown_columns"`

	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`

	CardID     int `json:"card_id"`
	PropertyID int `json:"property_id"`
	RoomID     int `json:"room_id"`
	UserID     int `json:"user_id"`

	PropertyName string `json:"property_name,omitempty" sql:"-"`
	PropertyCode string `json:"property_code,omitempty" sql:"-"`
	UserName     string `json:"user_name,omitempty" sql:"-"`
	UserEmail    string `json:"user_email,omitempty" sql:"-"`
	RoomName     string `json:"room_name,omitempty" sql:"-"`

	Comment  string `json:"comment"`
	PostAuth bool   `json:"post_auth"`

	Transition string `json:"transition"`
}

// BeforeInsert hook executed before database insert operation.
func (c *PropertyTransit) BeforeInsert(db orm.DB) error {
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
func (c *PropertyTransit) Validate() error {
	return validation.ValidateStruct(c)
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
