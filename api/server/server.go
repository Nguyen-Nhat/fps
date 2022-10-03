package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/user"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"

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
			logger.Errorf("Fail to open db, got: %v", err)
			return nil, fmt.Errorf("failed open DB: %w", err)
		}
		logger.Info("Connected to db")
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
	userServer := user.InitUserServer(s.db)
	userRouter := chi.NewRouter()
	userRouter.Use(loggerMiddleware, apiKeyMiddleware)
	userRouter.Post("/", userServer.CreateUserAPI())
	s.Router.Mount("/lfp/users", userRouter)

	// 3. File Award Point API
	fapServer := fileawardpoint.InitFileAwardPointServer(s.db)
	fapRouter := chi.NewRouter()
	fapRouter.Use(loggerMiddleware, apiKeyMiddleware)
	fapRouter.Post("/getListOrDetail", fapServer.GetDetailAPI())
	s.Router.Mount("/lfp/fileAwardPoint", fapRouter)

	// 4. Other APIs
	// ...
}

func (s *Server) Serve(cfg config.ServerListen) error {
	logger.Infof("Server is starting in port %v", cfg.Port)
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

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		startReqTime := time.Now()
		defer func() {
			logger.Infof("%s %s%s%s %s %d %dB in %s",
				r.Method,
				r.URL.Scheme,
				r.Host,
				r.URL.Path,
				r.Proto,
				ww.Status(),
				ww.BytesWritten(),
				time.Since(startReqTime),
			)
		}()
		next.ServeHTTP(ww, r)
	})
}

// todo middleware for get Author
// ...
