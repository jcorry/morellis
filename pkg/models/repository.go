package models

import "github.com/google/uuid"

//go:generate counterfeiter . UserRepository
type UserRepository interface {
	Insert(uid uuid.UUID, firstName NullString, lastName NullString, email NullString, phone string, statusID int, password string) (*User, error)
	Update(*User) (*User, error)
	Get(int) (*User, error)
	GetByUUID(uuid.UUID) (*User, error)
	GetByCredentials(Credentials) (*User, error)
	GetByPhone(string) (*User, error)
	GetByAuthToken(string) (*User, error)
	SaveAuthToken(string, int) error
	List(int, int, string) ([]*User, error)
	Delete(int) (bool, error)
	Count() int
	GetPermissions(userID int) ([]UserPermission, error)
	AddPermission(userID int, p Permission) (int, error)
	RemovePermission(userPermissionID int) (bool, error)
	RemoveAllPermissions(userID int) error
	AddIngredient(userID int64, ingredient *Ingredient, keyword string) (*UserIngredient, error)
	GetIngredients(userID int64) ([]*UserIngredient, error)
	RemoveUserIngredient(userIngredientID int64) error
}

//go:generate counterfeiter . StoreRepository
type StoreRepository interface {
	Insert(string, string, string, string, string, string, string, string, float64, float64) (*Store, error)
	Update(int, string, string, string, string, string, string, string, string, float64, float64) (*Store, error)
	Get(storeID int) (*Store, error)
	List(int, int, string) ([]*Store, error)
	Count() int
	ActivateFlavor(storeID int64, flavorID int64, position int) error
	DeactivateFlavor(storeID int64, flavorID int64) (bool, error)
	DeactivateFlavorAtPosition(storeID int64, position int) (bool, error)
}

//go:generate counterfeiter . FlavorRepository
type FlavorRepository interface {
	Count() int
	Get(int) (*Flavor, error)
	List(limit int, offset int, sortBy string, ingredientTerms []string) ([]*Flavor, error)
	Insert(*Flavor) (*Flavor, error)
	Update(int, *Flavor) (*Flavor, error)
	Delete(int) (bool, error)
}

//go:generate counterfeiter . IngredientRepository
type IngredientRepository interface {
	Get(ID int64) (*Ingredient, error)
	GetByName(string) (*Ingredient, error)
	Insert(*Ingredient) (*Ingredient, error)
	Search(limit int, offset int, order string, search []string) ([]*Ingredient, error)
}
