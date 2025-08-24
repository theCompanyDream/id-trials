package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserNanoID struct {
	ID string `gorm:"column:id;type:varchar(27);primaryKey" json:"id"`
	*UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserNanoID) TableName() string {
	return "users_nanoid"
}

func InputToNanoId(userCreate UserInput) *UserNanoID {
	var user UserNanoID
	copier.Copy(&user, &userCreate)
	if userCreate.Id != nil {
		user.ID = *userCreate.Id
	} else {
		user.ID = "" // or generate UUID here
	}
	return &user
}
