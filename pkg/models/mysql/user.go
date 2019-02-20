package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jcorry/morellis/pkg/models"
)

// UserModel wraps DB connection pool.
type UserModel struct {
	DB *sql.DB
}

const (
	DEFAULT_LIMIT int = 25
)

// Insert a new User
func (u *UserModel) Insert(firstName string, lastName string, email string, phone string) (*models.User, error) {
	created := time.Now()
	stmt := `INSERT INTO user (
		first_name,
		last_name,
		email,
		phone,
		status_id,
		created
	) VALUES (
		?,
		?,
		?,
		?,
		?,
		?
	)`
	result, err := u.DB.Exec(stmt, firstName, lastName, email, phone, models.USER_STATUS_UNVERIFIED, created)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Status:    models.USER_STATUS_UNVERIFIED.Slug(),
		Created:   created,
	}

	return user, nil
}

// Update a User identified by id
func (u *UserModel) Update(user *models.User) (*models.User, error) {
	stmt := `UPDATE user SET
			first_name = ?,
			last_name = ?,
			email = ?,
			phone = ?,
			status_id = ?
		WHERE id = ?`

	var userStatus models.UserStatus
	userStatus = models.USER_STATUS_VERIFIED
	userStatusID := userStatus.GetID(user.Status)

	_, err := u.DB.Exec(stmt, user.FirstName, user.LastName, user.Email, user.Phone, userStatusID, user.ID)
	if err != nil {
		return nil, err
	}

	user, err = u.Get(int(user.ID))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a single User by ID
func (u *UserModel) Get(id int) (*models.User, error) {
	stmt := `SELECT u.id, u.first_name, u.last_name, u.email, u.phone, s.name, u.created
			   FROM user AS u
		  LEFT JOIN ref_user_status AS s ON u.status_id = s.id
			  WHERE u.id = ?`

	user := &models.User{}
	err := u.DB.QueryRow(stmt, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Status, &user.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

// List Users limiting results by `limit` beginning at `offset` and ordered by `order`
func (u *UserModel) List(limit int, offset int, order string) ([]*models.User, error) {
	orderOpts := map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
		"email":     "email",
		"created":   "created",
		"status":    "s.name",
	}

	if val, ok := orderOpts[order]; ok {
		order = val
	} else {
		order = "created"
	}

	stmt := fmt.Sprintf(`SELECT u.id, first_name, last_name, email, phone, s.name, u.created
			   FROM user AS u
		  LEFT JOIN ref_user_status AS s ON u.status_id = s.id
		   ORDER BY %s
			  LIMIT ?,?`, order)

	if limit < 1 {
		limit = DEFAULT_LIMIT
	}

	rows, err := u.DB.Query(stmt, offset, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		u := &models.User{}
		err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.Status, &u.Created)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Delete the user identified by id.
func (u *UserModel) Delete(id int) (bool, error) {
	stmt, err := u.DB.Prepare(`DELETE FROM user WHERE id = ?`)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return false, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if affect > 0 {
		return true, nil
	}
	return false, nil
}

// Count gets the total count of user rows
func (u *UserModel) Count() int {
	var count int
	row := u.DB.QueryRow(`SELECT COUNT(*) FROM user`)

	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}

	return count
}