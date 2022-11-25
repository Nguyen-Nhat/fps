package fileprocessing

import (
	"database/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
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

// InitFileProcessingServer ...
func InitFileProcessingServer(db *sql.DB) *Server {
	repo := fileprocessing.NewRepo(db)
	service := fileprocessing.NewService(repo)
	return &Server{
		service: service,
	}
}
