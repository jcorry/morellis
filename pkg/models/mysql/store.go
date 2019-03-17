package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

type StoreModel struct {
	DB *sql.DB
}

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

// Get a single User by ID
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

func (s *StoreModel) Update(id int, name string, phone string, email string, url string, address string, city string, state string, zip string, lat float64, lng float64) (*models.Store, error) {
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

	_, err := s.DB.Exec(stmt, name, phone, email, url, address, city, state, zip, lat, lng, updated, id)
	if err != nil {
		return nil, err
	}

	store, err := s.Get(id)
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
