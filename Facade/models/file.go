package models

import (
	"strings"
	"time"

	guuid "github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-pg/pg/orm"
)

// File represents an authenticated application user
type File struct {
	ID        int       `json:"-"`
	UUID      string    `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`

	Name string `json:"name"`
	Mime string `json:"mime"`
	Size int64  `json:"size"`

	AccountID   int          `json:"-"`
	Account     *Account     `json:"-" pg:"fk:account_id"`
	AccountData *AccountPure `json:"user" sql:"-" pg:",discard_unknown_columns"`
}

// IsTableFile check weather file is an excel table
func (f *File) IsTableFile() bool {
	if !strings.EqualFold(f.Mime, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") &&
		!strings.EqualFold(f.Mime, "application/vnd.oasis.opendocument.spreadsheet") {
		return false
	}
	return true
}

// BeforeInsert hook executed before database insert operation.
func (r *File) BeforeInsert(db orm.DB) error {
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
func (r *File) BeforeUpdate(db orm.DB) error {
	//r.UpdatedAt = time.Now()
	return r.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (r *File) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (r *File) Validate() error {
	return validation.ValidateStruct(r)
}
