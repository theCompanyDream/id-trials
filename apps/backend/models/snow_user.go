package models

import (
	"strconv"

	"github.com/jinzhu/copier"
)

// UserDTO represents a user in the system.
type UserSnowflake struct {
	ID int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:false" json:"id"`
	UserBase
}

// TableName sets the table name for UserDTO to "users".
func (UserSnowflake) TableName() string {
	return "users_snowflake"
}

func InputToSnowFlake(userCreate UserInput) *UserSnowflake {
	var user UserSnowflake
	copier.Copy(&user, &userCreate)
	user.ID = 0

	if userCreate.Id != nil {
		return &user
	}

	if idInt, err := strconv.ParseInt(*userCreate.Id, 10, 64); err == nil {
		user.ID = idInt
	}

	return &user
}
