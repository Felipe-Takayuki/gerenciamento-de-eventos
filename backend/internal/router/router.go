package router

import (
	"database/sql"
	"net/http"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/config"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/database"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/service"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/webserver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

func NewRouter(db *sql.DB, cfg *config.Config) http.Handler {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	institutionDB := database.NewInstitutionDB(db)
	eventDB := database.NewEventDB(db)

	institutionService := service.NewInstitutionService(institutionDB)
	eventService := service.NewEventService(eventDB)

	webInstitutionService := webserver.NewWebInstitutionHandler(institutionService)
	webEventService := webserver.NewWebEventHandler(eventService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	corsConfig := cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	r.Use(cors.Handler(corsConfig))

	// Health Check (Público)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Autenticação (Público - apenas login)
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		webInstitutionService.LoginInstitution(w, r, tokenAuth)
	})
	r.Post("/login/institution", func(w http.ResponseWriter, r *http.Request) {
		webInstitutionService.LoginInstitution(w, r, tokenAuth)
	})

	// Rotas Públicas de Consulta de Eventos (Listagem Pública)
	r.Get("/events", webEventService.GetEvents)
	r.Get("/event", webEventService.GetEvents)
	r.Get("/event/search", webEventService.GetEvents)
	r.Get("/event/{event_id}", webEventService.GetEventByID)
	r.Get("/event/search/{event}", webEventService.GetEventByName)
	r.Get("/event/institution/{institution_id}", webEventService.GetEventByOwnerID)

	// Rotas Públicas de Consulta de Instituições
	r.Get("/institution/{institution_id}", webInstitutionService.GetInstitutionByID)

	// Rotas Protegidas (Requer JWT)
	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(tokenAuth))
		protected.Use(jwtauth.Authenticator)

		// Criação de Instituição: apenas instituição do tipo admin pode criar
		protected.Post("/institution", webInstitutionService.CreateInstitution)
		protected.Post("/create/institution", webInstitutionService.CreateInstitution)

		// Gestão de Eventos (Instituição proprietária ou admin)
		protected.Post("/event", webEventService.CreateEvent)
		protected.Put("/event/{event_id}", webEventService.EditEvent)
		protected.Delete("/event/{event_id}", webEventService.DeleteEvent)

		// Gestão de Salas no Evento
		protected.Post("/event/{event_id}/room", webEventService.AddRoomInEvent)
		protected.Get("/event/{event_id}/room", webEventService.GetRoomsByEventID)
		protected.Put("/event/{event_id}/room", webEventService.EditRoom)
		protected.Delete("/event/{event_id}/room", webEventService.DeleteRoom)
	})

	return r
}
