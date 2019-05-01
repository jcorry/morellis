package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoRecord                = errors.New("models: No matching record(s) found")
	ErrInvalidCredentials      = errors.New("models: Invalid credentials")
	ErrDuplicateEmail          = errors.New("models: Duplicate email")
	ErrDuplicateFlavor         = errors.New("models: Only one flavor may be active at a position at a time.")
	ErrInvalidPermission       = errors.New("models: Not a valid Permission")
	ErrDuplicateUserPermission = errors.New("models: User already has that Permission")
	ErrInvalidUser             = errors.New("models: Not a valid User")
)

// Credentials are used to authenticate with the API
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Flavor is an ice cream flavor served by Morellis at any of it's Stores.
type Flavor struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	Created     time.Time    `json:"created"`
}

type Ingredient struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// User is a user of the system
type User struct {
	ID          int64            `json:"-"`
	UUID        uuid.UUID        `json:"uuid"`
	FirstName   string           `json:"firstName"`
	LastName    string           `json:"lastName"`
	Email       string           `json:"email"`
	Phone       string           `json:"phone"`
	Status      string           `json:"status"`
	Permissions []UserPermission `json:"permissions"`
	Password    string           `json:"password,omitempty"`
	Created     time.Time        `json:"created"`
}

type UserPermission struct {
	UserPermissionID int `json:"userPermissionId,omitempty"`
	Permission       `json:"permission"`
}

type Permission struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
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
	switch slug {
	case "unverified":
		return USER_STATUS_UNVERIFIED
	case "Unverified":
		return USER_STATUS_UNVERIFIED
	case "verified":
		return USER_STATUS_VERIFIED
	case "Verified":
		return USER_STATUS_VERIFIED
	case "deleted":
		return USER_STATUS_DELETED
	case "Deleted":
		return USER_STATUS_DELETED
	}
	return USER_STATUS_UNVERIFIED
}

// Store is an instance of a Morelli's store
type Store struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
	URL     string    `json:"url"`
	Address string    `json:"address"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	Zip     string    `json:"zip"`
	Lat     float64   `json:"lat"`
	Lng     float64   `json:"lng"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"-"`
}

func (s *Store) AddressString() string {
	return fmt.Sprintf("%s %s, %s %s", s.Address, s.City, s.State, s.Zip)
}
