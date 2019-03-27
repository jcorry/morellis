package mysql

import (
	"testing"
	"time"
)

func TestStoreModel_List(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}
	db, teardown := newTestDB(t)
	defer teardown()

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
				if s.Name != tt.wantNames[i] {
					t.Errorf("want %s; got %s", tt.wantNames[i], s.Name)
				}
			}
		})
	}
}

func TestStoreModel_ActivateFlavor(t *testing.T) {

}

func TestStoreModel_GetActiveFlavors(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	m := StoreModel{db}

	// Insert into flavor_store a few rows for testing
	flavorStoreEntries := []struct {
		flavorID int64
		storeID  int64
		position int64
		active   bool
	}{
		{1, 1, 1, true},
		{2, 1, 2, true},
	}

	for _, entry := range flavorStoreEntries {
		stmt := `INSERT INTO flavor_store 
			(flavor_id, store_id, position, is_active, activated, deactivated) VALUES (?, ?, ?, ?, ?, ?)`

		active := 0
		now := time.Now()
		var activated time.Time
		deactivated := now.AddDate(0, 0, -1)

		if entry.active {
			active = 1
			activated = now.AddDate(0, -1, 0)
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
}
