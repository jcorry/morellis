package mock

import (
	"sort"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

type FlavorModel struct{}

var MockFlavors = []*models.Flavor{
	{
		ID:          1,
		Name:        "Ginger Lavender",
		Description: "Fragrant and unique, Ginger Lavender is a fan favorite and popular seasonal flavor! This delicate and delicious flavor is made with ginger ice cream and fresh, fragrant lavender, steeped in heavy cream.",
		Ingredients: []models.Ingredient{
			{
				ID:   5,
				Name: "Ginger",
			},
			{
				ID:   6,
				Name: "Lavender",
			},
		},
		Created: time.Now().Add(time.Second * 1),
	},
	{
		ID:          2,
		Name:        "Dark Chocolate Chili",
		Description: "Deep, dark and spicy! Our Dark Chocolate Chilli ice cream infuses chilli, cayenne pepper and a dark chocolate fudge swirl into our dark chocolate ice cream made with Guittard chocolate. Decadent and DELICIOUS!",
		Ingredients: []models.Ingredient{
			{
				ID:   3,
				Name: "Dark Chocolate",
			},
			{
				ID:   4,
				Name: "Chili",
			},
		},
		Created: time.Now().Add(time.Second * 2),
	},
	{
		ID:          3,
		Name:        "Coconut Jalapeño",
		Description: "One of our most unique flavors, it must be tasted to be believed! \n\nOur fresh made coconut ice cream is infused with just the right amount of fresh jalapenos. The experience of hot, sweet and cold hits your palate in pretty amazing ways; come try for yourself!",
		Ingredients: []models.Ingredient{
			{
				ID:   1,
				Name: "Coconut",
			},
			{
				ID:   2,
				Name: "Jalapeño",
			},
		},
		Created: time.Now().Add(time.Second * 3),
	},
	{
		ID:          4,
		Name:        "Smooth Monkey",
		Description: "Were not monkeying around with this flavor! Swing through the shop for a scoop of Morellis intense banana ice cream, filled with a thick, dark chocolate fudge swirl and drizzled with walnuts!",
		Ingredients: []models.Ingredient{
			{
				ID:   3,
				Name: "Dark Chocolate",
			},
			{
				ID:   7,
				Name: "Banana",
			},
			{
				ID:   8,
				Name: "Walnuts",
			},
		},
		Created: time.Now().Add(time.Second * 4),
	},
}

func (m *FlavorModel) Get(id int) (*models.Flavor, error) {
	for _, f := range MockFlavors {
		if int64(id) == f.ID {
			return f, nil
		}
	}
	return nil, models.ErrNoRecord
}

func (m *FlavorModel) List(limit int, offset int, order string) ([]*models.Flavor, error) {
	// sort them by name
	sort.Slice(MockFlavors, func(i, j int) bool {
		return MockFlavors[i].Name < MockFlavors[i].Name
	})
	return MockFlavors, nil
}

func (m *FlavorModel) Insert(flavor *models.Flavor) (*models.Flavor, error) {
	return flavor, nil
}

func (m *FlavorModel) Update(id int, flavor *models.Flavor) (*models.Flavor, error) {
	return nil, nil
}
func (m *FlavorModel) Delete(id int) (bool, error) {
	return true, nil
}

func (m *FlavorModel) AddIngredient(id int, ingredient *models.Ingredient) (*models.Ingredient, error) {
	return nil, nil
}

func (m *FlavorModel) RemoveIngredient(id int, ingredient *models.Ingredient) (*models.Ingredient, error) {
	return nil, nil
}

func (m *FlavorModel) Count() int {
	return len(MockFlavors)
}
