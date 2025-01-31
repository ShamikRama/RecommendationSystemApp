package user

import (
	"User_Service/internal/converter"
	"User_Service/internal/erorrs"
	"User_Service/internal/logger"
	"User_Service/internal/repository"
	"User_Service/internal/service/model"
	"User_Service/internal/utils"
	"User_Service/pkg/kafka"
	"context"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

type ServiceUser struct {
	repo   repository.User
	logger *logger.Logger
	kafka  kafka.Kafka
}

func NewUserService(repo repository.User, logger *logger.Logger, kafka kafka.Kafka) *ServiceUser {
	return &ServiceUser{
		repo:   repo,
		logger: logger,
		kafka:  kafka,
	}
}

func (r *ServiceUser) UpdateUser(ctx context.Context, userID int, updateInfo model.UpdateInfoUser) (int, error) {
	if len(updateInfo.Email) <= 0 || len(updateInfo.Password) <= 0 {
		r.logger.Error("Invalid input",
			zap.String("operation", "service.User.UpdateUser"))
		return 0, erorrs.ErrEmptyInput
	}

	user := converter.FromUpdateToUser(updateInfo)
	user.Id = userID
	user.Password = utils.GeneratePasswordHash(user.Password)
	id, err := r.repo.UpdateUser(ctx, user)
	if err != nil {
		r.logger.Error("Error to update the user",
			zap.String("operation", "service.User.UpdateUser"),
			zap.Error(err))

		if errors.Is(err, erorrs.ErrNoRows) {
			return 0, erorrs.ErrUserNotFound
		}
		return 0, err
	}

	key := strconv.Itoa(id)

	msg := &kafka.KafkaMessage{
		Id:    id,
		Event: "user_update",
		Data: map[string]interface{}{
			"email": updateInfo.Email,
		},
	}

	err = r.kafka.SendMessage(key, msg)
	if err != nil {
		r.logger.Error("Error sending msg to kafka",
			zap.String("operation", "service.User.UpdateUser"),
			zap.Error(err))
		return 0, err
	}

	r.logger.Info("Success sending msg to kafka")

	return id, nil
}
