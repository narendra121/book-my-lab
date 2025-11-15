package visits

import (
	"errors"
	"net/http"

	"booking.com/internal/dto"
	"booking.com/internal/svcs"
	"booking.com/internal/utils"
	"booking.com/pkg/constants"
	"github.com/gin-gonic/gin"
)

type VisitsHandler struct {
	VisitsSvc   *svcs.VisitsSvc
	UsrSvc      *svcs.UserSvc
	PropertySvc *svcs.PropertySvc
}

func NewVisitsHandler(visitsSvc *svcs.VisitsSvc) *VisitsHandler {
	return &VisitsHandler{VisitsSvc: visitsSvc}
}

func (u *VisitsHandler) ScheduleVisit(c *gin.Context) {
	var scheduleReq *dto.ScheduleReq
	if err := c.ShouldBindBodyWithJSON(&scheduleReq); err != nil {
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

}
func (u *VisitsHandler) UpdateVisit(c *gin.Context) {
}
func (u *VisitsHandler) FilterVisits(c *gin.Context) {

}
func (u *VisitsHandler) DeleteVisit(*gin.Context) {

}
