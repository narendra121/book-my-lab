package svcs

import (
	"errors"
	"fmt"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
)

type UserSvc struct {
	AppCfg *config.AppConfig
}

func NewUserSvc(cfg *config.AppConfig) *UserSvc {
	return &UserSvc{AppCfg: cfg}
}
func (u *UserSvc) CreateUser(user ...*model.User) error {
	usr := dao.User
	return usr.Save(user...)
}
func (u *UserSvc) UpdateUser(userName string, user *model.User) error {
	usr := dao.User
	_, err := u.GetUserWithDelFlag(userName)
	if err != nil {
		return err
	}
	_, err = usr.Where(usr.Deleted.Is(false)).Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).
		Updates(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserSvc) GetUserWithDelFlag(userName string) (*model.User, error) {
	usr := dao.User
	user, err := usr.Where(usr.Deleted.Is(false)).Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserSvc) GetUserWithoutDelFlag(userName string) (*model.User, error) {
	usr := dao.User
	user, err := usr.Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserSvc) GettAllUsers(page, limit int) ([]*model.User, error) {
	offset := (page - 1) * 10
	usr := dao.User
	users, _, err := usr.Where(usr.Deleted.Is(false)).FindByPage(offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserSvc) DelUser(userName string) error {
	usr := dao.User
	user, err := u.GetUserWithDelFlag(userName)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user already deleted")
	}
	_, err = usr.Where(usr.Deleted.Is(false)).Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).
		Updates(&model.User{Deleted: true})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserSvc) UpdateRefreshToken(userName, refreshToken string) error {
	usr := dao.User
	if _, err := usr.Where(usr.Deleted.Is(false)).Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).Select(usr.RefreshToken).Updates(&model.User{RefreshToken: refreshToken}); err != nil {
		return fmt.Errorf("failed to logged out, error:%v", err)
	}
	return nil
}
func (u *UserSvc) UpdateDelFlag(userName string, delFlag bool) error {
	usr := dao.User
	if _, err := usr.Where(usr.Email.Eq(userName)).Or(usr.Phone.Eq(userName)).Select(usr.Deleted).Updates(&model.User{Deleted: delFlag}); err != nil {
		return fmt.Errorf("failed to logged out, error:%v", err)
	}
	return nil
}
