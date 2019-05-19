package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

// IngredientModel is a wrapper for a DB struct and the methods.
type IngredientModel struct {
	DB *sql.DB
}

// Get retrieves a single Ingredient by its ID
func (m *IngredientModel) Get(ID int64) (*models.Ingredient, error) {
	var i = &models.Ingredient{}
	stmt := `SELECT id, name FROM ingredient WHERE id = ?`

	err := m.DB.QueryRow(stmt, ID).Scan(&i.ID, &i.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return i, nil
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

func (m *IngredientModel) Search(limit int, offset int, order string, search []string) ([]*models.Ingredient, error) {
	args := make([]interface{}, len(search))
	for i, term := range search {
		term = strings.ToLower(strings.TrimSpace(term))
		args[i] = fmt.Sprintf("%%%s%%", term)
	}

	stmt := `SELECT id, name FROM ingredient WHERE LOWER(name) LIKE ?`

	for i := 1; i <= len(args)-1; i++ {
		stmt += ` OR LOWER(name) LIKE ?`
	}

	// Only order by one of the available field names
	fields := []string{"id", "name", "created"}
	for _, field := range fields {
		if order == field {
			stmt += fmt.Sprintf(" ORDER BY %s", field)
		}
	}

	stmt += ` LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := []*models.Ingredient{}
	for rows.Next() {
		ingredient := &models.Ingredient{}
		err = rows.Scan(&ingredient.ID, &ingredient.Name)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
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
