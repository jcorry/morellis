package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:        1,
				FirstName: "Alice",
				LastName:  "Jones",
				Email:     "alice@example.com",
				Phone:     "867-5309",
				Status:    "Verified",
				Created:   time.Date(2019, 02, 24, 17, 25, 25, 0, time.UTC),
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			user, err := m.Get(tt.userID)

			if err != tt.wantError {
				t.Errorf("want %v, got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v, got %v", tt.wantUser, user)
			}
		})
	}
}
