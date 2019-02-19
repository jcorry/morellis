package models

type Ingredient struct {
	Base
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Flavors     []Flavor `json:"-";gorm:"many2many:flavors_ingredients;"`
}

type IngredientInList struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Flavors     []Flavor `json:"-";gorm:"many2many:flavors_ingredients;"`
}
