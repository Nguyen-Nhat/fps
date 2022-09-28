package server

import (
	"database/sql"
	"errors"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/user"
	"net/http"

	"github.com/go-chi/render"
)

type UserServer struct {
	service *user.ServiceImpl
}

// InitUserServer ...
func InitUserServer(db *sql.DB) *UserServer {
	repo := user.NewRepo(db)
	service := user.NewService(repo)
	return &UserServer{
		service: service,
	}
}

func (s *UserServer) CreateUserAPI() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Bind data
		data := &createUserRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// 2. Map request to internal DTO
		req := &user.CreateUserRequestDTO{
			Name: data.Name,
		}

		// 3. Call function of Service
		u, err := s.service.CreateUser(r.Context(), req)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err)) // TODO @dung.nx correct error status
		}

		// 4. Map result to response
		res := newUserResponse(u)
		render.Render(w, r, res)
	}
}

type createUserRequest struct {
	Name        string
	Email       string
	Password    string `json:"password"`
	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *createUserRequest) Bind(r *http.Request) error {
	// a.User is nil if no User fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if len(a.Name) == 0 {
		return errors.New("missing required name")
	}

	a.ProtectedID = ""
	return nil
}

// userResponse is the response payload for the User data model.
//
// In the userResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type userResponse struct {
	Name  string
	Email string
}

func newUserResponse(user *user.User) *userResponse {
	return &userResponse{
		Name:  user.Name,
		Email: user.Email,
	}
}

func (ur *userResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
