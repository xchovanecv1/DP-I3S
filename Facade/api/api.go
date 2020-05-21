// Package api configures an http server for administration and application resources.
package api

import (
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/robfig/cron/v3"
	"gitlab.com/IIIS/backend/facade/api/admin"
	"gitlab.com/IIIS/backend/facade/api/app"
	"gitlab.com/IIIS/backend/facade/api/gateway"
	"gitlab.com/IIIS/backend/facade/auth/jwt"
	"gitlab.com/IIIS/backend/facade/auth/pwdless"
	"gitlab.com/IIIS/backend/facade/database"
	"gitlab.com/IIIS/backend/facade/email"
	"gitlab.com/IIIS/backend/facade/logging"
)

// New configures application resources and routes.
func New(enableCORS bool, urlpfx string) (*chi.Mux, error) {
	logger := logging.NewLogger()

	if len(urlpfx) > 0 && urlpfx[0] != '/' {
		urlpfx = "/" + urlpfx
	} else {
		urlpfx = ""
	}

	db, err := database.DBConn()
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	mailer, err := email.NewMailer()
	if err != nil {
		logger.WithField("module", "email").Error(err)
		return nil, err
	}

	authStore := database.NewAuthStore(db)
	authResource, err := pwdless.NewResource(authStore, mailer)
	if err != nil {
		logger.WithField("module", "auth").Error(err)
		return nil, err
	}

	adminAPI, err := admin.NewAPI(db, mailer)
	if err != nil {
		logger.WithField("module", "admin").Error(err)
		return nil, err
	}

	appAPI, err := app.NewAPI(db)
	if err != nil {
		logger.WithField("module", "app").Error(err)
		return nil, err
	}

	gatewayAPI, err := gateway.NewAPI(db)
	if err != nil {
		logger.WithField("module", "gateway").Error(err)
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	//r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// use CORS middleware if client is not served by this api, e.g. from other domain or CDN
	if enableCORS {
		r.Use(corsConfig().Handler)
	}

	r.Mount(urlpfx+"/auth", authResource.Router())
	r.Group(func(r chi.Router) {
		r.Use(authResource.TokenAuth.Verifier())
		r.Use(jwt.Authenticator)
		r.Mount(urlpfx+"/admin", adminAPI.Router())
		r.Mount(urlpfx+"/api", appAPI.Router())
	})

	// No authorization required, TODO, validate allow only local requests
	r.Group(func(r chi.Router) {
		r.Mount(urlpfx+"/gateway", gatewayAPI.Router())
	})

	r.Get(urlpfx+"/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	client := "./public"
	r.Get("/*", SPAHandler(client))

	c := cron.New()
	c.AddFunc("* * * * *", func() {
		//adminAPI.Properties.
	})
	c.Start()

	return r, nil
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}

// SPAHandler serves the public Single Page Application.
func SPAHandler(publicDir string) http.HandlerFunc {
	handler := http.FileServer(http.Dir(publicDir))
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(publicDir, "index.html")
		serviceWorker := path.Join(publicDir, "service-worker.js")

		requestedAsset := path.Join(publicDir, r.URL.Path)
		if strings.Contains(requestedAsset, "service-worker.js") {
			requestedAsset = serviceWorker
		}
		if _, err := os.Stat(requestedAsset); err != nil {
			http.ServeFile(w, r, indexPage)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
