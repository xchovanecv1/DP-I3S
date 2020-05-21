package models

import (
	"strings"
	"time"

	guuid "github.com/google/uuid"

	"github.com/go-chi/jwtauth"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-pg/pg/orm"
	"gitlab.com/IIIS/backend/facade/auth/jwt"
)

// Account represents an authenticated application user
type Account struct {
	ID        int       `json:"-"`
	UUID      string    `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	LastLogin time.Time `json:"last_login,omitempty"`

	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Pass   string   `json:"pass"`
	Active bool     `sql:",notnull" json:"active"`
	Roles  []string `pg:",array" json:"roles,omitempty"`

	Card   *Card `pg:"fk:card_id"`
	CardID int   `json:"card_id,omitempty"`

	Token []jwt.Token `sql:"-" json:"token,omitempty"`
}

// AccountPure represents an authenticated application user
type AccountPure struct {
	ID        int       `json:"-"`
	UUID      string    `json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	LastLogin time.Time `json:"-"`

	Email  string   `json:"-"`
	Name   string   `json:"name"`
	Pass   string   `json:"-"`
	Active bool     `json:"-"`
	Roles  []string `pg:",array" json:"-"`

	Card   *Card `pg:"fk:card_id" json:"-"`
	CardID int   `json:"card_id,omitempty" json:"-"`

	Token []jwt.Token `sql:"-" json:"-,omitempty"`
}

/*
func (a *Account) Strip(db orm.DB) *AccountPure {
	ret := &AccountPure{
		UUID: a.UUID,
	}

	return ret
}*/

// BeforeInsert hook executed before database insert operation.
func (a *Account) BeforeInsert(db orm.DB) error {
	now := time.Now()

	if len(a.UUID) == 0 {
		id := guuid.New()
		a.UUID = id.String()
	}

	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}
	return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (a *Account) BeforeUpdate(db orm.DB) error {
	a.UpdatedAt = time.Now()
	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *Account) BeforeDelete(db orm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (a *Account) Validate() error {
	a.Email = strings.TrimSpace(a.Email)
	a.Email = strings.ToLower(a.Email)
	a.Name = strings.TrimSpace(a.Name)

	return validation.ValidateStruct(a,
		validation.Field(&a.Email, validation.Required, is.Email, is.LowerCase),
		validation.Field(&a.Name, validation.Required, is.ASCII),
	)
}

// CanLogin returns true if user is allowed to login.
func (a *Account) CanLogin() bool {
	return a.Active
}

// Claims returns the account's claims to be signed
func (a *Account) Claims() jwtauth.Claims {
	return jwtauth.Claims{
		"id":    a.UUID,
		"sub":   a.Name,
		"roles": a.Roles,
	}
}
