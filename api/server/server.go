package server

import (
	"context"
	"database/sql"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileprocessing"
	fps "git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"github.com/robfig/cron/v3"
	"net/http"

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
		// config db connection
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)

		debugDBConfig := cfg.Database.Debug
		if debugDBConfig.Enable {
			initJobPingDB(db, debugDBConfig)
		}

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

	// 4. File Processing API
	fpServer := fileprocessing.InitFileProcessingServer(s.db)
	fpRouter := chi.NewRouter()
	fpRouter.Use(middleware.LoggerMW, middleware.APIKeyMW, middleware.UserMW)
	fpRouter.Get("/getListProcessFiles", fpServer.GetFileProcessHistoryAPI())
	fpRouter.Post("/createProcessFile", fpServer.CreateProcessByFileAPI())
	s.Router.Mount("/v1", fpRouter)
	s.Router.Mount("/lfp/v1", fpRouter)
}

func (s *Server) Serve(cfg config.ServerListen) error {
	logger.Infof("Server is starting in port %v", cfg.Port)
	address := fmt.Sprintf(":%v", cfg.Port)
	return http.ListenAndServe(address, s.Router)
}

// initJobAccessDB ... job access DB each 1 minutes, we use this job for checking DB avoid loose connection DB
// ... will remove
func initJobPingDB(db *sql.DB, cfg config.DebugDBConfig) {
	fpRepo := fps.NewRepo(db)

	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	jobName := "Job ping DB for debugging"
	id, err := c.AddFunc(cfg.PingCron, func() {
		fpRepo.PingDB(context.Background(), 1)
	})
	if err != nil {
		logger.Errorf("Init Job %v failed: %v", jobName, err)
	}
	logger.Infof("Init %v Success with cron=\"%v\", ID=%v", jobName, cfg.PingCron, id)

	c.Start()
}
