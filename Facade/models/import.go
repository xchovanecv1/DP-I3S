package models

import (
	"time"

	guuid "github.com/google/uuid"

	"github.com/go-pg/pg/orm"
)

// Import represents an physical property being tracked
type Import struct {
	ID          int       `json:"-"`
	UUID        string    `json:"id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	CreatedByID int       `json:"created_by"`

	Resource string `json:"resource"`
	Elements []int  `json:"elements" pg:",array"`
}

// BeforeInsert hook executed before database insert operation.
func (r *Import) BeforeInsert(db orm.DB) error {
	now := time.Now()

	if len(r.UUID) == 0 {
		id := guuid.New()
		r.UUID = id.String()
	}

	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	return r.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (r *Import) BeforeUpdate(db orm.DB) error {
	//r.UpdatedAt = time.Now()
	return r.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (r *Import) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (r *Import) Validate() error {
	return nil
}
