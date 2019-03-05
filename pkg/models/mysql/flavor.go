package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

// FlavorModel is a wrapper for a DB struct and the methods.
type FlavorModel struct {
	DB *sql.DB
}

// Get a single Flavor by it's ID.
func (m *FlavorModel) Get(id int) (*models.Flavor, error) {
	stmt := `SELECT f.id, f.name, f.description, f.created, i.id, i.name
			   FROM flavor AS f
		  LEFT JOIN flavor_ingredient AS fi ON f.id = fi.flavor_id
		  LEFT JOIN ingredient AS i ON i.id = fi.ingredient_id
			  WHERE f.id = ?`

	flavor := &models.Flavor{}

	rows, err := m.DB.Query(stmt, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ingredient := &models.Ingredient{}
		err = rows.Scan(&flavor.ID, &flavor.Name, &flavor.Description, &flavor.Created, &ingredient.ID, &ingredient.Name)

		if err != nil {
			return nil, err
		}

		flavor.Ingredients = append(flavor.Ingredients, *ingredient)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return flavor, nil
}

// List {limit} number of Flavors starting at {offset}. If {order} matches a field name,
// results will be ordered by {order}.
func (m *FlavorModel) List(limit int, offset int, order string) ([]*models.Flavor, error) {

	stmt := fmt.Sprintf(`SELECT f.id, f.name, f.description, f.created, i.id, i.name
			   FROM flavor AS f
		       JOIN flavor_ingredient AS fi ON f.id = fi.flavor_id
		  LEFT JOIN ingredient AS i ON i.id = fi.ingredient_id
		   ORDER BY %s
			  LIMIT ?, ?`, "f.name")

	if limit < 1 {
		limit = DEFAULT_LIMIT
	}

	rows, err := m.DB.Query(stmt, offset, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flavors := []*models.Flavor{}
	flavor := &models.Flavor{}

	var (
		flavorId    int64
		id          int64
		name        string
		description string
		created     time.Time
		ingredient  models.Ingredient
	)

	for rows.Next() {
		err = rows.Scan(&id, &name, &description, &created, &ingredient.ID, &ingredient.Name)
		if err != nil {
			return nil, err
		}

		if flavor != nil && flavor.ID == id {
			flavor.Ingredients = append(flavor.Ingredients, ingredient)
		} else {
			flavor = &models.Flavor{
				ID:          id,
				Name:        name,
				Description: description,
				Created:     created,
				Ingredients: []models.Ingredient{ingredient},
			}
		}

		if flavor.ID != flavorId {
			flavors = append(flavors, flavor)
		}

		flavorId = flavor.ID
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return flavors, nil
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

	// now handle each of the ingredients
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
