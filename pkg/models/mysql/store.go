package mysql

import (
	"database/sql"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

type StoreModel struct {
	DB *sql.DB
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
	stmt := `SELECT id, name, phone, email, url, phone, address, city, state, zip, lat, lng, created, updated
			   FROM store
		  	  WHERE id = ?`

	store := &models.Store{}
	err := s.DB.QueryRow(stmt, id).Scan(&store.ID, &store.Name, &store.Phone, &store.Email, &store.URL, &store.Phone, &store.Address, &store.City, &store.State, &store.Zip, &store.Lat, &store.Lng, &store.Created, &store.Updated)

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
