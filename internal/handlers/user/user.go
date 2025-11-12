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

func (u *UserHandler) UpdateUser(c *gin.Context) {
	var updateReq *dto.UpdateUser
	if err := c.ShouldBindBodyWithJSON(&updateReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	reqUserName, ok := c.Get(constants.CurrentUserName)
	if !ok || reqUserName == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	userName := reqUserName.(string)
	if userName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	err := u.UserSvc.UpdateUser(userName, &model.User{FirstName: updateReq.FirstName, LastName: updateReq.LastName, Address: updateReq.Address})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusAccepted, utils.WriteAppResponse("user updated", nil, nil))
}
func (u *UserHandler) GetProfile(c *gin.Context) {
	reqUserName, ok := c.Get(constants.CurrentUserName)
	if !ok || reqUserName == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	userName := reqUserName.(string)
	if userName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	user, err := u.UserSvc.GetUserWithDelFlag(userName)
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

func (u *UserHandler) ListUsers(c *gin.Context) {
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
	err := u.UserSvc.UpdateUser(roleReq.UserName, &model.User{Role: roleReq.Role})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusAccepted, utils.WriteAppResponse("role updated", nil, nil))
}
func (u *UserHandler) DeleteUser(c *gin.Context) {
	reqUserName, ok := c.Get(constants.CurrentUserName)
	if !ok || reqUserName == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	userName := reqUserName.(string)
	if userName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid request"), nil))
		return
	}
	err := u.UserSvc.DelUser(userName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusAccepted, utils.WriteAppResponse("user deleted", nil, nil))
}
