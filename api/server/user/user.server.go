package user

import (
	"context"
	"database/sql"
	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	res "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/user"
	"net/http"

	"github.com/go-chi/render"
)

type (
	IUserServer interface {
		// APIs for Create User ----------------------------------------------------------------------------------------

		// CreateUserAPI is used for defining Router
		CreateUserAPI() func(http.ResponseWriter, *http.Request)
		// CreateUser is called in createUserAPI() and it handles logic of API
		CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)

		// APIs for <DO STH> User --------------------------------------------------------------------------------------

		// DoSthUserAPI() ...
		// ...
		// DoSthUserXXX()...
		// ...
	}

	// UserServer ...
	UserServer struct {
		service *user.ServiceImpl
	}
)

// InitUserServer ...
func InitUserServer(db *sql.DB) *UserServer {
	repo := user.NewRepo(db)
	service := user.NewService(repo)
	return &UserServer{
		service: service,
	}
}

// CreateUserAPI ...
func (s *UserServer) CreateUserAPI() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Bind data & validate input
		data := &CreateUserRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, error2.ErrRenderInvalidRequest(err))
			return
		}

		// 2. Handle request
		res, err := s.CreateUser(r.Context(), data)
		if err != nil {
			render.Render(w, r, error2.ToErrorResponse(err))
		}

		// 3. Render response
		render.Render(w, r, res)
	}
}

// CreateUser ...
func (s *UserServer) CreateUser(ctx context.Context, data *CreateUserRequest) (*res.BaseResponse[CreateUserResponse], error) {
	// 2. Map request to internal DTO
	req := &user.CreateUserRequestDTO{
		Name: data.Name,
	}

	// 3. Call function of Service
	u, err := s.service.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	// 4. Map result to response
	resp := toUserResponseByUserEntity(u)
	resp2 := res.ToResponse(resp)
	return resp2, nil
}
