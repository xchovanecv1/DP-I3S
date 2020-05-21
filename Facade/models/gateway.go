package models

import (
	"time"

	guuid "github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
)

// Gateway represents an authenticated application user
type Gateway struct {
	ID         int       `json:"id"`
	UUID       string    `json:"-"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	LastActive time.Time `json:"last_active,omitempty"`

	Name string `json:"name"`
	Addr string `json:"addr"`

	Room   *Room `json:"room"`
	RoomID int   `json:"room_id"`

	Active bool `json:"active"`
}

// GatewayForList represents an gateway for fk list
type GatewayForList struct {
	ID        int       `json:"id"`
	UUID      string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Name string `json:"name"`
	Addr string `json:"-"`

	Room   *Room `json:"-"`
	RoomID int   `json:"-"`

	Active bool `json:"active"`
}

// BeforeInsert hook executed before database insert operation.
func (r *Gateway) BeforeInsert(db orm.DB) error {
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
func (r *Gateway) BeforeUpdate(db orm.DB) error {
	r.UpdatedAt = time.Now()
	return r.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (r *Gateway) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (r *Gateway) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Addr, validation.Required, is.MAC),
	)
}
