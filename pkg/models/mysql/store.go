package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/jcorry/morellis/pkg/models"
)

type StoreModel struct {
	DB *sql.DB
}

// List stores. Length of list is defined by `limit`, beginning at `offset`. List is sorted by `order`.
func (s *StoreModel) List(limit int, offset int, order string) ([]*models.Store, error) {
	stmt := fmt.Sprintf(`SELECT s.id, s.name, s.phone, s.email, s.url, s.address, s.city, s.state, s.zip, s.lat, s.lng, s.created
								  FROM store AS s
							  ORDER BY %s
								 LIMIT ?, ?`, order)

	if limit < 1 {
		limit = DEFAULT_LIMIT
	}

	rows, err := s.DB.Query(stmt, offset, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stores []*models.Store

	for rows.Next() {
		store := &models.Store{}
		err = rows.Scan(&store.ID, &store.Name, &store.Phone, &store.Email, &store.URL, &store.Address, &store.City, &store.State, &store.Zip, &store.Lat, &store.Lng, &store.Created)
		if err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}

// Insert a new Store
func (s *StoreModel) Insert(name string, phone string, email string, url string, address string, city string, state string, zip string, lat float64, lng float64) (*models.Store, error) {
	created := time.Now()
	stmt := `INSERT INTO store (
		name,
		phone,
		email,
		url,
		address,
		city,
		state,
		zip,
		lat,
		lng,
		created
	) VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)`

	result, err := s.DB.Exec(stmt, name, phone, email, url, address, city, state, zip, lat, lng, created)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	store := &models.Store{
		ID:      id,
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
		Created: created,
	}

	return store, nil
}

// Get a single Store by ID
func (s *StoreModel) Get(id int) (*models.Store, error) {
	stmt := `SELECT id, name, phone, email, url, phone, address, city, state, zip, lat, lng, created
			   FROM store
		  	  WHERE id = ?`

	store := &models.Store{}
	err := s.DB.QueryRow(stmt, id).Scan(&store.ID, &store.Name, &store.Phone, &store.Email, &store.URL, &store.Phone, &store.Address, &store.City, &store.State, &store.Zip, &store.Lat, &store.Lng, &store.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return store, nil
}

// Update a Store identified by it's ID.
func (s *StoreModel) Update(ID int, name string, phone string, email string, url string, address string, city string, state string, zip string, lat float64, lng float64) (*models.Store, error) {
	updated := time.Now()
	stmt := `
	UPDATE store SET
		name = ?,
		phone = ?,
		email = ?,
		url = ?,
		address = ?,
		city = ?,
		state = ?,
		zip = ?,
		lat = ?,
		lng = ?,
		updated = ?
	WHERE id = ?`

	_, err := s.DB.Exec(stmt, name, phone, email, url, address, city, state, zip, lat, lng, updated, ID)
	if err != nil {
		return nil, err
	}

	store, err := s.Get(ID)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *StoreModel) Count() int {
	var count int
	stmt := `SELECT COUNT(id) FROM store WHERE 1`

	row := s.DB.QueryRow(stmt)
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}
	return count
}

// ActivateFlavor adds an active Flavor to the indicated Position at a Store, deactivating the Flavor
// currently occupying that Position.
func (s *StoreModel) ActivateFlavor(storeID int64, flavorID int64, position int) error {
	tx, _ := s.DB.Begin()
	defer tx.Rollback()

	stmt := `UPDATE flavor_store 
				SET is_active = NULL, deactivated = CURRENT_TIMESTAMP 
			  WHERE store_id = ?
				AND position = ?`

	_, err := s.DB.Exec(stmt, storeID, position)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt = `INSERT INTO flavor_store (store_id, flavor_id, position, is_active, activated)
			VALUES(?, ?, ?, 1, CURRENT_TIMESTAMP)`

	_, err = s.DB.Exec(stmt, storeID, flavorID, position)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		// If line 175 executed, this should never happen
		if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "uk_flavor_store_is_active_store_id_position_id") {
			tx.Rollback()
			return models.ErrDuplicateFlavor
		}
	}
	tx.Commit()

	return err
}

// DeactivateFlavor deactivates the Flavor identified by its ID at the indicated Store.
// note that if there are more than one instance of the flavor active at the store, all
// instances will be deactivated. To deactivate a single instance, use `DeactivateFlavorAtPosition`
func (s *StoreModel) DeactivateFlavor(storeID int64, flavorID int64) (bool, error) {
	stmt := `UPDATE flavor_store
				SET is_active = NULL, deactivated = CURRENT_TIMESTAMP
			  WHERE store_id = ?
				AND flavor_id = ?`
	res, err := s.DB.Exec(stmt, storeID, flavorID)
	if err != nil {
		return false, err
	}

	a, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	if a < 1 {
		return false, nil
	}

	return true, nil
}

// DeactivateFlavorAtPosition deactivates the Flavor in the indicated Position at the indicated Store.
// returns false if no rows were updated, true if rows were updated, error otherwise
func (s *StoreModel) DeactivateFlavorAtPosition(storeID int64, position int) (bool, error) {
	stmt := `UPDATE flavor_store
				SET is_active = NULL, deactivated = CURRENT_TIMESTAMP
			  WHERE store_id = ?
			    AND position = ?`

	res, err := s.DB.Exec(stmt, storeID, position)

	if err != nil {
		return false, err
	}

	a, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	if a < 1 {
		return false, nil
	}

	return true, nil
}

// GetActiveFlavors returns a collection of the currently active flavors at a store.
func (s *StoreModel) GetActiveFlavors(storeID int64) ([]*models.Flavor, error) {
	stmt := `SELECT f.id, f.name, f.description, f.created, i.id, i.name 
			   FROM flavor AS f
		  LEFT JOIN flavor_ingredient AS fi ON fi.flavor_id = f.id
		  LEFT JOIN ingredient AS i ON fi.ingredient_id = i.id
		  LEFT JOIN flavor_store AS fs ON fs.flavor_id = f.id
			  WHERE fs.store_id = ?
			    AND fs.is_active = 1
		   ORDER BY fs.position ASC`

	rows, err := s.DB.Query(stmt, storeID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flavors []*models.Flavor
	flavor := &models.Flavor{}

	var (
		flavorID    int64
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

		if flavor.ID != flavorID {
			flavors = append(flavors, flavor)
		}

		flavorID = flavor.ID
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return flavors, nil
}
