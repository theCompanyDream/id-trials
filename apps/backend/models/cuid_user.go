package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserCUID struct {
	ID string `gorm:"column:id;type:varchar(25);primaryKey" json:"id"`
	*UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserCUID) TableName() string {
	return "users_cuid"
}

func InputToCuid(userCreate UserInput) *UserCUID {
	var user UserCUID
	copier.Copy(&user, &userCreate)
	if userCreate.Id != nil {
		user.ID = *userCreate.Id
	} else {
		user.ID = "" // or generate UUID here
	}
	return &user
}

func (cuid *UserCUID) CuidToDTO() *UserDTO {
	var userDTO UserDTO
	copier.Copy(&userDTO, &cuid)
	return &userDTO
}
