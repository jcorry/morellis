package mysql

import (
	"testing"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

func TestStoreModel_List(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}
	db, teardown := newTestDB(t)
	defer teardown()

	m := StoreModel{db}

	stores := []*models.Store{
		{
			Name:    "Dunwoody Farmburger",
			Email:   "info@morellisicecream.com",
			Phone:   "404-622-0210",
			URL:     "http://morellisicecream.com",
			Address: "4514 Chamblee Dunwoody Rd",
			City:    "Dunwoody",
			State:   "GA",
			Zip:     "30338",
			Lat:     33.922714,
			Lng:     -84.315169,
		},
		{
			Name:    "Morellis On Moreland",
			Phone:   "404-622-0210",
			Email:   "info@morellisicecream.com",
			URL:     "http://www.morellisicecream.com/",
			Address: "749 Moreland Ave SE",
			City:    "Atlanta",
			State:   "GA",
			Zip:     "30316",
			Lat:     33.7339513,
			Lng:     -84.3496246,
		},
	}

	// setup the DB rows
	for _, s := range stores {
		_, err := m.Insert(s.Name, s.Phone, s.Email, s.URL, s.Address, s.City, s.State, s.Zip, s.Lat, s.Lng)
		if err != nil {
			t.Fatal(err.Error())
		}
		time.Sleep(time.Second)
	}

	tests := []struct {
		name      string
		limit     int
		wantNames []string
		wantError error
	}{
		{
			"Get all stores",
			0,
			[]string{"Dunwoody Farmburger", "Morellis On Moreland"},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := m.List(tt.limit, 0, "name")
			if err != tt.wantError {
				t.Errorf("want %q, got %s", tt.wantError, err)
			}

			if tt.limit > 0 && len(list) != tt.limit {
				t.Errorf("want list length %d; got %d", tt.limit, len(list))
			}

			if tt.limit == 0 && len(list) > DEFAULT_LIMIT {
				t.Errorf("want list length < %d, got list length %d", DEFAULT_LIMIT, len(list))
			}

			for i, s := range list {
				if s.Name != tt.wantNames[i] {
					t.Errorf("want %s; got %s", tt.wantNames[i], s.Name)
				}
			}
		})
	}
}
