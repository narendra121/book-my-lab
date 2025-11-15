package dto

type AddPropertyReq struct {
	Title        string  `gorm:"column:title;type:character varying(200);not null" json:"title"`
	Description  string  `gorm:"column:description;type:text" json:"description"`
	PropertyType string  `gorm:"column:property_type;type:property_type_enum;not null" json:"property_type"`
	Bedrooms     int32   `gorm:"column:bedrooms;type:integer" json:"bedrooms"`
	Bathrooms    int32   `gorm:"column:bathrooms;type:integer" json:"bathrooms"`
	AreaSqft     float64 `gorm:"column:area_sqft;type:numeric(10,2)" json:"area_sqft"`
	Price        float64 `gorm:"column:price;type:numeric(12,2)" json:"price"`
	City         string  `gorm:"column:city;type:character varying(100)" json:"city"`
	State        string  `gorm:"column:state;type:character varying(100)" json:"state"`
	Address      string  `gorm:"column:address;type:text" json:"address"`
}
type UpdatePropertyReq struct {
	ID           int64   `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	Title        string  `gorm:"column:title;type:character varying(200);not null" json:"title"`
	Description  string  `gorm:"column:description;type:text" json:"description"`
	PropertyType string  `gorm:"column:property_type;type:property_type_enum;not null" json:"property_type"`
	Bedrooms     int32   `gorm:"column:bedrooms;type:integer" json:"bedrooms"`
	Bathrooms    int32   `gorm:"column:bathrooms;type:integer" json:"bathrooms"`
	AreaSqft     float64 `gorm:"column:area_sqft;type:numeric(10,2)" json:"area_sqft"`
	Price        float64 `gorm:"column:price;type:numeric(12,2)" json:"price"`
	City         string  `gorm:"column:city;type:character varying(100)" json:"city"`
	State        string  `gorm:"column:state;type:character varying(100)" json:"state"`
	Address      string  `gorm:"column:address;type:text" json:"address"`
}

type GetProperty struct {
	ID int64 `uri:"id" binding:"required"`
}
