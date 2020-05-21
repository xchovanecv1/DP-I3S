package database

import (
	"github.com/go-pg/pg"
	"gitlab.com/IIIS/backend/facade/auth/jwt"
	"gitlab.com/IIIS/backend/facade/models"
)

// AccountStore implements database operations for account management by user.
type AccountStore struct {
	db *pg.DB
}

// NewAccountStore returns an AccountStore.
func NewAccountStore(db *pg.DB) *AccountStore {
	return &AccountStore{
		db: db,
	}
}

// Get an account by ID.
func (s *AccountStore) Get(id int) (*models.Account, error) {
	a := models.Account{ID: id}
	err := s.db.Model(&a).
		Where("account.id = ?id").
		Column("account.*").
		First()
	return &a, err
}

// GetByUUID an account by UUID.
func (s *AccountStore) GetByUUID(uuid string) (*models.Account, error) {
	a := models.Account{UUID: uuid}
	err := s.db.Model(&a).
		Where("account.uuid = ?uuid").
		Column("account.*").
		First()
	return &a, err
}

// GetByRole an account by Rome.
func (s *AccountStore) GetByRole(rl string) (*models.Account, error) {
	var role [1]string
	role[0] = rl
	a := models.Account{Roles: role[:]}
	err := s.db.Model(&a).
		Where("account.roles @> ?roles").
		Column("account.*").
		First()
	return &a, err
}

// GetAccountByCard returns an account by email.
func (s *AccountStore) GetAccountByCard(cid int) (*models.Account, error) {
	a := models.Account{CardID: cid}
	err := s.db.Model(&a).
		Column("account.*").
		Where("card_id = ?card_id").
		First()
	return &a, err
}

// Update an account.
func (s *AccountStore) Update(a *models.Account) error {
	_, err := s.db.Model(a).
		Column("email", "name").
		WherePK().
		Update()
	return err
}

// Delete an account.
func (s *AccountStore) Delete(a *models.Account) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&jwt.Token{}).
			Where("account_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		return tx.Delete(a)
	})
	return err
}

// UpdateToken updates a jwt refresh token.
func (s *AccountStore) UpdateToken(t *jwt.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		WherePK().
		Update()
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *AccountStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t)
	return err
}
