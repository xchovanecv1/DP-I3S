package admin

import (
	"context"
	"errors"
	"net/http"

	clog "log"

	"gitlab.com/IIIS/backend/facade/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/IIIS/backend/facade/database"
)

// The list of error types returned from account resource.
var (
	ErrCardValidation = errors.New("Card validation error")
)

// CardStore defines database operations for account management.
type CardStore interface {
	List(*database.CardFilter) ([]models.Card, int, error)
	ListLink(linkType string, linkID int) ([]models.Card, int, error)
	Create(*models.Card) error
	CreateFromLog(*models.Card, *models.CardLog) error
	Get(id int) (*models.Card, error)
	GetByUUID(uuid string) (*models.Card, error)
	Update(*models.Card) error
	Delete(*models.Card) error
}

// CardResource implements account management handler.
type CardResource struct {
	Store    CardStore
	CardLogs *database.CardLogStore
}

// NewCardResource creates and returns an account resource.
func NewCardResource(store CardStore, logs *database.CardLogStore) *CardResource {
	return &CardResource{
		Store:    store,
		CardLogs: logs,
	}
}

func (rs *CardResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", rs.list)
	r.Post("/", rs.create)
	r.Route("/logs", func(r chi.Router) {
		r.Get("/", rs.getCardLogs)
	})
	r.Route("/{CardID}", func(r chi.Router) {
		r.Use(rs.cardCtx)
		r.Get("/", rs.get)
		//r.Put("/", rs.update)
		r.Delete("/", rs.delete)
	})
	return r
}

/***/
// Card logs

//CardLogListResponse encapsulate response for card including response count
type CardLogListResponse struct {
	Cards []models.CardLog `json:"cards"`
}

//NewCardLogListResponse creates api response for card list
func NewCardLogListResponse(a []models.CardLog) *CardLogListResponse {
	resp := &CardLogListResponse{
		Cards: a,
	}
	return resp
}

func (rs *CardResource) getCardLogs(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewCardLogFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, err := rs.CardLogs.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, NewCardLogListResponse(al))
}

/***/

func (rs *CardResource) cardCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "CardID")
		card, err := rs.Store.GetByUUID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxCard, card)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type cardRequest struct {
	*models.Card
}

func (d *cardRequest) Bind(r *http.Request) error {
	return nil
}

type cardResponse struct {
	*models.Card
}

func newCardResponse(r *models.Card) *cardResponse {
	resp := &cardResponse{Card: r}
	return resp
}

//CardListResponse encapsulate response for card including response count
type CardListResponse struct {
	Cards []models.Card `json:"cards"`
	Count int           `json:"count"`
}

//NewCardListResponse creates api response for card list
func NewCardListResponse(a []models.Card, count int) *CardListResponse {
	resp := &CardListResponse{
		Cards: a,
		Count: count,
	}
	return resp
}

func (rs *CardResource) list(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewCardFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, NewCardListResponse(al, count))
}

func (rs *CardResource) get(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxCard).(*models.Card)
	render.Respond(w, r, newCardResponse(rm))
}

func (rs *CardResource) delete(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxCard).(*models.Card)
	if err := rs.Store.Delete(rm); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}

func (rs *CardResource) create(w http.ResponseWriter, r *http.Request) {
	data := &cardRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	clog.Println("%+v\n", data)

	if err := rs.Store.Create(data.Card); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrCardValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, newCardResponse(data.Card))
}

/*
func (rs *CardResource) update(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxCard).(*models.Card)
	data := &cardRequest{Card: rm}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Update(rm); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrCardValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, newCardResponse(rm))
}
*/
