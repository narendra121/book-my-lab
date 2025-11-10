package dto

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type CreateUser struct {
	FirstName string `gorm:"column:first_name;type:character varying(100);not null" json:"first_name"`
	LastName  string `gorm:"column:last_name;type:character varying(100);not null" json:"last_name"`
	Email     string `gorm:"column:email;type:character varying(100);not null" json:"email"`
	Phone     string `gorm:"column:phone;type:character varying(20)" json:"phone"`
	Password  string `gorm:"column:password;type:character varying(255)" json:"password"`
	Address   string `gorm:"column:address;type:character varying(255);not null" json:"address"`
}

type UpdateUser struct {
	UserName  string `json:"username"`
	FirstName string `gorm:"column:first_name;type:character varying(100);not null" json:"first_name"`
	LastName  string `gorm:"column:last_name;type:character varying(100);not null" json:"last_name"`
	Address   string `gorm:"column:address;type:character varying(255);not null" json:"address"`
}
