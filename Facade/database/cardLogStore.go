package database

import (
	"log"
	"net/url"
	"sort"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
	"gitlab.com/IIIS/backend/facade/models"
)

// CardLogStore implements database operations for account management by user.
type CardLogStore struct {
	db *pg.DB
}

// NewCardLogStore returns an DoorAccessStore.
func NewCardLogStore(db *pg.DB) *CardLogStore {
	return &CardLogStore{
		db: db,
	}
}

// AddLog adds new acces requesto to the DB
func (s *CardLogStore) AddLog(a *models.CardLog) error {
	log.Printf("Data to add: %v\n", a)
	err := s.db.Insert(a)
	log.Println(err)
	return err
}

// GetByCode retrieves Card object by its code property
func (s *CardLogStore) GetByCode(code string) (*models.CardLog, error) {
	a := models.CardLog{Code: code}
	err := s.db.Model(&a).
		Where("card_log.code = ?code").
		Column("card_log.*").
		First()
	return &a, err
}

// Create creates a new card.
func (s *CardLogStore) Create(a *models.CardLog) error {

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(a)

	})

	return err
}

// Delete an card.
func (s *CardLogStore) Delete(a *models.CardLog) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// DeleteByCode an card.
func (s *CardLogStore) DeleteByCode(cl *models.CardLog) error {
	_, err := s.db.Model(cl).Where("code = ?code").Delete()
	return err
}

// Get an account by ID.
func (s *CardLogStore) Get(id int) (*models.CardLog, error) {
	a := models.CardLog{ID: id}
	err := s.db.Model(&a).
		Where("card_log.id = ?id").
		Column("card_log.*").
		First()
	return &a, err
}

// ApplySearch applies an RoomSearch on an orm.Query.
func (f *CardLogFilter) ApplySearch(q *orm.Query) (*orm.Query, error) {

	if len(f.Query) > 0 {
		q = q.Where("LOWER(unaccent(code)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%") // ORDER  BY col_a ILIKE '%keyword%' OR NULL
	}

	return q, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *CardLogFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.ApplySearch)
	//q = q.Where("name LIKE ?", "%"+f.Query[0]+"%")
	q = q.Order("created_at DESC")
	//q = q.Order(f.Order...)
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *CardLogStore) List(f *CardLogFilter) ([]models.CardLog, error) {
	a := []models.CardLog{}

	var err error
	if len(f.Query) > 0 {
		//q = q.Where(, "%"+f.Query[0]+"%") // ORDER  BY col_a ILIKE '%keyword%' OR NULL
		_, err = s.db.Query(&a, "SELECT DISTINCT ON (code) * FROM card_logs WHERE LOWER(unaccent(code)) LIKE LOWER(unaccent(?)) ORDER BY code ASC, created_at DESC", "%"+f.Query[0]+"%")
	} else {
		_, err = s.db.Query(&a, "SELECT DISTINCT ON (code) * FROM card_logs ORDER BY code ASC, created_at DESC")
	}

	if err != nil {
		return nil, err
	}

	sort.Slice(a, func(i, j int) bool {
		return a[i].CreatedAt.After(a[j].CreatedAt)
	})

	return a, nil
}

// CardLogFilter provides pagination and filtering options on accounts.
type CardLogFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
	Query  []string
}

// NewCardLogFilter returns an CardFilter with options parsed from request url values.
func NewCardLogFilter(params interface{}) (*CardLogFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &CardLogFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
		Query:  p["q"],
	}
	return f, nil
}
