package models

type UserBase struct {
	UserName   string    `gorm:"column:user_name;type:varchar(20);not null" json:"user_name"`
	FirstName  string    `gorm:"column:first_name;type:varchar(40);not null" json:"first_name"`
	LastName   string    `gorm:"column:last_name;type:varchar(40);not null" json:"last_name"`
	Email      string    `gorm:"column:email;type:varchar(40);not null;unique" json:"email"`
	Department *string   `gorm:"column:department;type:varchar(25)" json:"department"`
}

type UserInput struct {
	// HashId is the public identifier for the user (UUID).
	// For create operations, this might be generated internally.
	Id *string `json:"id" validate:"omitempty,uuid4" form:"id"`

	// UserName is required when creating a new user.
	UserName *string `json:"user_name" validate:"omitempty,min=5,max=50" form:"user_name"`

	// FirstName is required when creating a new user.
	FirstName *string `json:"first_name" validate:"omitempty,min=3,max=255" form:"first_name"`

	// LastName is required when creating a new user.
	LastName *string `json:"last_name" validate:"omitempty,min=3,max=255" form:"last_name"`

	// Email is required when creating a new user.
	Email *string `json:"email" validate:"omitempty,email,max=255" form:"email"`

	// Department is optional.
	Department *string `json:"department" form:"department"`
}