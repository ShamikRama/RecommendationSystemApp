package converter

import (
	model2 "User_Service/internal/repository/model"
	"User_Service/internal/service/model"
)

func FromLoginToUser(login model.Login) model2.User {
	return model2.User{
		Email:    login.Email,
		Password: login.Password,
	}
}
