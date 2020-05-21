package admin

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.com/IIIS/backend/facade/database"
	"gitlab.com/IIIS/backend/facade/models"
)

// The list of error types returned from account resource.
var (
	ErrImportValidation = errors.New("import validation error")
)

// ImportStore defines database operations for account management.
type ImportStore interface {
	List(*database.ImportFilter) ([]models.Import, int, error)
	Create(*models.Import) error
	Get(id int) (*models.Import, error)
	GetByUUID(uuid string) (*models.Import, error)
	Update(*models.Import) error
	Delete(*models.Import) error
}

// ImportResource implements account management handler.
type ImportResource struct {
	Store ImportStore
}

// NewImportResource creates and returns an account resource.
func NewImportResource(store ImportStore) *ImportResource {
	return &ImportResource{
		Store: store,
	}
}

func (rs *ImportResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/actions", rs.actions)
	return r
}

type importAction struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type importActionType map[string][]importAction

func (rs *ImportResource) actions(w http.ResponseWriter, r *http.Request) {

	ret := make(map[string][]importAction)

	var roomActions [2]importAction

	roomActions[0] = importAction{
		Name: "Vyhľadať",
		Code: "F_RM_ID_CODE",
	}
	roomActions[1] = importAction{
		Name: "Vyhľadať alebo vytvoriť",
		Code: "FoC_RM_ID_CODE",
	}

	ret["room_id"] = roomActions[:]

	render.Respond(w, r, ret)
}

func (rs *ImportResource) logImport(r *http.Request, resource string, ids []int) error {

	user := r.Context().Value(ctxAccount).(*models.Account)

	data := &models.Import{
		CreatedByID: user.ID,
		Resource:    resource,
		Elements:    ids,
	}

	return rs.Store.Create(data)

}
