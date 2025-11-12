package svcs

import (
	"fmt"
	"time"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
	"booking.com/internal/utils"
	"booking.com/pkg/constants"
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
		Role:         constants.UserRole,
		UpdatedAt:    time.Now(),
	}}...)
	if err != nil {
		return fmt.Errorf("failed to store user data, error: %v", err)
	}
	return nil
}

func (u *UserSvc) UpdateUser(userName string, user *model.User) error {
	usr := dao.User
	_, err := usr.Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).Where(usr.Deleted.Is(false)).
		Updates(user)
	if err != nil {
		return fmt.Errorf("failed to update user data, error: %v", err)
	}
	return nil
}

func (u *UserSvc) GetUser(userName string) (*model.User, error) {
	usr := dao.User
	user, err := usr.Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).Where(usr.Deleted.Is(false)).First()
	if err != nil {
		return nil, fmt.Errorf("failed to get user data, error: %v", err)
	}
	return user, nil
}

func (u *UserSvc) GettAllUsers(page, limit int) ([]*model.User, error) {
	offset := (page - 1) * 10
	usr := dao.User
	users, _, err := usr.Where(usr.Deleted.Is(false)).FindByPage(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data, error: %v", err)
	}
	return users, nil
}

func (u *UserSvc) DelUser(userName string) error {
	usr := dao.User
	_, err := usr.Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).Where(usr.Deleted.Is(false)).
		Updates(&model.User{Deleted: true})
	if err != nil {
		return fmt.Errorf("failed to delete user data, error: %v", err)
	}
	return nil
}
