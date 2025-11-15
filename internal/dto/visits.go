package dto

import "time"

type ScheduleReq struct {
	PropertyID    int64 `json:"property_id"`
	BuyerUsername string
	ScheduledTime time.Time `json:"scheduled_time"`
	BuyerNote     string    `json:"buyer_note"`
}

type VisitFilterReq struct {
	Status          string
	PropertyID      int64
	BuyersUserName  string
	PartnerUserName string
}
