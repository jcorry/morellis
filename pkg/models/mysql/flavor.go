package mysql

import (
	"database/sql"

	"github.com/jcorry/morellis/pkg/models"
)

// FlavorModel is a wrapper for a DB struct and the methods.
type FlavorModel struct {
	DB *sql.DB
}

// Get a single Flavor by it's ID.
func (m *FlavorModel) Get(id int) (*models.Flavor, error) {
	return nil, nil
}

// List {limit} number of Flavors starting at {offset}. If {order} matches a field name,
// results will be ordered by {order}.
func (m *FlavorModel) List(limit int, offset int, order string) ([]*models.Flavor, error) {
	return nil, nil
}

// Insert a new Flavor with it's Ingredients.
func (m *FlavorModel) Insert(name string, ingredients []*models.Ingredient) (*models.Flavor, error) {
	return nil, nil
}

// Update a Flavor identified by its ID.
func (m *FlavorModel) Update(id int, name string, ingredients []*models.Ingredient) (*models.Flavor, error) {
	return nil, nil
}

// Delete a Flavor identified by ID.
func (m *FlavorModel) Delete(id int) (bool, error) {
	return false, nil
}

// AddIngredient adds an Ingredient to a Flavor.
func (m *FlavorModel) AddIngredient(ingredient *models.Ingredient) (*models.Flavor, error) {
	return nil, nil
}

// RemoveIngredient removes an Ingredient from a Flavor.
func (m *FlavorModel) RemoveIngredient(ingredient *models.Ingredient) (*models.Flavor, error) {
	return nil, nil
}
