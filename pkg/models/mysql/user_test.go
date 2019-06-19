package mysql

import (
	"database/sql"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/jcorry/morellis/pkg/models"
)

func TestUserModel_Get(t *testing.T) {
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
				FirstName: models.NullString{"Alice", true},
				LastName:  models.NullString{"Jones", true},
				Email:     models.NullString{"alice@example.com", true},
				Phone:     "867-5309",
				Status:    "verified",
				Created:   time.Date(2019, 02, 24, 17, 25, 25, 0, time.UTC),
			},
			wantError: nil,
		},
	}

	db, teardown := newTestDB(t)
	defer teardown()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := UserModel{db}

			user, err := m.Get(tt.userID)
			if err != nil {
				t.Errorf("Unexpected error getting user: %s", err)
			}

			// No way to generate this...has to come from DB
			tt.wantUser.UUID = user.UUID

			if err != tt.wantError {
				t.Errorf("want %v, got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v, got %v", tt.wantUser, user)
			}
		})
	}
}

func TestUserModel_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	result, err := m.Delete(1)
	if !result {
		t.Errorf("want true, got false")
	}

	if err != nil {
		t.Errorf("want nil err, got %s", err)
	}
}

func TestUserModel_List(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	m.Delete(1)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		t.Fatal("Failed to hash password")
	}

	users := []*models.User{
		{
			FirstName: models.NullString{"John", true},
			LastName:  models.NullString{"Corry", true},
			Email:     models.NullString{"jcorry@morellis.com", true},
			Phone:     "867-5309",
			Password:  string(hashedPassword),
		},
		{
			FirstName: models.NullString{"Garrett", true},
			LastName:  models.NullString{"Rap", true},
			Email:     models.NullString{"garrett@morellis.com", true},
			Phone:     "867-5309",
			Password:  string(hashedPassword),
		},
		{
			FirstName: models.NullString{"Brian", true},
			LastName:  models.NullString{"Morton", true},
			Email:     models.NullString{"brian@morellis.com", true},
			Phone:     "867-5309",
			Password:  string(hashedPassword),
		},
	}

	toD := []int64{}

	for _, u := range users {
		uid, err := uuid.NewRandom()
		if err != nil {
			t.Error("Failed to create UUID for user")
		}

		user, err := m.Insert(uid, u.FirstName, u.LastName, u.Email, u.Phone, int(models.USER_STATUS_VERIFIED), u.Password)
		if err != nil {
			t.Fatal("Failed to insert new user for test")
		}
		user.UUID = uid
		toD = append(toD, user.ID)
		time.Sleep(time.Millisecond * 1000)
	}

	tests := []struct {
		name          string
		order         string
		limit         int
		wantUserNames []string
		wantError     error
	}{
		{
			"Get all 3 users, no order",
			"created",
			3,
			[]string{"John", "Garrett", "Brian"},
			nil,
		},
		{
			"Get all 3 users, order by first name",
			"firstName",
			3,
			[]string{"Brian", "Garrett", "John"},
			nil,
		},
		{
			"Get 2 users, order by first name",
			"firstName",
			2,
			[]string{"Brian", "Garrett"},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := m.List(tt.limit, 0, tt.order)
			if err != tt.wantError {
				t.Errorf("want nil, got %s", err)
			}

			if len(list) != tt.limit {
				t.Errorf("want list length %d; got %d", tt.limit, len(list))
			}

			for i, u := range list {
				if u.FirstName.String != tt.wantUserNames[i] {
					t.Errorf("want %s; got %s", tt.wantUserNames[i], u.FirstName.String)
				}
			}

		})
	}

}

func TestUserModel_GetByUUID(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	u, err := m.Get(1)
	if err != nil {
		t.Error(err)
	}

	u, err = m.GetByUUID(u.UUID)
	if err != nil {
		t.Error(err)
	}
}

func TestUserModel_GetByCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		email     string
		password  string
		wantError error
	}{
		{
			name:      "Found user",
			email:     "alice@example.com",
			password:  "password",
			wantError: nil,
		},
		{
			name:      "Wrong email",
			email:     "bob@example.com",
			password:  "password",
			wantError: models.ErrInvalidCredentials,
		},
		{
			name:      "Wrong password",
			email:     "alice@example.com",
			password:  "P@SSW0rd",
			wantError: models.ErrInvalidCredentials,
		},
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := models.Credentials{
				Email:    tt.email,
				Password: tt.password,
			}

			_, err := m.GetByCredentials(c)
			if err != tt.wantError {
				t.Error(err)
			}
		})
	}
}

