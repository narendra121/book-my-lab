package svcs

import (
	"context"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
	"booking.com/pkg/constants"
)

type VisitsSvc struct {
	AppCfg *config.AppConfig
}

func (v *VisitsSvc) ScheduleVisit(visitReq *dto.ScheduleReq, propertySvc *PropertySvc) error {
	_, err := propertySvc.GetPropertyByID(visitReq.PropertyID, true)
	if err != nil {
		return err
	}
	return dao.Visit.Save(&model.Visit{
		PropertyID:     visitReq.PropertyID,
		BuyerUsername:  visitReq.BuyerUsername,
		RescheduleTime: visitReq.ScheduledTime,
		Status:         constants.Pending,
		BuyerNote:      visitReq.BuyerNote,
	})
}
func (v *VisitsSvc) FilterVisits(filterReq *dto.VisitFilterReq, propertySvc *PropertySvc) {

	vist := dao.Visit.WithContext(context.Background())
	if filterReq.Status != "" {
		vist = vist.Where(dao.Visit.Status.Eq(filterReq.Status))
	}
	if filterReq.PropertyID != 0 {
		vist = vist.Where(dao.Visit.ID.Eq(filterReq.PropertyID))
	}
	if filterReq.BuyersUserName != "" {
		vist = vist.Where(dao.Visit.BuyerUsername.Eq(filterReq.BuyersUserName))
	}
	if filterReq.PartnerUserName != "" {
		propertySvc.GetPropertyByID(filterReq.PropertyID,true)
		vist = vist.Where(dao.Visit.BuyerUsername.Eq(filterReq.BuyersUserName))
	}
	// vist.Join()
	// Where(vist.Deleted.Is(false)).

}
