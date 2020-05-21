package database

import (
	clog "log"
	"net/url"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
	"github.com/spf13/viper"
	"gitlab.com/IIIS/backend/facade/models"
)

// RoomStore implements database operations for account management by user.
type RoomStore struct {
	db *pg.DB
}

// NewRoomStore returns an RoomStore.
func NewRoomStore(db *pg.DB) *RoomStore {
	return &RoomStore{
		db: db,
	}
}

// Get an account by ID.
func (s *RoomStore) Get(id int) (*models.Room, error) {
	a := models.Room{ID: id}
	err := s.db.Model(&a).
		Where("room.id = ?id").
		Column("room.*").
		Relation("ActiveUser"). /*
			Relation("Accounts", func(q *orm.Query) (*orm.Query, error) {
				return q.Where("room.active_user = account.id"), nil
			}).*/
		First()
	return &a, err
}

// GetPure an account by ID.
func (s *RoomStore) GetPure(id int) (*models.Room, error) {
	a := models.Room{ID: id}
	err := s.db.Model(&a).
		Where("room.id = ?id").
		Column("room.*"). /*
			Relation("Accounts", func(q *orm.Query) (*orm.Query, error) {
				return q.Where("room.active_user = account.id"), nil
			}).*/
		First()
	return &a, err
}

// GetCodePure an room by Code.
func (s *RoomStore) GetCodePure(code string) (*models.Room, error) {
	a := models.Room{Code: code}
	err := s.db.Model(&a).
		Where("room.code = ?code").
		Column("room.*"). /*
			Relation("Accounts", func(q *orm.Query) (*orm.Query, error) {
				return q.Where("room.active_user = account.id"), nil
			}).*/
		First()
	return &a, err
}

// GetByUUID an account by UUID.
func (s *RoomStore) GetByUUID(uuid string) (*models.Room, error) {
	a := models.Room{UUID: uuid}
	err := s.db.Model(&a).
		Where("room.uuid = ?uuid").
		Column("room.*").
		Relation("Gateways").
		First()

	clog.Println("Err: %v", err)
	return &a, err
}

// ValidateActiveUser check for valid active user id in ROOM within expiry period
func (s *RoomStore) ValidateActiveUser(id int) error {
	a := models.Room{ID: id}

	err := s.db.Model(&a).
		Where("room.id = ?id").
		Column("room.*").
		First()

	if err != nil {
		return err
	}

	to := viper.GetDuration("room_auth_expiry")
	d := a.UpdatedAt
	d.Add(to)
	/*
		n := time.Now()

			// timeout
			if n.After(d) {
				a.ActiveUserID = 0
				fmt.Println("Clearing active user")
			}
	*/
	err = s.UpdateActive(&a)

	return err
}

// Update an account.
func (s *RoomStore) Update(a *models.Room) error {
	_, err := s.db.Model(a).
		Column("name").
		WherePK().
		Update()
	return err
}

// UpdateActive an account.
func (s *RoomStore) UpdateActive(a *models.Room) error {
	_, err := s.db.Model(a).
		Column("active_user_id").
		WherePK().
		Update()
	return err
}

// Delete an account.
func (s *RoomStore) Delete(a *models.Room) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// ApplySearch applies an RoomSearch on an orm.Query.
func (f *RoomFilter) ApplySearch(q *orm.Query) (*orm.Query, error) {

	if len(f.Query) > 0 {
		q = q.
			WhereOr("LOWER(unaccent(name)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%").
			WhereOr("LOWER(unaccent(code)) LIKE LOWER(unaccent(?))", "%"+f.Query[0]+"%")
		words := strings.Fields(f.Query[0])
		if len(words) > 1 {
			for _, element := range words {
				// index is the index where we are
				// element is the element from someSlice for where we are
				q = q.
					WhereOr("LOWER(unaccent(name)) LIKE LOWER(unaccent(?))", "%"+element+"%").
					WhereOr("LOWER(unaccent(code)) LIKE LOWER(unaccent(?))", "%"+element+"%")
			}
		}
	}

	return q, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *RoomFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.Filter.Filters)
	q = q.Apply(f.ApplySearch)
	//q = q.Where("name LIKE ?", "%"+f.Query[0]+"%")
	q = q.Order(f.Order...)
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *RoomStore) List(f *RoomFilter) ([]models.Room, int, error) {
	a := []models.Room{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// RoomFilter provides pagination and filtering options on accounts.
type RoomFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
	Query  []string
}

// NewRoomFilter returns an RoomFilter with options parsed from request url values.
func NewRoomFilter(params interface{}) (*RoomFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &RoomFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
		Query:  p["q"],
	}
	return f, nil
}

// Create creates a new room.
func (s *RoomStore) Create(a *models.Room) error {
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
