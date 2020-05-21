package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/IIIS/backend/facade/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/IIIS/backend/facade/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// The list of error types returned from account resource.
var (
	ErrAccountValidation = errors.New("account validation error")
)

// AccountStore defines database operations for account management.
type AccountStore interface {
	List(*database.AccountFilter) ([]models.Account, int, error)
	Create(*models.Account) error
	Get(id int) (*models.Account, error)
	GetByUUID(uuid string) (*models.Account, error)
	GetByRole(role string) (*models.Account, error)
	Update(*models.Account) error
	Delete(*models.Account) error
}

// AccountResource implements account management handler.
type AccountResource struct {
	Store    AccountStore
	Cards    *database.CardStore
	CardLogs *database.CardLogStore
}

// NewAccountResource creates and returns an account resource.
func NewAccountResource(store AccountStore, cs *database.CardStore, cls *database.CardLogStore) *AccountResource {
	return &AccountResource{
		Store:    store,
		Cards:    cs,
		CardLogs: cls,
	}
}

func (rs *AccountResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", rs.list)
	r.Post("/", rs.create)
	r.Route("/{accountID}", func(r chi.Router) {
		r.Use(rs.accountCtx)
		r.Get("/", rs.get)
		r.Put("/", rs.update)
		r.Delete("/", rs.delete)
		r.Route("/cards", func(r chi.Router) {
			r.Get("/", rs.getCards)
			r.Delete("/", rs.deleteCardLink)
			//r.Post("/", rs.addCardToProp)
			r.Post("/log", rs.addCardFromLog)
		})
	})
	return r
}

func (rs *AccountResource) accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "accountID")
		account, err := rs.Store.GetByUUID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxAccount, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *AccountResource) getCards(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)

	cards, count, err := rs.Cards.ListLink("Account", acc.ID)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewCardListResponse(cards, count))
}

func (rs *AccountResource) addCardFromLog(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)

	data := &cardLogRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	cl := data.CardLog

	if cl.ID == 0 {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}

	log, err := rs.CardLogs.Get(cl.ID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	card := &models.Card{
		Link:   "Account",
		LinkID: acc.ID,
	}

	err = rs.Cards.CreateFromLog(card, log)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = rs.CardLogs.DeleteByCode(log)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	/*
		if err := rs.Store.Update(rm); err != nil {
			switch err.(type) {
			case validation.Errors:
				render.Render(w, r, ErrValidation(ErrPropertyValidation, err.(validation.Errors)))
				return
			}
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
	*/
	render.Respond(w, r, card)
}

func (rs *AccountResource) deleteCardLink(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)

	data := &cardRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	fmt.Printf("%+v\n", data.Card)

	cl := data.Card

	if len(cl.UUID) == 0 {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}

	card, err := rs.Cards.GetByUUID(cl.UUID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if card.Link != "Account" || card.LinkID != acc.ID {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}

	card.Link = ""
	card.LinkID = 0

	err = rs.Cards.Update(card)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, http.NoBody)
}

type accountRequest struct {
	*models.Account
}

func (d *accountRequest) Bind(r *http.Request) error {
	return nil
}

type accountResponse struct {
	*models.Account
}

func newAccountResponse(a *models.Account) *accountResponse {
	resp := &accountResponse{Account: a}
	return resp
}

type accountListResponse struct {
	Accounts []models.Account `json:"accounts"`
	Count    int              `json:"count"`
}

func newAccountListResponse(a []models.Account, count int) *accountListResponse {
	resp := &accountListResponse{
		Accounts: a,
		Count:    count,
	}
	return resp
}

func (rs *AccountResource) list(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewAccountFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newAccountListResponse(al, count))
}

func (rs *AccountResource) create(w http.ResponseWriter, r *http.Request) {
	data := &accountRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Create(data.Account); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrAccountValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, newAccountResponse(data.Account))
}

func (rs *AccountResource) get(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)
	render.Respond(w, r, newAccountResponse(acc))
}

func (rs *AccountResource) update(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)
	data := &accountRequest{Account: acc}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Update(acc); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrAccountValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, newAccountResponse(acc))
}

func (rs *AccountResource) delete(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)
	if err := rs.Store.Delete(acc); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}
