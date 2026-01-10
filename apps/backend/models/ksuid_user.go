package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserKSUID struct {
	ID string `gorm:"column:id;type:varchar(27);primaryKey" json:"id"`
	*UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserKSUID) TableName() string {
	return "users_ksuid"
}

func InputToKSUID(userCreate UserInput) *UserKSUID {
	var user UserKSUID
	copier.Copy(&user, &userCreate)
	if userCreate.Id != nil {
		user.ID = *userCreate.Id
	} else {
		user.ID = ""
	}
	return &user
}

func (ksuid *UserKSUID) KsuidToDTO() *UserDTO {
	var userDTO UserDTO
	copier.Copy(&userDTO, &ksuid)
	return &userDTO
}
