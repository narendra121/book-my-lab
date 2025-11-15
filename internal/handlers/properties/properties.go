package properties

import (
	"errors"
	"net/http"

	"booking.com/internal/dto"
	"booking.com/internal/svcs"
	"booking.com/internal/utils"
	"booking.com/pkg/constants"
	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	PropertySvc *svcs.PropertySvc
}

func NewPropertyHandler(propertySvc *svcs.PropertySvc) *PropertyHandler {
	return &PropertyHandler{PropertySvc: propertySvc}
}
func (p *PropertyHandler) AddProperties(c *gin.Context) {
	var propertiesReq []dto.AddPropertyReq
	if err := c.ShouldBindBodyWithJSON(&propertiesReq); err != nil {
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
	err := p.PropertySvc.AddProperties(userName, propertiesReq...)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusCreated, utils.WriteAppResponse("properties added", nil, nil))
}
func (p *PropertyHandler) UpdateProperty(c *gin.Context) {
	var propertiesReq dto.UpdatePropertyReq
	if err := c.ShouldBindBodyWithJSON(&propertiesReq); err != nil {
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
	err := p.PropertySvc.UpdateProperty(propertiesReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusCreated, utils.WriteAppResponse("property updated", nil, nil))
}

func (p *PropertyHandler) GetPropertyByID(c *gin.Context) {
	var updateReq *dto.GetProperty
	if err := c.ShouldBindUri(&updateReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	property, err := p.PropertySvc.GetPropertyByID(updateReq.ID, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusFound, utils.WriteAppResponse("", nil, property))
}

func (p *PropertyHandler) GetFilteredProperties(c *gin.Context) {
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
	excludeSelf := c.DefaultQuery("exclude_user", "false") == "true"
	city := c.Query("city")
	state := c.Query("state")
	status := c.Query("status")
	properties, err := p.PropertySvc.GetFilteredProperties(userName, excludeSelf, city, state, status, true)
	if err != nil || len(properties) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusFound, utils.WriteAppResponse("", nil, properties))
}

func (p *PropertyHandler) GetAllProperties(c *gin.Context) {
	properties, err := p.PropertySvc.GetFilteredProperties("", true, "", "", "", true)
	if err != nil || len(properties) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusFound, utils.WriteAppResponse("", nil, properties))
}
