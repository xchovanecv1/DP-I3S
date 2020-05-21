package admin

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gitlab.com/IIIS/backend/facade/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/IIIS/backend/facade/database"
)

// The list of error types returned from account resource.
var (
	ErrRoomValidation = errors.New("room validation error")
)

// RoomStore defines database operations for account management.
type RoomStore interface {
	List(*database.RoomFilter) ([]models.Room, int, error)
	Create(*models.Room) error
	Get(id int) (*models.Room, error)
	GetByUUID(uuid string) (*models.Room, error)
	GetCodePure(code string) (*models.Room, error)
	Update(*models.Room) error
	Delete(*models.Room) error
}

// RoomResource implements account management handler.
type RoomResource struct {
	Store   RoomStore
	Gateway *database.GatewayStore
}

// NewRoomResource creates and returns an account resource.
func NewRoomResource(store RoomStore, gw *database.GatewayStore) *RoomResource {
	return &RoomResource{
		Store:   store,
		Gateway: gw,
	}
}

func (rs *RoomResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", rs.list)
	r.Post("/", rs.create)
	r.Route("/{roomID}", func(r chi.Router) {
		r.Use(rs.roomCtx)
		r.Get("/", rs.get)
		r.Put("/", rs.update)
		r.Delete("/", rs.delete)

		r.Route("/gateway", func(r chi.Router) {
			r.Get("/", rs.getGateways)
			r.Post("/", rs.addGateway)
		})
	})
	return r
}

func (rs *RoomResource) roomCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "roomID")
		iid, err := strconv.Atoi(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		room, err := rs.Store.Get(iid)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxRoom, room)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type roomRequest struct {
	*models.Room
}

func (d *roomRequest) Bind(r *http.Request) error {
	return nil
}

type roomResponse struct {
	*models.Room
}

func newRoomResponse(r *models.Room) *roomResponse {
	resp := &roomResponse{Room: r}
	return resp
}

type roomListResponse struct {
	Rooms []models.Room `json:"rooms"`
	Count int           `json:"count"`
}

func newRoomListResponse(a []models.Room, count int) *roomListResponse {
	resp := &roomListResponse{
		Rooms: a,
		Count: count,
	}
	return resp
}

func (rs *RoomResource) list(w http.ResponseWriter, r *http.Request) {
	f, err := database.NewRoomFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newRoomListResponse(al, count))
}

func (rs *RoomResource) get(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)
	render.Respond(w, r, newRoomResponse(rm))
}

func (rs *RoomResource) getGateways(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)

	gws, err := rs.Gateway.GetByRoom(rm.ID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, gws)
}

type gwRequest struct {
	*models.Gateway
}

func (d *gwRequest) Bind(r *http.Request) error {
	return nil
}

func (rs *RoomResource) addGateway(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)

	data := &gwRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	gw, err := rs.Gateway.Get(data.ID)

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if gw.ID == 0 || gw.RoomID != 0 {
		render.Render(w, r, ErrUnprocessableEntity)
		return
	}

	gw.RoomID = rm.ID

	err = rs.Gateway.Update(gw)
	if err != nil || gw.ID == 0 {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, gw)
}

func (rs *RoomResource) delete(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)
	if err := rs.Store.Delete(rm); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}

func (rs *RoomResource) create(w http.ResponseWriter, r *http.Request) {
	data := &roomRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Create(data.Room); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrRoomValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, newRoomResponse(data.Room))
}

func (rs *RoomResource) update(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value(ctxRoom).(*models.Room)
	data := &roomRequest{Room: rm}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Update(rm); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrRoomValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, newRoomResponse(rm))
}
