package models

import (
	"fmt"

	u "github.com/jcorry/morellis/utils"

	"github.com/jinzhu/gorm"
)

// A struct to represent a flavor
type Flavor struct {
	gorm.Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Ingredients []Ingredient `gorm:"many2many:flavors_ingredients;"json:"ingredients"`
}

func (flavor *Flavor) Validate() (map[string]interface{}, bool) {
	if len(flavor.Name) == 0 {
		return u.Message(false, "Name is required."), false
	}

	if len(flavor.Description) == 0 {
		return u.Message(false, "Description is required."), false
	}

	if len(flavor.Ingredients) > 0 {
		for index, ingredient := range flavor.Ingredients {
			if len(ingredient.Name) == 0 {
				return u.Message(false, fmt.Sprintf("Ingredient name is required (ingredient[%d].", index)), false
			}
		}
	}

	return u.Message(false, "Validation passed"), true
}

// Create a new flavor
func (flavor *Flavor) Create() map[string]interface{} {
	if resp, ok := flavor.Validate(); !ok {
		return resp
	}

	GetDB().Model(&flavor).Create(&flavor)
	GetDB().Model(&flavor).Save(&flavor)

	if flavor.ID <= 0 {
		return u.Message(false, "Failed to create flavor, DB connection error.")
	}

	response := u.Message(true, "Flavor has been created")
	response["flavor"] = flavor
	return response
}

// Get a single flavor by ID
func GetFlavor(id uint) *Flavor {
	flavor := &Flavor{}
	err := GetDB().Table("flavors").Where("id = ?", id).First(flavor).Error
	if err != nil {
		return nil
	}

	return flavor
}

// Get all of the flavors in the flavors table
func GetFlavors() []*Flavor {
	flavors := make([]*Flavor, 0)
	err := GetDB().Table("flavors").Find(&flavors).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return flavors
}
