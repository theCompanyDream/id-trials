package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserUlid struct {
	ID         string `gorm:"column:id;type:varchar(26);primaryKey" json:"id"`
	*UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserUlid) TableName() string {
	return "user_ulid"
}

func InputToDTO(userCreate UserInput) *UserUlid {
	var user UserUlid
	copier.Copy(&user, &userCreate)
	return &user
}
