package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserUUID struct {
	ID string `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	*UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserUUID) TableName() string {
	return "users_uuid"
}

func InputToUUID(userCreate UserInput) *UserUUID {
	var user UserUUID
	copier.Copy(&user, &userCreate)
	if userCreate.Id != nil {
		user.ID = *userCreate.Id
	} else {
		user.ID = "" // or generate UUID here
	}
	return &user
}
