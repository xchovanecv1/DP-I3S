// Package admin ties together administration resources and handlers.
package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-pg/pg"

	"gitlab.com/IIIS/backend/facade/auth/authorize"
	"gitlab.com/IIIS/backend/facade/auth/jwt"
	"gitlab.com/IIIS/backend/facade/database"
	"gitlab.com/IIIS/backend/facade/email"
	"gitlab.com/IIIS/backend/facade/logging"
	"gitlab.com/IIIS/backend/facade/models"
)

const (
	roleAdmin = "admin"
)

type ctxKey int

const (
	ctxAccount  ctxKey = iota
	ctxRoom     ctxKey = iota
	ctxCard     ctxKey = iota
	ctxGateway  ctxKey = iota
	ctxFile     ctxKey = iota
	ctxProperty ctxKey = iota
)

// Mailer defines methods to send account emails.
type Mailer interface {
	UnatuhNorification(name, email string, c email.ContentUnauthNotif) error
}

// API provides admin application resources and handlers.
type API struct {
	Accounts   *AccountResource
	Rooms      *RoomResource
	Cards      *CardResource
	Gateways   *GatewayResource
	Files      *FileResource
	Properties *PropertyResource
	Imports    *ImportResource
	Mailer     *email.Mailer
}

// NewAPI configures and returns admin application API.
func NewAPI(db *pg.DB, mailer *email.Mailer) (*API, error) {

	cardLogStore := database.NewCardLogStore(db)
	cardStore := database.NewCardStore(db)

	accountStore := database.NewAdmAccountStore(db)
	accounts := NewAccountResource(accountStore, cardStore, cardLogStore)

	importStore := database.NewImportStore(db)
	imports := NewImportResource(importStore)

	gatewayStore := database.NewGatewayStore(db)
	roomStore := database.NewRoomStore(db)
	rooms := NewRoomResource(roomStore, gatewayStore)

	cards := NewCardResource(cardStore, cardLogStore)

	gateways := NewGatewayResource(gatewayStore, roomStore)

	filesStore := database.NewFileStore(db)
	files := NewFileResource(filesStore, accountStore)

	propertiesStore := database.NewPropertyStore(db)
	propertiesTransitionStore := database.NewPropertyTransitStore(db)
	properties := NewPropertyResource(propertiesStore, propertiesTransitionStore, accountStore, files, rooms, cards, cardLogStore, imports)

	api := &API{
		Accounts:   accounts,
		Rooms:      rooms,
		Cards:      cards,
		Gateways:   gateways,
		Files:      files,
		Properties: properties,
		Imports:    imports,
		Mailer:     mailer,
	}

	c := cron.New()
	c.AddFunc("*/5 * * * *", func() {
		nw := time.Now()
		nt := viper.GetDuration("prop_unauth_notif")
		nw = nw.Add(-nt)

		tras, _ := api.Properties.Transition.GetExitAfterDuration(nw)

		adm, _ := api.Accounts.Store.GetByRole("admin")

		for i, t := range *tras {
			(*tras)[i].PostAuth = true
			if t.UserID != 0 {
				go func(t models.PropertyTransit) {
					content := email.ContentAuthUser{
						Email:    t.UserEmail,
						Name:     t.UserName,
						ItemName: t.PropertyName,
						RoomName: t.RoomName,
					}
					if err := api.Mailer.AuthUserNotif(t.UserName, t.UserEmail, content); err != nil {
						fmt.Println(err)
					}

				}(t)
				go func(t models.PropertyTransit) {
					content := email.ContentAdminAuth{
						Email:    adm.Email,
						Name:     adm.Name,
						UserName: t.UserName,
						ItemName: t.PropertyName,
						RoomName: t.RoomName,
					}
					if err := api.Mailer.AdminAuthNotif(adm.Name, adm.Email, content); err != nil {
						fmt.Println(err)
					}

				}(t)
			} else {
				go func(t models.PropertyTransit) {
					content := email.ContentUnauthNotif{
						Email:    adm.Email,
						Name:     adm.Name,
						ItemName: t.PropertyName,
						RoomName: t.RoomName,
					}
					if err := api.Mailer.UnatuhNorification(adm.Name, adm.Email, content); err != nil {
						fmt.Println(err)
					}

				}(t)
			}
		}

		api.Properties.Transition.UpdateTransits(tras)
	})
	c.Start()

	return api, nil
}

func (a *API) accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := jwt.ClaimsFromCtx(r.Context())
		fmt.Println(claims.ID)
		account, err := a.Accounts.Store.GetByUUID(claims.ID)
		if err != nil {
			// account deleted while access token still valid
			render.Render(w, r, ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxAccount, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Router provides admin application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(authorize.RequiresRole(roleAdmin))
	r.Use(a.accountCtx)
	/*r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Admin"))
	})*/

	r.Mount("/accounts", a.Accounts.router())
	r.Mount("/import", a.Imports.router())
	r.Mount("/rooms", a.Rooms.router())
	r.Mount("/cards", a.Cards.router())
	r.Mount("/gateways", a.Gateways.router())
	r.Mount("/files", a.Files.router())
	r.Mount("/props", a.Properties.router())
	return r
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
