package converter

import (
	serv "github.com/moremoneymod/auth/internal/model"
	repo "github.com/moremoneymod/auth/internal/repository/user/model"
)

func ToUserFromRepo(repoUser *repo.User) *serv.User {
	return &serv.User{
		Username: repoUser.Username,
		Password: repoUser.Password,
		Role:     repoUser.Role,
	}
}
