package database

import (
	"log"
	"net/url"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
	"gitlab.com/IIIS/backend/facade/models"
)

// ImportStore implements database operations for account management by user.
type ImportStore struct {
	db *pg.DB
}

// NewImportStore returns an DoorAccessStore.
func NewImportStore(db *pg.DB) *ImportStore {
	return &ImportStore{
		db: db,
	}
}

// AddImport adds new acces requesto to the DB
func (s *ImportStore) AddImport(a *models.Import) error {
	log.Printf("Data to add: %v\n", a)
	err := s.db.Insert(a)
	log.Println(err)
	return err
}

// Create creates a new import.
func (s *ImportStore) Create(a *models.Import) error {
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

// CreateMultiple creates a new properties.
func (s *ImportStore) CreateMultiple(a []*models.Import) error {

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		for _, s := range a {
			err := tx.Insert(s)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// Delete an import.
func (s *ImportStore) Delete(a *models.Import) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// Get an account by ID.
func (s *ImportStore) Get(id int) (*models.Import, error) {
	a := models.Import{ID: id}
	err := s.db.Model(&a).
		Where("import.id = ?id").
		Column("import.*").
		First()
	return &a, err
}

// GetByUUID an account by UUID.
func (s *ImportStore) GetByUUID(uuid string) (*models.Import, error) {
	a := models.Import{UUID: uuid}
	err := s.db.Model(&a).
		Where("import.uuid = ?uuid").
		Column("import.*").
		Relation("Room").
		First()
	return &a, err
}

// ApplySearch applies an RoomSearch on an orm.Query.
func (f *ImportFilter) ApplySearch(q *orm.Query) (*orm.Query, error) {

	if len(f.Query) > 0 {
		q = q.Where("LOWER(unaccent(code)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%") // ORDER  BY col_a ILIKE '%keyword%' OR NULL
	}

	return q, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *ImportFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.ApplySearch)
	//q = q.Where("name LIKE ?", "%"+f.Query[0]+"%")
	q = q.Order(f.Order...)
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *ImportStore) List(f *ImportFilter) ([]models.Import, int, error) {
	a := []models.Import{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		Relation("Room").
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// ImportFilter provides pagination and filtering options on accounts.
type ImportFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
	Query  []string
}

// NewImportFilter returns an ImportFilter with options parsed from request url values.
func NewImportFilter(params interface{}) (*ImportFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &ImportFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
		Query:  p["q"],
	}
	return f, nil
}

// Update an import.
func (s *ImportStore) Update(a *models.Import) error {
	_, err := s.db.Model(a).
		Column("name").
		WherePK().
		Update()
	return err
}
