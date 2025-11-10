package users

import (
	"fmt"
	"net/http"
	"time"

	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

func Add(w http.ResponseWriter, r *http.Request) {
	var user *dto.CreateUser
	err := utils.ParseHttpRequest(r, &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user request to store, error: %v", err), http.StatusInternalServerError)
		return
	}
	salt := utils.GetUUID()
	hashedPass, err := utils.HashPassword(user.Password + salt)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to store user data, error: %v", err), http.StatusInternalServerError)
		return
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
		http.Error(w, fmt.Sprintf("failed to store user data, error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeTextPlain)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user are created"))
}

func Put(w http.ResponseWriter, r *http.Request) {
	var user *dto.UpdateUser
	err := utils.ParseHttpRequest(r, &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user request to store, error: %v", err), http.StatusInternalServerError)
		return
	}
	u := dao.User
	_, err = u.Where(u.Email.Eq(user.UserName)).Or(u.Phone.Eq(user.UserName)).
		Updates(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update user data, error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeTextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user are updated"))
}
