package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record(s) found")

// Flavor is an ice cream flavor served by Morellis at any of it's Stores.
type Flavor struct {
	ID          int32
	Name        string
	Ingredients []string
	Created     time.Time
}

// User is a user of the system
type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
}

type UserStatus int

const (
	USER_STATUS_UNVERIFIED UserStatus = 1
	USER_STATUS_VERIFIED   UserStatus = 2
	USER_STATUS_DELETED    UserStatus = 3
)

// Slug returns a textual slug for the given UserStatus
func (status UserStatus) Slug() string {
	names := make(map[UserStatus]string)
	names[USER_STATUS_UNVERIFIED] = "unverified"
	names[USER_STATUS_VERIFIED] = "verified"
	names[USER_STATUS_DELETED] = "deleted"

	return names[status]
}

// GetID returns the UserStatus for a given textual slug
func (status UserStatus) GetID(slug string) UserStatus {
	names := make(map[string]UserStatus)
	names["unverified"] = USER_STATUS_UNVERIFIED
	names["verified"] = USER_STATUS_VERIFIED
	names["deleted"] = USER_STATUS_DELETED

	return names[slug]
}

// Store is an instance of a Morelli's store
type Store struct {
	ID    int
	Name  string
	Phone string
	Location
}

// Location is a geographic location
type Location struct {
	Address string
	City    string
	State   string
	Zip     string
	Lat     float64
	Lng     float64
}
