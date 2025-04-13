package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserUUID struct {
	ID         string `gorm:"column:id;type:varchar(27);primaryKey" json:"id"`
	UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserUUID) TableName() string {
	return "user_uuid"
}

func InputToUUID(userCreate UserInput) *UserUUID {
	var user UserUUID
	copier.Copy(&user, &userCreate)
	return &user
}
