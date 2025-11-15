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
	filterReq dto.PropertFilterReq,
	withDelFlag bool,
) ([]*model.Property, error) {

	ctx := context.Background()
	pr := dao.Property.WithContext(ctx)
	usr := dao.User.WithContext(ctx)

	pr = pr.LeftJoin(usr, dao.User.Username.EqCol(dao.Property.PartnerUsername))

	if filterReq.Id != 0 {
		pr = pr.Where(dao.Property.ID.Eq(filterReq.Id))
	}
	if filterReq.Title != "" {
		pr = pr.Where(dao.Property.Title.Like("%" + filterReq.Title + "%"))
	}
	if filterReq.PartnerName != "" {
		pr = pr.Where(dao.Property.PartnerUsername.Like("%" + filterReq.PartnerName + "%"))
	}

	if filterReq.From_Price != 0 && filterReq.To_Price != 0 {
		pr = pr.Where(dao.Property.Price.Between(filterReq.From_Price, filterReq.To_Price))
	} else if filterReq.From_Price != 0 {
		pr = pr.Where(dao.Property.Price.Between(0, filterReq.From_Price))
	} else if filterReq.To_Price != 0 {
		pr = pr.Where(dao.Property.Price.Between(0, filterReq.To_Price))
	}

	if filterReq.City != "" {
		pr = pr.Where(dao.Property.City.Like("%" + filterReq.City + "%"))
	}

	if filterReq.State != "" {
		pr = pr.Where(dao.Property.State.Like("%" + filterReq.State + "%"))
	}

	if filterReq.Status != "" {
		pr = pr.Where(dao.Property.Status.Eq(filterReq.Status))
	}

	if filterReq.ExcludeSelf {
		pr = pr.Where(dao.Property.PartnerUsername.Neq(userName))
	}
	if withDelFlag {
		pr = pr.Where(
			dao.Property.Deleted.Is(false),
			dao.User.Deleted.Is(false), // safe now because join is already applied
		)
	}

	return pr.Find()
}
func (p *PropertySvc) DeletePropertyByID(id int64, deleteFlag bool) error {
	pr := dao.Property.WithContext(context.Background())
	_, err := pr.Where(dao.Property.ID.Eq(id)).Select(dao.Property.Deleted).Updates(&model.Property{Deleted: deleteFlag})
	return err
}
