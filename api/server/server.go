package server

import (
	"database/sql"
	"fmt"
	fileprocessing "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/processingfile"
	"net/http"

	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/middleware"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/user"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"

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
		logger.Infof("Connected to db %v", cfg.Database.MySQL.DBName)
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
	userRouter.Use(middleware.LoggerMW, middleware.APIKeyMW, middleware.UserMW)
	userRouter.Post("/", userServer.CreateUserAPI())
	s.Router.Mount("/lfp/users", userRouter)

	// 3. File Award Point API
	fapServer := fileawardpoint.InitFileAwardPointServer(s.db)
	fapRouter := chi.NewRouter()
	fapRouter.Use(middleware.LoggerMW, middleware.APIKeyMW, middleware.UserMW)
	fapRouter.Post("/getListOrDetail", fapServer.GetDetailAPI())
	fapRouter.Get("/getList", fapServer.GetListAPI())
	fapRouter.Post("/create", fapServer.CreateFileAwardPointAPI())
	s.Router.Mount("/lfp/v1/fileAwardPoint", fapRouter)

	// 4. File Processing API
	fpServer := fileprocessing.InitFileProcessingServer(s.db)
	fpRouter := chi.NewRouter()
	fapRouter.Use(middleware.LoggerMW, middleware.APIKeyMW, middleware.UserMW)
	fapRouter.Get("/getList", fpServer.GetFileProcessHistoryAPI())
	s.Router.Mount("/fps/v1", fpRouter)
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
