package database

import (
	"log"
	"net/url"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
	"gitlab.com/IIIS/backend/facade/models"
)

// CardStore implements database operations for account management by user.
type CardStore struct {
	db *pg.DB
}

// NewCardStore returns an DoorAccessStore.
func NewCardStore(db *pg.DB) *CardStore {
	return &CardStore{
		db: db,
	}
}

// AddCard adds new acces requesto to the DB
func (s *CardStore) AddCard(a *models.Card) error {
	log.Printf("Data to add: %v\n", a)
	err := s.db.Insert(a)
	log.Println(err)
	return err
}

// GetByCode retrieves Card object by its code property
func (s *CardStore) GetByCode(code string) (*models.Card, error) {
	a := models.Card{Code: code}
	err := s.db.Model(&a).
		Where("card.code = ?code").
		Column("card.*").
		First()
	return &a, err
}

// Create creates a new card.
func (s *CardStore) Create(a *models.Card) error {
	/*
		count, _ := s.db.Model(a).
			Where("name = ?name").
			Count()

		if count != 0 {
			return ErrUniqueEmailConstraint
		}
	*/
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(a)
		/*
			err := tx.Insert(a)
			if err != nil {
				return err
			}
			p := &models.Profile{
				AccountID: a.ID,
			}
			return tx.Insert(p)*/
	})

	return err
}

// CreateFromLog creates a new card from card logs.
func (s *CardStore) CreateFromLog(c *models.Card, l *models.CardLog) error {
	/*
		count, _ := s.db.Model(a).
			Where("name = ?name").
			Count()

		if count != 0 {
			return ErrUniqueEmailConstraint
		}
	*/
	// Copy useful data
	c.Code = l.Code
	c.Type = l.Type

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(c)
	})

	return err
}

// Delete an card.
func (s *CardStore) Delete(a *models.Card) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// Get an account by ID.
func (s *CardStore) Get(id int) (*models.Card, error) {
	a := models.Card{ID: id}
	err := s.db.Model(&a).
		Where("card.id = ?id").
		Column("card.*").
		First()
	return &a, err
}

// GetByUUID an account by UUID.
func (s *CardStore) GetByUUID(uuid string) (*models.Card, error) {
	a := models.Card{UUID: uuid}
	err := s.db.Model(&a).
		Where("card.uuid = ?uuid").
		Column("card.*").
		First()
	return &a, err
}

// ApplySearch applies an RoomSearch on an orm.Query.
func (f *CardFilter) ApplySearch(q *orm.Query) (*orm.Query, error) {

	if len(f.Query) > 0 {
		q = q.Where("LOWER(unaccent(code)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%") // ORDER  BY col_a ILIKE '%keyword%' OR NULL
	}

	return q, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *CardFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.ApplySearch)
	//q = q.Where("name LIKE ?", "%"+f.Query[0]+"%")
	q = q.Order(f.Order...)
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *CardStore) List(f *CardFilter) ([]models.Card, int, error) {
	a := []models.Card{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// CardFilter provides pagination and filtering options on accounts.
type CardFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
	Query  []string
}

// NewCardFilter returns an CardFilter with options parsed from request url values.
func NewCardFilter(params interface{}) (*CardFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &CardFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
		Query:  p["q"],
	}
	return f, nil
}

// ListLink applies a filter and returns paginated array of matching results and total count.
func (s *CardStore) ListLink(linkType string, linkID int) ([]models.Card, int, error) {
	a := []models.Card{}
	count, err := s.db.Model(&a).
		Where("link = ?", linkType).
		Where("link_id = ?", linkID).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// Update an card.
func (s *CardStore) Update(a *models.Card) error {
	_, err := s.db.Model(a).
		Column("link").
		Column("link_id").
		Where("uuid = ?uuid").
		Update()
	return err
}
