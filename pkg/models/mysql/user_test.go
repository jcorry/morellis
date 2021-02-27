package mysql_test

import (
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/jcorry/morellis/pkg/models"
	repo "github.com/jcorry/morellis/pkg/models/mysql"
)

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		phone string
		exp   string
	}{
		{
			phone: "+16785928804",
			exp:   "16785928804",
		},
		{
			phone: "(678)592-8804",
			exp:   "6785928804",
		},
		{
			phone: "678.592.8804",
			exp:   "6785928804",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.phone, func(t *testing.T) {
			require.Equal(t, tt.exp, repo.NormalizePhone(tt.phone))
		})
	}
}

func TestUserModel(t *testing.T) {
	db := repo.NewTestDB(t)
	rdb := repo.NewTestRedis(t)

	var r models.UserRepository

	r = &repo.UserModel{
		DB:    db,
		Redis: rdb,
	}

	var user *models.User
	var err error

	t.Run("insert a user", func(t *testing.T) {
		user, err = r.Insert(uuid.New(), models.NullString{String: "Testy"}, models.NullString{String: "McTestFace"}, models.NullString{String: "testy@testface.com"}, "404-533-1212", int(models.USER_STATUS_VERIFIED), "password")
		require.NoError(t, err)
	})

	t.Run("insert a user, same phone", func(t *testing.T) {
		_, err := r.Insert(uuid.New(), models.NullString{String: "Testy"}, models.NullString{String: "McTestFace"}, models.NullString{String: "testy@testface.com"}, "404-533-1212", int(models.USER_STATUS_VERIFIED), "password")
		require.EqualError(t, err, models.ErrDuplicatePhone.Error())
	})

	t.Run("insert a user, same email", func(t *testing.T) {
		_, err := r.Insert(uuid.New(), models.NullString{String: "Testy"}, models.NullString{String: "McTestFace"}, models.NullString{String: "testy@testface.com"}, "404-533-1010", int(models.USER_STATUS_VERIFIED), "password")
		require.EqualError(t, err, models.ErrDuplicateEmail.Error())
	})

	t.Run("get the user", func(t *testing.T) {
		user, err = r.Get(int(user.ID))
		require.NoError(t, err)
	})

	t.Run("get by phone", func(t *testing.T) {
		_, err := r.GetByPhone(user.Phone)
		require.NoError(t, err)
	})

	var token string
	t.Run("generate auth token", func(t *testing.T) {
		guid := uuid.New()
		token = base64.StdEncoding.EncodeToString([]byte(guid.String()))

		err := r.SaveAuthToken(token, int(user.ID))
		require.NoError(t, err)
	})

	t.Run("get by auth token", func(t *testing.T) {
		u, err := r.GetByAuthToken(token)
		require.NoError(t, err)
		require.Equal(t, user, u)
	})

	t.Run("get user list", func(t *testing.T) {
		l, err := r.List(10, 0, "")
		require.NoError(t, err)
		require.Len(t, l, 1)
	})

	var user2 *models.User
	t.Run("insert another user", func(t *testing.T) {
		user2, err = r.Insert(uuid.New(), models.NullString{String: "Testy"}, models.NullString{String: "McTestFace"}, models.NullString{String: "another@testface.com"}, "404-535-2233", int(models.USER_STATUS_VERIFIED), "password")
		require.NoError(t, err)
	})

	t.Run("get user list", func(t *testing.T) {
		l, err := r.List(10, 0, "")
		require.NoError(t, err)
		require.Len(t, l, 2)
	})

	t.Run("delete a user", func(t *testing.T) {
		ok, err := r.Delete(int(user2.ID))
		require.NoError(t, err)
		require.True(t, ok)
	})
}
