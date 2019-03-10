package mock

import (
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

var mockStore = &models.Store{
	ID:      1,
	Name:    "Test Store",
	Phone:   "867-5309",
	Email:   "test@store.com",
	URL:     "http://www.testystore.com",
	Address: "123 Testy Ave",
	City:    "Test",
	State:   "TS",
	Zip:     "01010",
	Lat:     0.0,
	Lng:     0.0,
}

var MockStores = []*models.Store{
	mockStore,
	{
		ID:      2,
		Name:    "Another Store",
		Phone:   "867-5309",
		Email:   "test@store-two.com",
		URL:     "http://www.testystore.com",
		Address: "1427 Testy St",
		City:    "Anothertest",
		State:   "TS",
		Zip:     "02022",
		Lat:     22.22222,
		Lng:     122.22222,
	},
}

type StoreModel struct{}

func (m *StoreModel) Insert(name string, phone string, email string, url string, address string, city string, state string, zip string, lat float64, lng float64) (*models.Store, error) {
	store := &models.Store{
		ID:      101,
		Name:    name,
		Phone:   phone,
		Email:   email,
		URL:     url,
		Address: address,
		City:    city,
		State:   state,
		Zip:     zip,
		Lat:     lat,
		Lng:     lng,
		Created: time.Now(),
	}

	return store, nil
}

func (m *StoreModel) Update(id int, name string, phone string, email string, url string, address string, city string, state string, zip string, lat float64, lng float64) (*models.Store, error) {
	store := &models.Store{
		ID:      int64(id),
		Name:    name,
		Phone:   phone,
		Email:   email,
		URL:     url,
		Address: address,
		City:    city,
		State:   state,
		Zip:     zip,
		Lat:     lat,
		Lng:     lng,
		Created: time.Now(),
	}

	return store, nil
}

func (m *StoreModel) Get(id int) (*models.Store, error) {
	mockStore.ID = int64(id)
	return mockStore, nil
}

func (m *StoreModel) List(limit int, offset int, order string) ([]*models.Store, error) {
	return MockStores, nil
}

func (m *StoreModel) Count() int {
	return len(MockStores)
}
