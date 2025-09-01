package cmd

import "github.com/brianvoe/gofakeit/v6"

var departments = []string{
	"Engineering",
	"Sales",
	"Marketing",
	"HR",
	"Operations",
	"Finance",
	"IT",
	"Legal",
}

func randomDepartment() *string {
	dept := gofakeit.RandomString(departments)
	return &dept
}
