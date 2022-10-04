package user

import "git.teko.vn/loyalty-system/loyalty-file-processing/internal/user"

func toUserResponseByUserEntity(user *user.User) *CreateUserResponse {
	return &CreateUserResponse{
		Name:  user.Name,
		Email: user.Email,
	}
}
