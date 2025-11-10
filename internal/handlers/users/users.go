package users

import (
	"fmt"
	"net/http"

	"booking.com/internal/dto"
	"booking.com/internal/svcs"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

type UserHandler struct {
	UserSvc *svcs.UserSvc
}

func NewUserHandler(userSvc *svcs.UserSvc) *UserHandler {
	return &UserHandler{UserSvc: userSvc}
}

func (u *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	var user *dto.CreateUser
	err := utils.ParseHttpRequest(r, &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user request to store, error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := u.UserSvc.CreateUser(user); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeTextPlain)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user are created"))
}

func (u *UserHandler) Put(w http.ResponseWriter, r *http.Request) {
	var user *dto.UpdateUser
	err := utils.ParseHttpRequest(r, &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid user request to store, error: %v", err), http.StatusInternalServerError)
		return
	}
	err = u.UserSvc.UpdateUser(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeTextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user are updated"))
}
