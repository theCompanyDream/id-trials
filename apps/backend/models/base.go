package models

type UserBase struct {
	UserName   string    `gorm:"column:user_name;type:varchar(20);not null" json:"user_name"`
	FirstName  string    `gorm:"column:first_name;type:varchar(40);not null" json:"first_name"`
	LastName   string    `gorm:"column:last_name;type:varchar(40);not null" json:"last_name"`
	Email      string    `gorm:"column:email;type:varchar(40);not null;unique" json:"email"`
	Department *string   `gorm:"column:department;type:varchar(25)" json:"department"`
}