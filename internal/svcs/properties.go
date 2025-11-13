package svcs

import (
	"context"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
)

type PropertySvc struct {
	AppCfg *config.AppConfig
}

func NewPropertySvc(cfg *config.AppConfig) *PropertySvc {
	return &PropertySvc{AppCfg: cfg}
}
func (p *PropertySvc) AddProperties(userName string, properties ...dto.AddPropertyReq) error {
	pr := dao.Property
	daoProperties := make([]*model.Property, 0)
	for _, property := range properties {
		daoProperty := &model.Property{
			PartnerUsername: userName,
			Title:           property.Title,
			Description:     property.Description,
			PropertyType:    property.PropertyType,
			Bedrooms:        property.Bedrooms,
			Bathrooms:       property.Bathrooms,
			AreaSqft:        property.AreaSqft,
			Price:           property.Price,
			City:            property.City,
			State:           property.State,
			Address:         property.Address,
		}
		daoProperties = append(daoProperties, daoProperty)
	}
	return pr.Create(daoProperties...)
}

func (p *PropertySvc) UpdateProperty(userName string, property dto.UpdatePropertyReq) error {
	_, err := p.GetPropertyByID(property.ID)
	if err != nil {
		return err
	}
	daoProperty := &model.Property{
		PartnerUsername: userName,
		Title:           property.Title,
		Description:     property.Description,
		PropertyType:    property.PropertyType,
		Bedrooms:        property.Bedrooms,
		Bathrooms:       property.Bathrooms,
		AreaSqft:        property.AreaSqft,
		Price:           property.Price,
		City:            property.City,
		State:           property.State,
		Address:         property.Address,
	}
	pr := dao.Property
	return pr.Save(daoProperty)
}

func (p *PropertySvc) GetPropertyByID(id int64) (*model.Property, error) {
	pr := dao.Property
	property, err := pr.Where(pr.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return property, nil
}
func (p *PropertySvc) GetPropertiesByUserName(userName string) ([]*model.Property, error) {
	pr := dao.Property
	properties, err := pr.Where(pr.PartnerUsername.Eq(userName)).Find()
	if err != nil {
		return nil, err
	}
	return properties, nil
}
func (s *PropertySvc) GetFilteredProperties(userName string, excludeSelf bool, city, state, status string) ([]*model.Property, error) {
	q := dao.Property.WithContext(context.Background())

	if excludeSelf {
		q = q.Where(dao.Property.PartnerUsername.Neq(userName))
	} else {
		q = q.Where(dao.Property.PartnerUsername.Eq(userName))
	}

	if city != "" {
		q = q.Where(dao.Property.City.Like("%" + city + "%"))
	}
	if state != "" {
		q = q.Where(dao.Property.State.Like("%" + state + "%"))
	}
	if status != "" {
		q = q.Where(dao.Property.Status.Eq(status))
	}

	return q.Find()
}
