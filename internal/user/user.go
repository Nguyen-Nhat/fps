package user

import "git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"

type User struct {
	ent.User
}

// mapper ...

func toUserFromCreateDTO(user *CreateUserRequestDTO) *User {
	u := ent.User{
		Name:  user.Name,
		Email: user.Email,
	}

	return &User{u}
}
