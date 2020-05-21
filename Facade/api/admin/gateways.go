package admin

import (
	"context"
	"errors"
	clog "log"
	"net/http"

	"gitlab.com/IIIS/backend/facade/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/IIIS/backend/facade/database"
)

// The list of error types returned from account resource.
var (
	ErrGatewayValidation = errors.New("gateway validation error")
)

// GatewayStore defines database operations for account management.
type GatewayStore interface {
	List(*database.GatewayFilter) ([]models.Gateway, int, error)
	ListUnused(*database.GatewayFilter) ([]models.Gateway, int, error)
	Create(*models.Gateway) error
	Get(id int) (*models.Gateway, error)
	GetByUUID(uuid string) (*models.Gateway, error)
	Update(*models.Gateway) error
	Delete(*models.Gateway) error
	DeleteRoom(*models.Gateway) (int, error)
}

// GatewayResource implements account management handler.
type GatewayResource struct {
	Store GatewayStore
	Room  RoomStore
}

// NewGatewayResource creates and returns an account resource.
func NewGatewayResource(store GatewayStore, room RoomStore) *GatewayResource {
	return &GatewayResource{
		Store: store,
		Room:  room,
	}
}

func (rs *GatewayResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", rs.list)
	r.Get("/unused", rs.listUnused)
	r.Post("/", rs.create)
	r.Route("/{gatewayID}", func(r chi.Router) {
		r.Use(rs.gatewayCtx)
		r.Get("/", rs.get)
		r.Put("/", rs.update)
		r.Delete("/", rs.delete)
		r.Route("/room", func(r chi.Router) {
			r.Delete("/", rs.deleteRoom)
			r.Route("/{roomID}", func(r chi.Router) {
				r.Use(rs.roomCtx)
				r.Put("/", rs.setRoom)
			})
		})
	})
	return r
}

func (rs *GatewayResource) gatewayCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "gatewayID")
		gateway, err := rs.Store.GetByUUID(id)
		log(r).Println("kktk: ", gateway)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxGateway, gateway)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *GatewayResource) roomCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "roomID")
		room, err := rs.Room.GetByUUID(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxRoom, room)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type gatewayRequest struct {
	*models.Gateway
}

func (d *gatewayRequest) Bind(r *http.Request) error {
	return nil
}

type gatewayResponse struct {
	*models.Gateway
}

func newGatewayResponse(r *models.Gateway) *gatewayResponse {
	resp := &gatewayResponse{Gateway: r}
	return resp
}

type gatewayListResponse struct {
	Gateways []models.Gateway `json:"gateways"`
	Count    int              `json:"count"`
}

func newGatewayListResponse(a []models.Gateway, count int) *gatewayListResponse {
	resp := &gatewayListResponse{
		Gateways: a,
		Count:    count,
	}
	return resp
}

func (rs *GatewayResource) list(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewGatewayFilter(r.URL.Query(), false)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newGatewayListResponse(al, count))
}

func (rs *GatewayResource) listUnused(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewGatewayFilter(r.URL.Query(), true)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	al, count, err := rs.Store.ListUnused(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newGatewayListResponse(al, count))
}

func (rs *GatewayResource) setRoom(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)
	render.Respond(w, r, newRoomResponse(rm))
}

func (rs *GatewayResource) deleteRoom(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxGateway).(*models.Gateway)

	ar, err := rs.Store.DeleteRoom(rm)

	clog.Println("deleteRooom", ar, err)
	if err != nil {
		clog.Println("ERR deleteRooom", err)
		render.Render(w, r, ErrRender(err))
	}

	render.Respond(w, r, http.NoBody)
}

func (rs *GatewayResource) get(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxGateway).(*models.Gateway)
	render.Respond(w, r, newGatewayResponse(rm))
}

func (rs *GatewayResource) delete(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxGateway).(*models.Gateway)
	if err := rs.Store.Delete(rm); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}

func (rs *GatewayResource) create(w http.ResponseWriter, r *http.Request) {
	data := &gatewayRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Create(data.Gateway); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrGatewayValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, newGatewayResponse(data.Gateway))
}

func (rs *GatewayResource) update(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxGateway).(*models.Gateway)
	data := &gatewayRequest{Gateway: rm}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Update(rm); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrGatewayValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, newGatewayResponse(rm))
}
