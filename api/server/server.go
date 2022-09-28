package server

import (
	"database/sql"
	"fmt"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xo/dburl"
)

type Server struct {
	Router *chi.Mux
	db     *sql.DB
}

type Option func(*Server)

func NewServer(cfg config.Config, opts ...Option) (*Server, error) {

	srv := &Server{
		Router: chi.NewRouter(),
	}
	for _, opt := range opts {
		opt(srv)
	}

	// init db if necessary
	dbConf := cfg.Database.MySQL
	if srv.db == nil {
		db, err := dburl.Open(dbConf.DatabaseURI())
		if err != nil {
			return nil, fmt.Errorf("failed open DB: %w", err)
		}
		srv.db = db
	}

	srv.initRoutes()
	return srv, nil
}

func (s *Server) initRoutes() {
	// 1. System API
	healthzRouter := chi.NewRouter()
	healthzRouter.Get("/ready", ready)
	healthzRouter.Get("/liveness", liveness)
	s.Router.Mount("/healthz", healthzRouter)

	// 2. User API
	userServer := InitUserServer(s.db)
	userRouter := chi.NewRouter()
	userRouter.Use(middleware.Logger, apiKeyMiddleware)
	userRouter.Post("/", userServer.CreateUserAPI())
	s.Router.Mount("/lfp/users", userRouter)

	// 3. Other APIs
	// ...
}

func (s *Server) Serve(cfg config.ServerListen) error {
	fmt.Printf("Server is starting in port %v\n", cfg.Port)
	address := fmt.Sprintf(":%v", cfg.Port)
	return http.ListenAndServe(address, s.Router)
}

func WithDB(db *sql.DB) Option {
	return func(s *Server) {
		s.db = db
	}
}

// apiKeyMiddleware ...
func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		fmt.Printf("API KEY = %v\n", apiKey)
		next.ServeHTTP(w, r)
	})
}

// todo middleware for get Author
// ...
