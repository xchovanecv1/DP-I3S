package database

import (
	"log"
	"net/url"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
	"gitlab.com/IIIS/backend/facade/models"
)

// PropertyStore implements database operations for account management by user.
type PropertyStore struct {
	db *pg.DB
}

// NewPropertyStore returns an DoorAccessStore.
func NewPropertyStore(db *pg.DB) *PropertyStore {
	return &PropertyStore{
		db: db,
	}
}

// AddProperty adds new acces requesto to the DB
func (s *PropertyStore) AddProperty(a *models.Property) error {
	log.Printf("Data to add: %v\n", a)
	err := s.db.Insert(a)
	log.Println(err)
	return err
}

// GetByCode retrieves Property object by its code property
func (s *PropertyStore) GetByCode(code string) (*models.Property, error) {
	a := models.Property{Code: code}
	err := s.db.Model(&a).
		Where("property.code = ?code").
		Column("property.*").
		First()
	return &a, err
}

// Create creates a new property.
func (s *PropertyStore) Create(a *models.Property) error {
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
func (s *PropertyStore) CreateMultiple(a []*models.Property) error {

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

// Delete an property.
func (s *PropertyStore) Delete(a *models.Property) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// Get an account by ID.
func (s *PropertyStore) Get(id int) (*models.Property, error) {
	a := models.Property{ID: id}
	err := s.db.Model(&a).
		Where("property.id = ?id").
		Column("property.*").
		First()
	return &a, err
}

// GetByUUID an account by UUID.
func (s *PropertyStore) GetByUUID(uuid string) (*models.Property, error) {
	a := models.Property{UUID: uuid}
	err := s.db.Model(&a).
		Where("property.uuid = ?uuid").
		Column("property.*").
		Relation("Room").
		First()
	return &a, err
}

// ApplySearch applies an RoomSearch on an orm.Query.
func (f *PropertyFilter) ApplySearch(q *orm.Query) (*orm.Query, error) {

	if len(f.Query) > 0 {
		q = q.Where("LOWER(unaccent(property.name)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%")
		q = q.WhereOr("LOWER(unaccent(property.code)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%")
		q = q.WhereOr("LOWER(unaccent(room.name)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%")
		q = q.WhereOr("LOWER(unaccent(room.code)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%") // ORDER  BY col_a ILIKE '%keyword%' OR NULL
	}

	return q, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *PropertyFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.ApplySearch)
	//q = q.Where("name LIKE ?", "%"+f.Query[0]+"%")
	q = q.Order(f.Order...)
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *PropertyStore) List(f *PropertyFilter) ([]models.Property, int, error) {
	a := []models.Property{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		Relation("Room").
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// PropertyFilter provides pagination and filtering options on accounts.
type PropertyFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
	Query  []string
}

// NewPropertyFilter returns an PropertyFilter with options parsed from request url values.
func NewPropertyFilter(params interface{}) (*PropertyFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &PropertyFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
		Query:  p["q"],
	}
	return f, nil
}

// Update an property.
func (s *PropertyStore) Update(a *models.Property) error {
	_, err := s.db.Model(a).
		Column("name").
		Column("code").
		Column("price").
		Column("acquired_at").
		Column("room_id").
		WherePK().
		Update()
	return err
}
