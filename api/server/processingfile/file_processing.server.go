package fileprocessing

import (
	"database/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"net/http"
)

type (
	IServer interface {
	}

	// Server ...
	Server struct {
		service *fileprocessing.ServiceImpl
	}
)

var _ IServer = &Server{}

func (s *Server) GetFileProcessHistoryAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// InitFileProcessingServer ...
func InitFileProcessingServer(db *sql.DB) *Server {
	repo := fileprocessing.NewRepo(db)
	service := fileprocessing.NewService(repo)
	return &Server{
		service: service,
	}
}
