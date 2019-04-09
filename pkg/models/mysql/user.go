package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/go-sql-driver/mysql"

	"github.com/jcorry/morellis/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// UserModel wraps DB connection pool.
type UserModel struct {
	DB *sql.DB
}

const (
	DEFAULT_LIMIT int = 25
)

// Insert a new User
func (u *UserModel) Insert(uid uuid.UUID, firstName string, lastName string, email string, phone string, password string) (*models.User, error) {
	created := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	stmt := `INSERT INTO user (
        uuid,
		first_name,
		last_name,
		email,
		phone,
		status_id,
		hashed_password,
		created
	) VALUES (
	    ?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)`

	result, err := u.DB.Exec(stmt, uid.String(), firstName, lastName, email, phone, models.USER_STATUS_UNVERIFIED, hashedPassword, created)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "uk_user_email") {
				return nil, models.ErrDuplicateEmail
			}
		}
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
	userStatusID := userStatus.GetID(user.Status)

	_, err := u.DB.Exec(stmt, user.FirstName, user.LastName, user.Email, user.Phone, userStatusID, user.ID)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "uk_user_email") {
				return nil, models.ErrDuplicateEmail
			}
		}
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
	stmt := `SELECT u.id, u.uuid, u.first_name, u.last_name, u.email, u.phone, s.slug, u.created
			   FROM user AS u
		  LEFT JOIN ref_user_status AS s ON u.status_id = s.id
			  WHERE u.id = ?`

	user := &models.User{}
	err := u.DB.QueryRow(stmt, id).Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Status, &user.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a single User by UUID
func (u *UserModel) GetByUUID(uuid uuid.UUID) (*models.User, error) {
	stmt := `SELECT u.id, u.uuid, u.first_name, u.last_name, u.email, u.phone, s.slug, u.created
			   FROM user AS u
		  LEFT JOIN ref_user_status AS s ON u.status_id = s.id
			  WHERE u.uuid = ?`

	user := &models.User{}
	err := u.DB.QueryRow(stmt, uuid).Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Status, &user.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserModel) GetByCredentials(c models.Credentials) (*models.User, error) {
	var pwHash []byte = nil

	stmt := `SELECT u.id, u.uuid, u.first_name, u.last_name, u.email, u.hashed_password, u.phone, s.slug, u.created
			   FROM user AS u
		  LEFT JOIN ref_user_status AS s ON u.status_id = s.id
			  WHERE u.email = ?`

	user := &models.User{}

	err := u.DB.QueryRow(stmt, c.Email).Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &pwHash, &user.Phone, &user.Status, &user.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(pwHash, []byte(c.Password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, models.ErrInvalidCredentials
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

	stmt := fmt.Sprintf(`SELECT u.id, u.uuid, first_name, last_name, email, phone, s.slug, u.created
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
		err = rows.Scan(&u.ID, &u.UUID, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.Status, &u.Created)
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

func (u *UserModel) GetPermissions(ID int) ([]models.Permission, error) {
	var permissions []models.Permission

	stmt := `SELECT p.name
			   FROM permission AS p
		  LEFT JOIN permission_user AS pu ON pu.permission_id = p.id
		  LEFT JOIN user AS u ON pu.user_id = u.id
			  WHERE u.id = ?
		   ORDER BY name DESC`

	rows, err := u.DB.Query(stmt, ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Permission
		err = rows.Scan(&p)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

// AddPermission adds a Permission to a User
func (u *UserModel) AddPermission(userID int, p models.Permission) (bool, error) {
	if !u.checkValidPermission(p) {
		return false, models.ErrInvalidPermission
	}

	if !u.checkValidUser(userID) {
		return false, models.ErrInvalidUser
	}

	stmt := `INSERT INTO permission_user (user_id, permission_id, created)
				  VALUES (?, (SELECT id FROM permission WHERE name = ?), CURRENT_TIMESTAMP)`

	res, err := u.DB.Exec(stmt, userID, p)

	if err != nil {
		return false, err
	}

	a, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if a == 0 {
		return false, nil
	}

	return true, nil
}

// RemovePermission removes a Permission from a User
func (u *UserModel) RemovePermission(userID int, p models.Permission) (bool, error) {
	stmt := `DELETE FROM permission_user 
	  			   WHERE user_id = ? 
				     AND permission_id = (
						SELECT id FROM permission WHERE name = ?
					 )`
	res, err := u.DB.Exec(stmt, userID, p)

	if err != nil {
		return false, err
	}

	a, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	if a == 0 {
		return false, nil
	}

	return true, nil
}

func (u *UserModel) checkValidPermission(p models.Permission) bool {
	var isValid bool
	stmt := `SELECT IF(COUNT(*), 'true', 'false') 
			   FROM permission 
			  WHERE name = ?`

	err := u.DB.QueryRow(stmt, p).Scan(&isValid)

	if err != nil {
		return false
	}

	return isValid
}

func (u *UserModel) checkValidUser(userID int) bool {
	var isValid bool
	stmt := `SELECT IF(COUNT(*), 'true', 'false')
			   FROM user
			  WHERE id = ? 
		  		AND status_id = ?`

	err := u.DB.QueryRow(stmt, userID, models.USER_STATUS_VERIFIED).Scan(&isValid)
	if err != nil {
		return false
	}

	return isValid
}