func TestUserModel_GetPermissions(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	tests := []struct {
		name     string
		minCount int
		wantErr  error
	}{
		{"With Permissions", 1, nil},
		{"No Permissions", 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				ID: 1,
			}

			if tt.minCount > 0 {
				stmt := `INSERT INTO permission_user (user_id, permission_id)
					SELECT ?, id FROM permission WHERE name LIKE "self%"`

				_, err := m.DB.Exec(stmt, user.ID)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Setup complete...test it!
			permissions, err := m.GetPermissions(int(user.ID))
			if err != tt.wantErr {
				t.Errorf("Want err %v; Got err %v", tt.wantErr, err)
			}

			if len(permissions) < tt.minCount {
				t.Errorf("Want at least %d items; Got %d items", tt.minCount, len(permissions))
			}

			if err != nil {
				t.Fatal(err)
			}

			stmt := `DELETE FROM permission_user WHERE user_id = ?`

			_, err = m.DB.Exec(stmt, user.ID)

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserModel_AddPermission(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	tests := []struct {
		name    string
		userID  int
		perm    models.Permission
		wantErr error
		wantRes int
	}{
		{"Valid Permission", 1, models.Permission{1, "flavor:read"}, nil, 1},
		{"Invalid Permission", 1, models.Permission{0, "foo:write"}, models.ErrInvalidPermission, 0},
		{"Invalid User", 100, models.Permission{0, "user:read"}, models.ErrInvalidUser, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := m.AddPermission(tt.userID, tt.perm)
			if err != tt.wantErr {
				t.Errorf("Want err %v; Got err %v", tt.wantErr, err)
			}

			if err != nil && !(res >= tt.wantRes) {
				t.Errorf("Want res %v; Got res %v", tt.wantRes, res)
			}
		})
	}
}

func TestUserModel_RemovePermission(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	tests := []struct {
		name       string
		userID     int
		permission models.Permission
		wantRes    bool
		wantErr    error
	}{
		{"Valid Permission", 1, models.Permission{1, "self:write"}, true, nil},
		{"Invalid Permission", 1, models.Permission{0, "self:foo"}, false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up by adding the permission
			var userPermissionID int
			if m.checkValidPermission(tt.permission) {
				res, err := m.AddPermission(tt.userID, tt.permission)
				if err != nil {
					t.Fatal(err)
				}
				if res == 0 {
					t.Fatalf("Permission %v was not inserted", tt.permission)
				}
				userPermissionID = res
			}

			// Setup complete! Test it out!
			res, err := m.RemovePermission(userPermissionID)
			if err != tt.wantErr {
				t.Errorf("Want err %v; Got err %v", tt.wantErr, err)
			}

			if res != tt.wantRes {
				t.Errorf("Want res %v; Got res %v", tt.wantRes, res)
			}
		})
	}
}

func TestUserModel_checkValidPermission(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	cols := []string{`IF (COUNT(*), 'true', 'false')`}

	m := UserModel{DB: db}

	tests := []struct {
		name       string
		permission string
		sqlErr     error
		wantRes    bool
	}{
		{
			"Valid permission",
			"user:read",
			nil,
			true,
		},
		{
			"Invalid permission",
			"foo:bar",
			nil,
			false,
		},
		{
			"SQL Error",
			"user:read",
			sql.ErrNoRows,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := sqlmock.NewRows(cols).AddRow(strconv.FormatBool(tt.wantRes))

			querySql := `^SELECT IF\(COUNT\(\*\), 'true', 'false'\) FROM permission WHERE name = (.+)$`

			if tt.sqlErr == nil {
				mock.ExpectQuery(querySql).WithArgs(tt.permission).WillReturnRows(row)
			} else {
				mock.ExpectQuery(querySql).WithArgs(tt.permission).WillReturnError(sql.ErrNoRows)
			}

			p := models.Permission{
				ID:   0,
				Name: tt.permission,
			}

			res := m.checkValidPermission(p)

			if res != tt.wantRes {
				t.Errorf("Unexpected result, got %s, want %s", strconv.FormatBool(res), strconv.FormatBool(tt.wantRes))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserModel_AddIngredient(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		execErr  *mysql.MySQLError
		queryErr error
		wantErr  error
	}{
		{
			"Successful insert",
			nil,
			nil,
			nil,
		},
		{
			"Duplicate err",
			&mysql.MySQLError{Number: 1062, Message: "...uk_ingredient_user_ingredient..."},
			nil,
			models.ErrDuplicateUserIngredient,
		},
		{
			"Unspecified err",
			&mysql.MySQLError{Number: 1045, Message: "Foo"},
			nil,
			&mysql.MySQLError{Number: 1045, Message: "Foo"},
		},
		{
			"Unable to select err",
			nil,
			sql.ErrNoRows,
			sql.ErrNoRows,
		},
	}

	for _, tt := range tests {

		userID := int64(1)
		ingredient := models.Ingredient{14, "Chocolate"}
		keyword := "choc"
		lastInsertId := int64(27)

		m := UserModel{DB: db}

		t.Run(tt.name, func(t *testing.T) {

			querySql := `^INSERT INTO ingredient_user \(ingredient_id, user_id, keyword\) VALUES ((.+), (.+), (.+))`
			if tt.execErr == nil {
				mock.ExpectExec(querySql).WithArgs(ingredient.ID, userID, keyword).WillReturnResult(sqlmock.NewResult(lastInsertId, 1))
			} else {
				mock.ExpectExec(querySql).WithArgs(ingredient.ID, userID, keyword).WillReturnError(tt.execErr)
			}

			if tt.execErr == nil {
				row := sqlmock.NewRows([]string{"id", "created"}).AddRow(lastInsertId, time.Now())
				querySql = `^SELECT id, created FROM ingredient_user WHERE id = (.+)$`

				if tt.queryErr == nil {
					mock.ExpectQuery(querySql).WithArgs(lastInsertId).WillReturnRows(row)
				} else {
					mock.ExpectQuery(querySql).WithArgs(lastInsertId).WillReturnError(sql.ErrNoRows)
				}

			}

			_, err := m.AddIngredient(userID, &ingredient, keyword)
			if tt.execErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Didn't get expected error; want %s; got %s.", tt.wantErr, err)
				}
			} else if tt.queryErr != nil {
				if err != tt.queryErr {
					t.Errorf("Didn't get expected error; want %s; got %s.", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}

}

func TestUserModel_GetIngredients(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	cols := []string{"id", "created", "id", "name"}

	tests := []struct {
		name     string
		userId   int64
		numRows  int
		wantErr  error
		wantRows *sqlmock.Rows
	}{
		{
			"2 rows",
			17,
			2,
			nil,
			sqlmock.NewRows(cols).AddRow(1, randomTimestamp(), 12, "caramel").AddRow(2, randomTimestamp(), 12, "fudge"),
		},
		{
			"Deleted",
			17,
			0,
			sql.ErrNoRows,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := `^SELECT iu.id, iu.created, i.id, i.name
			   FROM ingredient_user iu
	      LEFT JOIN ingredient i ON iu.ingredient_id = i.id
              WHERE user_id = (.+) AND deleted = 0$`

			if tt.wantErr == nil {
				mock.ExpectQuery(query).WithArgs(tt.userId).WillReturnRows(tt.wantRows)
			} else {
				mock.ExpectQuery(query).WithArgs(tt.userId).WillReturnError(tt.wantErr)
			}

			m := UserModel{DB: db}

			userIngedients, err := m.GetIngredients(tt.userId)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			if tt.wantErr == sql.ErrNoRows {
				if err != models.ErrNoRecord {
					t.Errorf("Got unexpected error: %s", err)
				}
			}

			if tt.numRows != len(userIngedients) {
				t.Errorf("Unexpected return length: want %d, got %d", tt.numRows, len(userIngedients))
			}

		})

	}
}

func TestUserModel_RemoveUserIngredient(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		affected      int64
		wantErr       error
		wantResultErr error
	}{
		{
			"Successful Delete",
			1,
			nil,
			nil,
		},
		{
			"Failed delete",
			1,
			mysql.ErrMalformPkt,
			mysql.ErrMalformPkt,
		},
		{
			"None Affected",
			0,
			nil,
			models.ErrNoneAffected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := `^UPDATE ingredient_user 
						  SET deleted = (.+)
					    WHERE id = (.+) 
						  AND deleted = 0$`

			userIngredientId := int64(7)

			if tt.wantErr == nil {
				mock.ExpectExec(query).WithArgs(AnyInt64{}, userIngredientId).WillReturnResult(sqlmock.NewResult(1, tt.affected))
			} else {
				mock.ExpectExec(query).WithArgs(AnyInt64{}, userIngredientId).WillReturnError(tt.wantErr)
			}

			m := UserModel{DB: db}

			err := m.RemoveUserIngredient(userIngredientId)
			if err != tt.wantResultErr {
				t.Errorf("Unexpected error: want %s; Got %s", tt.wantResultErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}

}
