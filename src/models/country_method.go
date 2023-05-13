package models

import (
	"gorm.io/gorm"
)

type CountryMethod struct {
	gorm.Model
	CountryName string `json:"country_name"`
	MethodName  string `json:"method_name"`
}
