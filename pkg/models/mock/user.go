package mock

import (
	"time"

	"github.com/google/uuid"

	"github.com/jcorry/morellis/pkg/models"
)

var mockUser = &models.User{
	ID:        1,
	FirstName: models.NullString{String: "Testy", Valid: true},
	LastName:  models.NullString{String: "McTestFace", Valid: true},
	Email:     models.NullString{String: "test@example.com", Valid: true},
	Phone:     "867-5309",
	Status:    models.USER_STATUS_VERIFIED.Slug(),
	Created:   time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(uid uuid.UUID, firstName models.NullString, lastName models.NullString, email models.NullString, phone string, statusID int, password string) (*models.User, error) {
	user := &models.User{
		ID:        1,
		UUID:      uid,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Status:    models.USER_STATUS_VERIFIED.Slug(),
		Created:   time.Now(),
	}
	switch email.String {
	case "dupe@example.com":
		return nil, models.ErrDuplicateEmail
	default:
		return user, nil
	}
}

func (m *UserModel) Update(user *models.User) (*models.User, error) {
	switch user.Email.String {
	case "dupe@example.com":
		return nil, models.ErrDuplicateEmail
	default:
		return user, nil
	}
}

func (m *UserModel) Get(ID int) (*models.User, error) {
	if ID == 0 {
		return nil, models.ErrNoRecord
	}

	mockUser.ID = int64(ID)
	return mockUser, nil
}

func (m *UserModel) GetByUUID(ID uuid.UUID) (*models.User, error) {
	if ID.String() == "" || ID.String() == "e6fc6b5a-882c-40ba-b860-b11a413ec2df" {
		return nil, models.ErrNoRecord
	}

	if ID.String() == "df97802e-79e8-11e9-8f9e-2a86e4085a59" {
		mockUser.ID = 1001
	} else {
		mockUser.ID = 1
	}

	mockUser.UUID = ID
	mockUser.Permissions = []models.UserPermission{
		{
			Permission: models.Permission{Name: "user:read"},
		},
		{
			Permission: models.Permission{Name: "user:write"},
		},
		{
			Permission: models.Permission{Name: "self:read"},
		},
		{
			Permission: models.Permission{Name: "self:write"},
		},
	}

	return mockUser, nil
}

func (m *UserModel) GetByCredentials(credentials models.Credentials) (*models.User, error) {
	mockUser.Email.String = credentials.Email
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	mockUser.UUID = uid

	if credentials.Email == "noauth@example.com" {
		return nil, models.ErrInvalidCredentials
	}

	return mockUser, nil
}

func (m *UserModel) List(limit int, offset int, order string) ([]*models.User, error) {
	return nil, nil
}

func (m *UserModel) Delete(ID int) (bool, error) {
	return true, nil
}

func (m *UserModel) Count() int {
	return 4
}

func (u *UserModel) GetPermissions(ID int) ([]models.UserPermission, error) {
	return []models.UserPermission{
		{
			Permission: models.Permission{Name: "user:read"},
		},
		{
			Permission: models.Permission{Name: "user:write"},
		},
	}, nil
}

// AddPermission adds a Permission to a User
func (u *UserModel) AddPermission(userID int, p models.Permission) (int, error) {
	return 112, nil
}

// RemovePermission removes a Permission from a User
func (u *UserModel) RemovePermission(userPermissionID int) (bool, error) {
	return true, nil
}

func (u *UserModel) RemoveAllPermissions(userID int) error {
	return nil
}

// AddIngredient creates a UserIngredient association. This is used for allowing Users to
// save Ingredient preferences for notifications.
func (u *UserModel) AddIngredient(userID int64, ingredient *models.Ingredient, keyword string) (*models.UserIngredient, error) {
	if userID == 1001 {
		return nil, models.ErrDuplicateUserIngredient
	}

	userIngredient := &models.UserIngredient{
		UserIngredientID: 10,
		Ingredient:       ingredient,
		Created:          time.Now(),
	}

	return userIngredient, nil
}

// GetIngredients gets all of the UserIngredient associations for the User
func (u *UserModel) GetIngredients(userID int64) ([]*models.UserIngredient, error) {
	if userID >= 1000 {
		return nil, models.ErrNoRecord
	}

	uiSlice := []*models.UserIngredient{}

	ingredientNames := []string{"vanilla", "almonds", "jalapeno", "coconut", "chocolate"}

	for i := 0; i < 3; i++ {
		ui := &models.UserIngredient{
			UserIngredientID: int64(i),
			Ingredient: &models.Ingredient{
				ID:   int64(i * 3),
				Name: ingredientNames[i],
			},
			Created: time.Now(),
		}
		uiSlice = append(uiSlice, ui)
	}
	return uiSlice, nil
}

// RemoveIngredient removes the UserIngredient association
func (u *UserModel) RemoveUserIngredient(userIngredientID int64) error {
	if userIngredientID >= 1000 {
		return models.ErrNoneAffected
	}

	return nil
}
