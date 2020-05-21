package database

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/spf13/viper"
	"gitlab.com/IIIS/backend/facade/models"
)

// PropertyTransitStore implements database operations for account management by user.
type PropertyTransitStore struct {
	db *pg.DB
}

// NewPropertyTransitStore returns an DoorAccessStore.
func NewPropertyTransitStore(db *pg.DB) *PropertyTransitStore {
	return &PropertyTransitStore{
		db: db,
	}
}

// AddPropertyTransit adds new acces requesto to the DB
func (s *PropertyTransitStore) AddPropertyTransit(a *models.PropertyTransit) error {
	log.Printf("Data to add: %v\n", a)
	err := s.db.Insert(a)
	log.Println(err)
	return err
}

// Create creates a new PropertyTransit.
func (s *PropertyTransitStore) Create(a *models.PropertyTransit) error {
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

// CreateMultiple asd
func (s *PropertyTransitStore) CreateMultiple(a *[]models.PropertyTransit) error {
	/*
		count, _ := s.db.Model(a).
			Where("name = ?name").
			Count()

		if count != 0 {
			return ErrUniqueEmailConstraint
		}
	*/
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		var err error
		for _, t := range *a {
			err = tx.Insert(&t)
			if err != nil {
				return err
			}
		}
		return nil
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

// Delete an PropertyTransit.
func (s *PropertyTransitStore) Delete(a *models.PropertyTransit) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// Get an account by ID.
func (s *PropertyTransitStore) Get(id int) (*models.PropertyTransit, error) {
	a := models.PropertyTransit{ID: id}
	err := s.db.Model(&a).
		Where("property_transit.id = ?id").
		Column("property_transit.*").
		First()
	return &a, err
}

// GetLatestByProp an account by ID.
func (s *PropertyTransitStore) GetLatestByProp(id int) (*models.PropertyTransit, error) {
	a := models.PropertyTransit{PropertyID: id}
	err := s.db.Model(&a).
		Column("property_transit.*").
		Where("property_id = ?property_id").
		Last()
	return &a, err
}

// GetLastestUserExits an account by ID.
func (s *PropertyTransitStore) GetLastestUserExits(uid int) (*[]models.PropertyTransit, error) {

	ret := []models.PropertyTransit{}
	a := []models.PropertyTransit{}
	_, err := s.db.Query(&a, "SELECT DISTINCT ON (pt.property_id) pt.*, property.name as \"property_name\", property.code as \"property_code\" FROM property_transits pt LEFT JOIN properties as \"property\" ON pt.property_id = property.id WHERE pt.user_id = ? ORDER BY pt.property_id DESC, pt.id DESC", uid)

	for _, t := range a {
		if t.Transition == "EXIT" {
			ret = append(ret, t)
		}
	}

	return &ret, err
}

// GetExitAfterDuration an account by ID.
func (s *PropertyTransitStore) GetExitAfterDuration(t time.Time) (*[]models.PropertyTransit, error) {

	ret := []models.PropertyTransit{}
	a := []models.PropertyTransit{}
	_, err := s.db.Query(&a, "SELECT DISTINCT ON (pt.property_id) pt.*, property.name as \"property_name\", property.code as \"property_code\", account.name as \"user_name\", account.email as \"user_email\", room.name as \"room_name\" FROM property_transits pt LEFT JOIN properties as \"property\" ON pt.property_id = property.id  LEFT JOIN rooms as \"room\" ON pt.room_id = room.id LEFT JOIN accounts as \"account\" ON pt.user_id = account.id WHERE pt.created_at <= ? AND pt.post_auth = false ORDER BY pt.property_id DESC, pt.id DESC", t)

	for _, t := range a {
		if t.Transition == "EXIT" {
			ret = append(ret, t)
		}
	}

	return &ret, err
}

// UpdateUserID an account by ID.
func (s *PropertyTransitStore) UpdateUserID(a *models.PropertyTransit) error {
	_, err := s.db.Model(a).
		Column("user_id").
		WherePK().
		Update()
	return err
}

// UpdateTransits an account by ID.
func (s *PropertyTransitStore) UpdateTransits(a *[]models.PropertyTransit) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		for _, pt := range *a {
			err := tx.Update(&pt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// AuthPriorTransits authentificate any old transitions within room_auth_expiry period
func (s *PropertyTransitStore) AuthPriorTransits(rid int, acid int) error {
	a := []models.PropertyTransit{}

	tp := -viper.GetDuration("prop_auth_period")
	nw := time.Now()
	nw = nw.Add(tp)

	err := s.db.Model(&a).
		Where("user_id IS NULL").
		Where("room_id = ?", rid).
		Where("created_at > ?", nw).
		Select()

	if err != nil {
		return err
	}
	for i := range a {
		a[i].UserID = acid
		a[i].PostAuth = true
	}

	return s.UpdateTransits(&a)
}

/*
// GetByUUID an account by UUID.
func (s *PropertyTransitStore) GetByUUID(uuid string) (*models.PropertyTransit, error) {
	a := models.PropertyTransit{UUID: uuid}
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
func (s *PropertyTransitStore) List(f *CardFilter) ([]models.Card, int, error) {
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

// NewPropertyTansitFilter returns an CardFilter with options parsed from request url values.
func NewPropertyTansitFilter(params interface{}) (*CardFilter, error) {
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
*/
/*
// Update an card.
func (s *PropertyTransitStore) Update(a *models.Card) error {
	_, err := s.db.Model(a).
		Column("name").
		WherePK().
		Update()
	return err
}*/
