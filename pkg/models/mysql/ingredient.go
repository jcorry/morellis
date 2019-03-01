package mysql

import (
	"database/sql"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

// IngredientModel is a wrapper for a DB struct and the methods.
type IngredientModel struct {
	DB *sql.DB
}

// GetByName retrieves an Ingredient by its Name.
func (m *IngredientModel) GetByName(name string) (*models.Ingredient, error) {
	var ingredient = &models.Ingredient{}
	stmt := `SELECT id, name FROM ingredient WHERE LOWER(name) = ?`
	err := m.DB.QueryRow(stmt, name).Scan(&ingredient.ID, &ingredient.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return ingredient, nil
}

// Insert inserts a new Ingredient into the DB
func (m *IngredientModel) Insert(ingredient *models.Ingredient) (*models.Ingredient, error) {
	created := time.Now()
	stmt := `INSERT INTO ingredient (name, created) VALUES (?, ?)`
	res, err := m.DB.Exec(stmt, ingredient.Name, created)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	ingredient.ID = id

	return ingredient, err
}
