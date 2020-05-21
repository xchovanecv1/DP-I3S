package database

import (
	clog "log"
	"net/url"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
	"gitlab.com/IIIS/backend/facade/models"
)

// GatewayStore implements database operations for account management by user.
type GatewayStore struct {
	db *pg.DB
}

// NewGatewayStore returns an GatewayStore.
func NewGatewayStore(db *pg.DB) *GatewayStore {
	return &GatewayStore{
		db: db,
	}
}

// Get an account by ID.
func (s *GatewayStore) Get(id int) (*models.Gateway, error) {
	a := models.Gateway{ID: id}
	err := s.db.Model(&a).
		Where("gateway.id = ?id").
		Column("gateway.*").
		First()
	return &a, err
}

// GetByRoom by room id
func (s *GatewayStore) GetByRoom(id int) (*[]models.Gateway, error) {
	a := []models.Gateway{}
	err := s.db.Model(&a).
		Where("gateway.room_id = ?", id).
		Column("gateway.*").
		Select()
	return &a, err
}

// GetByUUID an account by UUID.
func (s *GatewayStore) GetByUUID(uuid string) (*models.Gateway, error) {
	a := models.Gateway{UUID: uuid}
	err := s.db.Model(&a).
		Where("gateway.uuid = ?uuid").
		Column("gateway.*").
		Relation("Room").
		First()

	clog.Println("Err: %v", err)
	return &a, err
}

// GetByAddress retrieves Gateway object by its address property
func (s *GatewayStore) GetByAddress(addr string) (*models.Gateway, error) {
	a := models.Gateway{Addr: addr}
	err := s.db.Model(&a).
		Where("gateway.addr = ?addr").
		Column("gateway.*").
		First()
	return &a, err
}

// Update an account.
func (s *GatewayStore) Update(a *models.Gateway) error {
	_, err := s.db.Model(a).
		Column("name").
		Column("last_active").
		Column("room_id").
		WherePK().
		Update()

	return err
}

// SetLastActive an gaeway.
func (s *GatewayStore) SetLastActive(addr string) error {

	gw, err := s.GetByAddress(addr)
	now := time.Now()
	if err != nil {
		if gw.ID == 0 {
			gw.Addr = addr
			gw.LastActive = now
			err = s.Create(gw)

			return err
		} else {
			return err
		}
	}

	gw.LastActive = now

	err = s.Update(gw)

	return err
}

// Delete an account.
func (s *GatewayStore) Delete(a *models.Gateway) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// DeleteRoom sets assigned room in table to
func (s *GatewayStore) DeleteRoom(a *models.Gateway) (int, error) {
	gw, err := s.db.Model(a).
		Set("room_id = null").
		Where("uuid = ?uuid").
		Update()
	return gw.RowsAffected(), err
}

// ApplySearch applies an GatewaySearch on an orm.Query.
func (f *GatewayFilter) ApplySearch(q *orm.Query) (*orm.Query, error) {

	if len(f.Query) > 0 && len(f.Query[0]) > 0 {
		q = q.Where("LOWER(unaccent(name)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%") // ORDER  BY col_a ILIKE '%keyword%' OR NULL
		q = q.WhereOr("LOWER(unaccent(addr)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%")
	}

	return q, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *GatewayFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.ApplySearch)

	if f.Unused {
		q = q.Where("room_id IS NULL")
	} else {
		q = q.Where("room_id IS NOT NULL")
	}
	//q = q.Where("name LIKE ?", "%"+f.Query[0]+"%")
	q = q.Order(f.Order...)
	q = q.Order("last_active DESC")
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *GatewayStore) List(f *GatewayFilter) ([]models.Gateway, int, error) {
	a := []models.Gateway{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *GatewayStore) ListUnused(f *GatewayFilter) ([]models.Gateway, int, error) {
	a := []models.Gateway{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// GatewayFilter provides pagination and filtering options on accounts.
type GatewayFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
	Query  []string
	Unused bool
}

// NewGatewayFilter returns an GatewayFilter with options parsed from request url values.
func NewGatewayFilter(params interface{}, unused bool) (*GatewayFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &GatewayFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
		Query:  p["q"],
		Unused: unused,
	}
	return f, nil
}

// Create creates a new gateway.
func (s *GatewayStore) Create(a *models.Gateway) error {
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
