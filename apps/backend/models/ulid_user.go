package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserUlid struct {
	ID string `gorm:"column:id;type:varchar(26);primaryKey" json:"id"`
	*UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserUlid) TableName() string {
	return "users_ulid"
}

func InputToUlid(userCreate UserInput) *UserUlid {
	var user UserUlid
	copier.Copy(&user, &userCreate)
	if userCreate.Id != nil {
		user.ID = *userCreate.Id
	} else {
		user.ID = "" // or generate UUID here
	}
	return &user
}

func (ulid *UserUlid) UlidToDTO() *UserDTO {
	var userDTO UserDTO
	copier.Copy(&userDTO, &ulid)
	return &userDTO
}
