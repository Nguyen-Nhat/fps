package user

import (
	"errors"
	"net/http"
)

// Request DTO =========================================================================================================

// CreateUserRequest ...
type CreateUserRequest struct {
	Name        string
	Email       string
	Password    string `json:"password"`
	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *CreateUserRequest) Bind(_ *http.Request) error {
	// a.User is nil if no User fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if len(a.Name) == 0 {
		return errors.New("missing required name")
	}

	a.ProtectedID = ""
	return nil
}

// CreateUserResponse is the response payload for the User data model.
//
// In the userResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type CreateUserResponse struct {
	Name  string
	Email string
}

func (ur *CreateUserResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
