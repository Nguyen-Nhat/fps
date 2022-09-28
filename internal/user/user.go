package user

type User struct {
	// ID of the migrations.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Active holds the value of the "active" field.
	Active bool `json:"active,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
	// PhoneNumber holds the value of the "phone_number" field.
	PhoneNumber string `json:"phone_number,omitempty"`

	PasswordHash string `json:"-"`
}
