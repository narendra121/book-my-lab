package user

import (
	"errors"
	"net/http"
	"strconv"

	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
	"booking.com/internal/svcs"
	"booking.com/internal/utils"
	"booking.com/pkg/constants"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserSvc *svcs.UserSvc
}

func NewUserHandler(userSvc *svcs.UserSvc) *UserHandler {
	return &UserHandler{UserSvc: userSvc}
}

func (u *UserHandler) Add(c *gin.Context) {
	var user dto.CreateUser
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	if err := u.UserSvc.CreateUser(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusCreated, utils.WriteAppResponse("user created", nil, nil))
}

func (u *UserHandler) Put(c *gin.Context) {
	var user *dto.UpdateUser
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	err := u.UserSvc.UpdateUser(user.UserName, &model.User{FirstName: user.FirstName, LastName: user.LastName, Address: user.Address})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusAccepted, utils.WriteAppResponse("user updated", nil, nil))
}
func (u *UserHandler) Get(c *gin.Context) {
	var profileReq dto.UserProfile
	if err := c.ShouldBindUri(&profileReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	if profileReq.UserName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("username not provided"), nil))
		return
	}
	user, err := u.UserSvc.GetUser(profileReq.UserName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	if user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.WriteAppResponse("", errors.New("user not found"), nil))
		return
	}
	c.JSON(http.StatusFound, utils.WriteAppResponse("", nil, user))
}

func (u *UserHandler) GetAll(c *gin.Context) {
	if role, ok := c.Get(constants.Role); !ok || role != constants.AdminRole {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("user don't have access to get all users", nil, nil))
		return
	}
	page, limit := 1, 10
	if p, err := strconv.Atoi(c.Query("page")); err != nil && p != 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.Query("limit")); err != nil && l != 0 {
		limit = l
	}
	users, err := u.UserSvc.GettAllUsers(page, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	if len(users) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.WriteAppResponse("", errors.New("users data not found"), nil))
		return
	}
	c.JSON(http.StatusFound, utils.WriteAppResponse("", nil, users))
}

func (u *UserHandler) UpdateRole(c *gin.Context) {
	if role, ok := c.Get(constants.Role); !ok || role != constants.AdminRole {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("user don't have access to modify role", nil, nil))
		return
	}
	var roleReq *dto.UserRoleReq
	if err := c.ShouldBindBodyWithJSON(&roleReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	user, err := u.UserSvc.GetUser(roleReq.UserName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	if user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.WriteAppResponse("", errors.New("user not found"), nil))
		return
	}
	err = u.UserSvc.UpdateUser(roleReq.UserName, &model.User{Role: roleReq.Role})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusAccepted, utils.WriteAppResponse("role updated", nil, nil))
}
func (u *UserHandler) Delete(c *gin.Context) {
	var delReq dto.UserProfile
	if err := c.ShouldBindUri(&delReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	user, err := u.UserSvc.GetUser(delReq.UserName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	if user == nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.WriteAppResponse("user already deleted", nil, nil))
		return
	}
	err = u.UserSvc.DelUser(delReq.UserName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusAccepted, utils.WriteAppResponse("user deleted", nil, nil))
}
