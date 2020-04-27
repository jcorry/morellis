package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

func TestStoreModel_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}
	db, teardown, err := newTestDB(t)
	if err != nil {
		t.Errorf("unexpected err getting test DB:\n%v", err)
	}
	t.Cleanup(func() {
		teardown()
	})

	name := "New Store"
	phone := "867-5309"
	email := "example@example.com"
	url := "http://example.com"
	address := "123 Any Ln"
	city := "Anywhere"
	state := "CA"
	zip := "10111"
	lat := 32.476
	lng := -89.234

	m := StoreModel{db}

	store, err := m.Insert(name, phone, email, url, address, city, state, zip, lat, lng)
	if err != nil {
		t.Errorf("unexpected err inserting store: %v", err)
	}
	if store.Name != name {
		t.Errorf("Want %s; Got %s", name, store.Name)
	}

	if store.ID <= 0 {
		t.Error("Store ID wasn't set")
	}
}

func TestStoreModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}
	db, teardown, err := newTestDB(t)
	if err != nil {
		t.Errorf("unexpected err getting test DB:\n%v", err)
	}
	t.Cleanup(func() {
		teardown()
	})

	m := StoreModel{db}

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{"valid record", 1, nil},
		{"no record found", 1000, sql.ErrNoRows},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.Get(tt.id)
			if err != nil {
				if fmt.Sprintf("%T", err) != fmt.Sprintf("%T", tt.wantErr) {
					t.Errorf("Got %s; want %s", err, tt.wantErr)
				}
			}
		})
	}
}

func TestStoreModel_List(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}
	db, teardown, err := newTestDB(t)
	if err != nil {
		t.Errorf("unexpected err getting test DB:\n%v", err)
	}
	t.Cleanup(func() {
		teardown()
	})

	m := StoreModel{db}

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
				s := s
				if !strings.Contains(strings.Join(tt.wantNames, " "), s.Name) {
					t.Errorf("want %s; got %s", tt.wantNames[i], s.Name)
				}
			}
		})
	}
}

func TestStoreModel_ActivateFlavor(t *testing.T) {
	db, teardown, err := newTestDB(t)
	if err != nil {
		t.Errorf("unexpected err getting test DB:\n%v", err)
	}
	t.Cleanup(func() {
		teardown()
	})

	s := StoreModel{db}
	f := FlavorModel{db}

	store, err := s.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	flavor, err := f.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	err = s.ActivateFlavor(store.ID, flavor.ID, 1)
	if err != nil {
		t.Fatal(err)
	}

	activeFlavors, err := s.GetActiveFlavors(store.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(activeFlavors) != 1 {
		t.Errorf("Want 1; Got %d", len(activeFlavors))
	}

	err = s.ActivateFlavor(store.ID, flavor.ID, 1)
	if err != nil && err != models.ErrDuplicateFlavor {
		t.Errorf("Error: %v", err)
		t.Errorf("Want models.ErrDuplicateFlavor; Got %T", err)
	}
}

func TestStoreModel_GetActiveFlavors(t *testing.T) {
	db, teardown, err := newTestDB(t)
	if err != nil {
		t.Errorf("unexpected err getting test DB:\n%v", err)
	}
	t.Cleanup(func() {
		teardown()
	})

	m := StoreModel{db}

	// Insert into flavor_store a few rows for testing
	flavorStoreEntries := []struct {
		storeID  int64
		flavorID int64
		position int64
		active   bool
	}{
		{1, 1, 1, false},
		{1, 1, 1, true},
		{1, 2, 2, true},
	}

	// Insert some relationship data to test
	for _, entry := range flavorStoreEntries {
		stmt := `INSERT INTO flavor_store 
			(flavor_id, store_id, position, is_active, activated, deactivated) VALUES (?, ?, ?, ?, ?, ?)`

		active := 0

		activated := time.Now()
		deactivated := activated.AddDate(0, 0, -1)

		if entry.active {
			active = 1
			activated = activated.AddDate(0, -1, 0)
		}

		if entry.active {
			_, err := m.DB.Exec(stmt, entry.flavorID, entry.storeID, entry.position, active, activated, nil)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			_, err := m.DB.Exec(stmt, entry.flavorID, entry.storeID, entry.position, active, activated, deactivated)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	tests := []struct {
		name         string
		storeID      int64
		wantRowCount int
		wantError    error
	}{
		{"Get active rows", 1, 2, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flavors, err := m.GetActiveFlavors(tt.storeID)
			if err != nil {
				t.Errorf("Got MySQL error: %#v", err)
			}

			if len(flavors) != tt.wantRowCount {
				t.Errorf("Want row count %d; got %d", tt.wantRowCount, len(flavors))
			}
		})
	}
}

func TestStoreModel_DeactivateFlavor(t *testing.T) {
	db, teardown, err := newTestDB(t)
	if err != nil {
		t.Errorf("unexpected err getting test DB:\n%v", err)
	}
	t.Cleanup(func() {
		teardown()
	})

	s := StoreModel{db}

	tests := []struct {
		name     string
		storeID  int64
		flavorID int64
		wantErr  error
		wantRes  bool
	}{
		{"Deactivate a non existent row", 1, 1, nil, false},
		{"Deactivate an existing row", 1, 1, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantRes {
				// If we want a true result, we have to seed a row...
				err := s.ActivateFlavor(tt.storeID, tt.flavorID, 1)
				if err != nil {
					t.Fatal(err)
				}
			}
			res, err := s.DeactivateFlavor(tt.storeID, tt.flavorID)
			if err != tt.wantErr {
				t.Errorf("Want err %s, Got err %s", tt.wantErr, err)
			}
			if res != tt.wantRes {
				t.Errorf("Want result %v, Got result %v", tt.wantRes, res)
			}
		})
	}
}

func TestStoreModel_DeactivateFlavorAtPosition(t *testing.T) {
	db, teardown, err := newTestDB(t)
	if err != nil {
		t.Errorf("unexpected err getting test DB:\n%v", err)
	}
	t.Cleanup(func() {
		teardown()
	})

	s := StoreModel{db}

	tests := []struct {
		name     string
		storeID  int64
		position int
		wantErr  error
		wantRes  bool
	}{
		{"Deactivate a non existent row", 1, 1, nil, false},
		{"Deactivate an existing row", 1, 1, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantRes {
				// If we want a true result, we have to seed a row...
				err := s.ActivateFlavor(tt.storeID, 1, tt.position)
				if err != nil {
					t.Fatal(err)
				}
			}
			res, err := s.DeactivateFlavorAtPosition(tt.storeID, tt.position)
			if err != tt.wantErr {
				t.Errorf("Want err %s, Got err %s", tt.wantErr, err)
			}
			if res != tt.wantRes {
				t.Errorf("Want result %v, Got result %v", tt.wantRes, res)
			}
		})
	}
}
