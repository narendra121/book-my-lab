package svcs

import (
	"fmt"
	"time"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
	"booking.com/pkg/utils"
)

type UserSvc struct {
	AppCfg *config.AppConfig
}

func NewUserSvc(cfg *config.AppConfig) *UserSvc {
	return &UserSvc{AppCfg: cfg}
}

func (u *UserSvc) CreateUser(user *dto.CreateUser) error {
	salt := utils.GetUUID()
	hashedPass, err := utils.HashPassword(user.Password + salt)
	if err != nil {
		return fmt.Errorf("failed to store user data, error: %v", err)
	}
	err = dao.Q.User.Create([]*model.User{{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Phone:        user.Phone,
		PasswordHash: hashedPass,
		Salt:         salt,
		Address:      user.Address,
		UpdatedAt:    time.Now(),
	}}...)
	if err != nil {
		return fmt.Errorf("failed to store user data, error: %v", err)
	}
	return nil
}

func (u *UserSvc) UpdateUser(user *dto.UpdateUser) error {
	usr := dao.User
	_, err := usr.Where(usr.Email.Eq(user.UserName)).Or(usr.Phone.Eq(user.UserName)).
		Updates(user)
	if err != nil {
		return fmt.Errorf("failed to update user data, error: %v", err)
	}
	return nil
}
