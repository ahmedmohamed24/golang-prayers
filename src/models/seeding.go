package models

import "gorm.io/gorm"

// each seeding has its name
// this approach help us to prevent running the same seeder twice
type Seeding struct {
	*gorm.Model
	Name string `json:"seeder_name"`
}
