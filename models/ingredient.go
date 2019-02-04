package models

import "github.com/jinzhu/gorm"

type Ingredient struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Flavors     []Flavor `json:"-";gorm:"many2many:flavors_ingredients;"`
}
