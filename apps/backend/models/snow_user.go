package models

import (
	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserSnowflake struct {
	ID         string `gorm:"column:id;type:varchar(27);primaryKey" json:"id"`
	UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserSnowflake) TableName() string {
	return "users_snowflake"
}

func InputToSnowFlake(userCreate UserInput) *UserSnowflake {
	var user UserSnowflake
	copier.Copy(&user, &userCreate)
	return &user
}
