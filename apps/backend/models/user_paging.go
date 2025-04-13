package models

// Paging defines the structure for pagination information
// @Description Paging
type Paging struct {
	// The current page number
	Page *int `json:"page"`
	// The total number of items available
	PageCount *int `json:"page_count"`
	// The number of items per page
	PageSize *int `json:"page_size"`
}

// UserDTOPaging defines the structure for user pagination which includes the users and pagination details
// @Description UserDTOPaging
type UserPaging struct {
	Paging // Embedded Paging struct
	// A list of users
	Users []UserDTO `json:"users"`
}
