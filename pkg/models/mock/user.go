package mock

import (
	"time"

	"github.com/google/uuid"

	"github.com/jcorry/morellis/pkg/models"
)

var mockUser = &models.User{
	ID:        1,
	FirstName: "Testy",
	LastName:  "McTestFace",
	Email:     "test@example.com",
	Phone:     "867-5309",
	Status:    models.USER_STATUS_VERIFIED.Slug(),
	Created:   time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(firstName string, lastName string, email string, phone string, password string) (*models.User, error) {
	user := &models.User{
		ID:        1,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Status:    models.USER_STATUS_VERIFIED.Slug(),
		Created:   time.Now(),
	}
	switch email {
	case "dupe@example.com":
		return nil, models.ErrDuplicateEmail
	default:
		return user, nil
	}
}

func (m *UserModel) Update(user *models.User) (*models.User, error) {
	switch user.Email {
	case "dupe@example.com":
		return nil, models.ErrDuplicateEmail
	default:
		return user, nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	if id == 0 {
		return nil, models.ErrNoRecord
	}

	mockUser.ID = int64(id)
	return mockUser, nil
}

func (m *UserModel) GetByUUID(id uuid.UUID) (*models.User, error) {
	if id.String() == "" {
		return nil, models.ErrNoRecord
	}

	mockUser.UUID = id
	return mockUser, nil
}

func (m *UserModel) List(limit int, offset int, order string) ([]*models.User, error) {
	return nil, nil
}

func (m *UserModel) Delete(id int) (bool, error) {
	return true, nil
}

func (m *UserModel) Count() int {
	return 4
}

func (m *UserModel) Authenticate(email string, password string) (*models.User, error) {
	mockUser.Email = email

	return mockUser, nil
}
