package mock

import "github.com/jcorry/morellis/pkg/models"

type IngredientModel struct{}

func (m *IngredientModel) Get(ID int64) (*models.Ingredient, error) {
	if ID < 100 {
		return &models.Ingredient{
			ID:   ID,
			Name: "Vanilla",
		}, nil
	} else {
		return nil, models.ErrNoRecord
	}
}

func (m *IngredientModel) GetByName(name string) (*models.Ingredient, error) {
	return nil, nil
}

func (m *IngredientModel) Search(limit int, offset int, order string, search []string) ([]*models.Ingredient, error) {
	return nil, nil
}

func (m *IngredientModel) Insert(ingredient *models.Ingredient) (*models.Ingredient, error) {
	return nil, nil
}
