package mysql

import (
	"database/sql"
	"time"

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
func (m *FlavorModel) Insert(flavor *models.Flavor) (*models.Flavor, error) {
	created := time.Now()
	tx, _ := m.DB.Begin()
	defer tx.Rollback()

	stmt := `INSERT INTO flavor (name, description, created) VALUES (?, ?, ?)`
	res, err := tx.Exec(stmt, flavor.Name, flavor.Description, created)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	flavorId, err := res.LastInsertId()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	flavor.ID = flavorId
	flavor.Created = created

	ingredientsModel := IngredientModel{DB: m.DB}

	for index, ingredient := range flavor.Ingredients {
		i, err := ingredientsModel.GetByName(ingredient.Name)
		if err == models.ErrNoRecord {
			i, err = ingredientsModel.Insert(&ingredient)
		}

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		flavor.Ingredients[index].ID = i.ID

		stmt = `INSERT INTO flavor_ingredient (flavor_id, ingredient_id) VALUES (?, ?)`
		_, err = tx.Exec(stmt, flavorId, i.ID)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return flavor, nil
}

// Update a Flavor identified by its ID.
func (m *FlavorModel) Update(id int, flavor *models.Flavor) (*models.Flavor, error) {
	return nil, nil
}

// Delete a Flavor identified by ID.
func (m *FlavorModel) Delete(id int) (bool, error) {
	return false, nil
}

// AddIngredient adds an Ingredient to a Flavor.
func (m *FlavorModel) AddIngredient(flavorId int, ingredient *models.Ingredient) (*models.Ingredient, error) {

	// Add Ingredient ID to flavor_ingredient table
	return ingredient, nil
}

// RemoveIngredient removes an Ingredient from a Flavor.
func (m *FlavorModel) RemoveIngredient(id int, ingredient *models.Ingredient) (*models.Ingredient, error) {
	return nil, nil
}
