package user

// CreateUserRequestDTO ...
type CreateUserRequestDTO struct {
	Name     string
	Email    string
	Password string
}

// UserDTO ...
type UserDTO struct {
	ID    int32
	Name  string
	Email string
}
