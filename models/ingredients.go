package models

import "github.com/jinzhu/gorm"

type Ingredient struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Flavors     []Flavor `json:"flavors";gorm:"many2many:flavor_ingredients;save_associations:true"`
}
