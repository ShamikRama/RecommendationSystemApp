package converter

import (
	model2 "User_Service/internal/repository/model"
	"User_Service/internal/service/model"
)

func FromUpdateToUser(updateInfo model.UpdateInfoUser) model2.User {
	return model2.User{
		Email:    updateInfo.Email,
		Password: updateInfo.Password,
	}
}
