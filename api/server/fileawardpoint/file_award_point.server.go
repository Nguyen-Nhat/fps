package fileawardpoint

import (
	"context"
	"database/sql"
	"fmt"
	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/error"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"github.com/go-chi/render"
	"net/http"
)

type (
	IServer interface {
		// APIs for Create User ----------------------------------------------------------------------------------------

		// GetDetailAPI is used for defining Router
		GetDetailAPI() func(http.ResponseWriter, *http.Request)
		// GetDetail is called in GetDetailAPI() and it handles logic of API
		GetDetail(context.Context, *GetFileAwardPointDetailRequest) (*GetFileAwardPointDetailResponse, error)

		// APIs for <DO STH> User --------------------------------------------------------------------------------------

		// DoSthUserAPI() ...
		// ...
		// DoSthUserXXX()...
		// ...
	}

	// Server ...
	Server struct {
		service *fileawardpoint.ServiceImpl
	}
)

func (s Server) GetDetailAPI() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Bind data & validate input
		data := &GetFileAwardPointDetailRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, error2.ErrInvalidRequest(err))
			return
		}
		fmt.Printf("data = %v", data)

		// 2. Handle request
		res, err := s.GetDetail(r.Context(), data)
		if err != nil {
			render.Render(w, r, error2.ErrInvalidRequest(err)) // TODO @dung.nx correct error status
		}

		// 3. Render response
		render.Render(w, r, res)
	}
}

func (s Server) GetDetail(ctx context.Context, request *GetFileAwardPointDetailRequest) (*GetFileAwardPointDetailResponse, error) {
	// 1. Map the request to internal DTO if input for Service too complex
	req := &fileawardpoint.GetFileAwardPointDetailDTO{
		Id: request.Id,
	}

	// 2. Call function of Service
	u, err := s.service.GetFileAwardPoint(ctx, req)
	if err != nil {
		return nil, err
	}

	// 4. Map result to response
	res := toFapDetailResponseByEntity(u)
	return res, nil
}

var _ IServer = &Server{}

// InitFileAwardPointServer ...
func InitFileAwardPointServer(db *sql.DB) *Server {
	repo := fileawardpoint.NewRepo(db)
	service := fileawardpoint.NewService(repo)
	return &Server{
		service: service,
	}
}
