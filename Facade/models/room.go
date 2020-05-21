package models

import (
	"time"

	guuid "github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"
)

// Room represents an authenticated application user
type Room struct {
	ID        int       `json:"id"`
	UUID      string    `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Gateways []*Gateway `json:"gateways"`

	Name string `json:"name"`
	Code string `json:"code"`

	ActiveUserID int      `json:"-"`
	ActiveUser   *Account `json:"user" pg:"fk:active_user_id"`
}

// BeforeInsert hook executed before database insert operation.
func (r *Room) BeforeInsert(db orm.DB) error {
	now := time.Now()

	if len(r.UUID) == 0 {
		id := guuid.New()
		r.UUID = id.String()
	}

	if r.CreatedAt.IsZero() {
		r.UpdatedAt = now
		r.CreatedAt = now
	}
	return r.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (r *Room) BeforeUpdate(db orm.DB) error {
	r.UpdatedAt = time.Now()
	return r.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (r *Room) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (r *Room) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
	)
}
