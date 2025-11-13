package svcs

import (
	"context"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/utils"
)

type UserSvc struct {
	AppCfg *config.AppConfig
}

func NewUserSvc(cfg *config.AppConfig) *UserSvc {
	return &UserSvc{AppCfg: cfg}
}
func (u *UserSvc) CreateUser(user *model.User) error {
	usr := dao.User
	return usr.Save(user)
}
func (u *UserSvc) UpdateUser(userName string, user *model.User) error {
	usr := dao.User
	_, err := u.FilterUsers(userName, "", "", true)
	if err != nil {
		return err
	}
	_, err = usr.Where(usr.Username.Eq(userName)).
		Updates(user)
	if err != nil {
		return err
	}
	return nil
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
	user, err := u.FilterUsers(userName, "", "", true)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.ErrUserAlreadyDeleted
	}
	_, err = usr.Where(usr.Username.Eq(userName)).
		Updates(&model.User{Deleted: true})
	if err != nil {
		return err
	}
	return nil
}

func (u *UserSvc) UpdateRefreshToken(userName, refreshToken string) error {
	usr := dao.User
	if _, err := usr.Where(usr.Deleted.Is(false), usr.Username.Eq(userName)).Select(usr.RefreshToken).Updates(&model.User{RefreshToken: refreshToken}); err != nil {
		return err
	}
	return nil
}
func (u *UserSvc) UpdateDelFlag(userName string, delFlag bool) error {
	usr := dao.User
	if _, err := usr.Where(usr.Deleted.Is(true), usr.Username.Eq(userName)).Select(usr.Deleted).Updates(&model.User{Deleted: delFlag}); err != nil {
		return err
	}
	return nil
}
func (u *UserSvc) FilterUsers(userName, email, phone string, useDelFlag bool) ([]*model.User, error) {
	q := dao.User.WithContext(context.Background())

	if useDelFlag {
		q = q.Where(dao.User.Deleted.Is(false))
	}
	if userName != "" {
		q = q.Where(dao.User.Username.Eq(userName))
	}
	if email != "" {
		q = q.Where(dao.User.Email.Eq(email))
	}
	if phone != "" {
		q = q.Where(dao.User.Phone.Eq(phone))
	}
	return q.Find()
}
func (u *UserSvc) GetUserWithEmailOrPhone(email, phone string, useDelFlag bool) (*model.User, error) {
	q := dao.User.WithContext(context.Background())
	if useDelFlag {
		q = q.Where(dao.User.Deleted.Is(false))
	}
	q = q.Where(dao.User.Email.Eq(email)).
		Or(dao.User.Phone.Eq(phone))
	return q.First()
}
func (u *UserSvc) GetUserWithEmailAndPhone(email, phone string, useDelFlag bool) (*model.User, error) {
	q := dao.User.WithContext(context.Background())
	if useDelFlag {
		q = q.Where(dao.User.Deleted.Is(false))
	}
	q = q.Where(dao.User.Email.Eq(email), dao.User.Phone.Eq(phone))
	return q.First()
}

func (u *UserSvc) GetUserByUserName(userName string, useDelFlag bool) (*model.User, error) {
	q := dao.User.WithContext(context.Background())
	if useDelFlag {
		q = q.Where(dao.User.Deleted.Is(false))
	}
	q = q.Where(dao.User.Username.Eq(userName))
	return q.First()
}
