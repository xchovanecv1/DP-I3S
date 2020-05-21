package models

import (
	"time"

	guuid "github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"
)

// Property represents an physical property being tracked
type Property struct {
	ID        int       `json:"-"`
	UUID      string    `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Name       string    `json:"name" import:"true"`
	Code       string    `json:"code" import:"true"`
	Price      float32   `json:"price" import:"true"`
	AcquiredAt time.Time `json:"acquired_at,omitempty" import:"true"`

	RoomID int   `json:"room_id" import:"true"`
	Room   *Room `json:"room,omitempty" pg:"fk:room_id"`
}

// BeforeInsert hook executed before database insert operation.
func (r *Property) BeforeInsert(db orm.DB) error {
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
func (r *Property) BeforeUpdate(db orm.DB) error {
	r.UpdatedAt = time.Now()
	return r.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (r *Property) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (r *Property) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("Názov zariadenia je povinné pole")),
		validation.Field(&r.Code, validation.Required.Error("Kódové označenie zariadenia je povinné pole")),
	)
}
