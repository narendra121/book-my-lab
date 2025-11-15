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

func (p *PropertySvc) UpdateProperty(property dto.UpdatePropertyReq) error {
	_, err := p.GetPropertyByID(property.ID, true)
	if err != nil {
		return err
	}

	daoProperty := &model.Property{
		// ID:           prop.ID,
		Title:        property.Title,
		Description:  property.Description,
		PropertyType: property.PropertyType,
		Bedrooms:     property.Bedrooms,
		Bathrooms:    property.Bathrooms,
		AreaSqft:     property.AreaSqft,
		Price:        property.Price,
		City:         property.City,
		State:        property.State,
		Address:      property.Address,
	}

	ctx := context.Background()
	pr := dao.Property.WithContext(ctx)

	_, err = pr.Where(dao.Property.ID.Eq(property.ID), dao.Property.Deleted.Is(false)).
		Updates(daoProperty)
	return err
}
func (p *PropertySvc) GetPropertyByID(id int64, withDelFlag bool) (*model.Property, error) {
	pr := dao.Property.WithContext(context.Background())
	usr := dao.User
	pr = pr.Join(usr, usr.Username.EqCol(dao.Property.PartnerUsername))
	if withDelFlag {
		pr = pr.Where(dao.Property.Deleted.Is(false), usr.Deleted.Is(false))
	}
	property, err := pr.First()
	if err != nil {
		return nil, err
	}
	return property, nil
}

func (p *PropertySvc) GetPropertiesByUserName(userName string, withDelFlag bool) ([]*model.Property, error) {
	pr := dao.Property.WithContext(context.Background())
	usr := dao.User
	pr = pr.Join(usr, usr.Username.EqCol(dao.Property.PartnerUsername))
	if withDelFlag {
		pr = pr.Where(dao.Property.Deleted.Is(false), usr.Deleted.Is(false))
	}
	properties, err := pr.Find()
	if err != nil {
		return nil, err
	}
	return properties, nil
}
func (s *PropertySvc) GetFilteredProperties(
	userName string,
	excludeSelf bool,
	city, state, status string,
	withDelFlag bool,
) ([]*model.Property, error) {

	ctx := context.Background()
	pr := dao.Property.WithContext(ctx)
	usr := dao.User.WithContext(ctx)

	pr = pr.LeftJoin(usr, dao.User.Username.EqCol(dao.Property.PartnerUsername))

	if excludeSelf {
		pr = pr.Where(dao.Property.PartnerUsername.Neq(userName))
	}
	if withDelFlag {
		pr = pr.Where(
			dao.Property.Deleted.Is(false),
			dao.User.Deleted.Is(false), // safe now because join is already applied
		)
	}

	if city != "" {
		pr = pr.Where(dao.Property.City.Like("%" + city + "%"))
	}

	if state != "" {
		pr = pr.Where(dao.Property.State.Like("%" + state + "%"))
	}

	if status != "" {
		pr = pr.Where(dao.Property.Status.Eq(status))
	}

	return pr.Find()
}
