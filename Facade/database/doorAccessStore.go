package database

import (
	"log"

	"github.com/go-pg/pg"
	"gitlab.com/IIIS/backend/facade/models"
)

// DoorAccessStore implements database operations for account management by user.
type DoorAccessStore struct {
	db *pg.DB
}

// NewDoorAccessStore returns an DoorAccessStore.
func NewDoorAccessStore(db *pg.DB) *DoorAccessStore {
	return &DoorAccessStore{
		db: db,
	}
}

// AddAccess adds new acces requesto to the DB
func (s *DoorAccessStore) AddAccess(a *models.DoorAccess) error {
	err := s.db.Insert(a)
	log.Println(err)
	return err
}

/*
// Get an account by ID.
func (s *DoorAccessStore) Get(id int) (*models.Account, error) {
	a := models.Account{ID: id}
	err := s.db.Model(&a).
		Where("account.id = ?id").
		Column("account.*", "Token").
		First()
	return &a, err
}
// Update an account.
func (s *DoorAccessStore) Update(a *models.Account) error {
	_, err := s.db.Model(a).
		Column("email", "name").
		WherePK().
		Update()
	return err
}

// Delete an account.
func (s *DoorAccessStore) Delete(a *models.Account) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&jwt.Token{}).
			Where("account_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(&models.Profile{}).
			Where("account_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		return tx.Delete(a)
	})
	return err
}

// UpdateToken updates a jwt refresh token.
func (s *DoorAccessStore) UpdateToken(t *jwt.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		WherePK().
		Update()
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *DoorAccessStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t)
	return err
}

*/
