package database

import (
	clog "log"

	"github.com/go-pg/pg"
	"gitlab.com/IIIS/backend/facade/models"
)

// FileStore implements database operations for account management by user.
type FileStore struct {
	db *pg.DB
}

// NewRoomStore returns an RoomStore.
func NewFileStore(db *pg.DB) *FileStore {
	return &FileStore{
		db: db,
	}
}

// Get an account by ID.
func (s *FileStore) Get(id int) (*models.File, error) {
	a := models.File{ID: id}
	err := s.db.Model(&a).
		Where("file.id = ?id").
		Column("file.*").
		Relation("Account"). /*
			Relation("Accounts", func(q *orm.Query) (*orm.Query, error) {
				return q.Where("room.active_user = account.id"), nil
			}).*/
		First()
	return &a, err
}

// Get an account by ID.
func (s *FileStore) GetPure(id int) (*models.File, error) {
	a := models.File{ID: id}
	err := s.db.Model(&a).
		Where("room.id = ?id").
		Column("room.*"). /*
			Relation("Accounts", func(q *orm.Query) (*orm.Query, error) {
				return q.Where("room.active_user = account.id"), nil
			}).*/
		First()
	return &a, err
}

// GetByUUID an account by UUID.
func (s *FileStore) GetByUUID(uuid string) (*models.File, error) {
	a := models.File{UUID: uuid}
	err := s.db.Model(&a).
		Where("file.uuid = ?uuid").
		Column("file.*").
		Relation("Account").
		First()

	clog.Println("Err: %v", err)
	return &a, err
}

// Update an account.
func (s *FileStore) Update(a *models.File) error {
	_, err := s.db.Model(a).
		Column("name").
		WherePK().
		Update()
	return err
}

// Delete an account.
func (s *FileStore) Delete(a *models.File) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Delete(a)
	})
	return err
}

// Create an account.
func (s *FileStore) Create(a *models.File) error {

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(a)
	})

	return err
}

/*
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

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(a)

	})

	return err
}
*/
