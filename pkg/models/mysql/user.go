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
	PW_HASH_COST  int = 12
)

// Insert a new User
func (u *UserModel) Insert(uid uuid.UUID, firstName string, lastName string, email string, phone string, statusID int, password string) (*models.User, error) {
	created := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), PW_HASH_COST)
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

	result, err := u.DB.Exec(stmt, uid.String(), firstName, lastName, email, phone, statusID, hashedPassword, created)
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

	userStatus := models.UserStatus(statusID)

	user := &models.User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Status:    userStatus.Slug(),
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
	stmt := `SELECT u.id, u.uuid, u.first_name, u.last_name, u.email, u.phone, s.slug, u.created, pu.id AS "pu_id", p.id AS "p_id", p.name
			   FROM user AS u
		  LEFT JOIN ref_user_status AS s ON u.status_id = s.id
		  	   JOIN permission_user AS pu ON pu.user_id = u.id
		       JOIN permission AS p ON pu.permission_id = p.id
			  WHERE u.uuid = ?`

	rows, err := u.DB.Query(stmt, uuid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &models.User{}

	for rows.Next() {
		p := &models.UserPermission{}

		err = rows.Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Status, &user.Created, &p.UserPermissionID, &p.ID, &p.Name)

		if err != nil {
			return nil, err
		}
		user.Permissions = append(user.Permissions, *p)
	}

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
	tx, _ := u.DB.Begin()

	stmt := `DELETE FROM permission_user WHERE user_id = ?`
	_, err := u.DB.Exec(stmt, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	stmt = `DELETE FROM user WHERE id = ?`

	res, err := u.DB.Exec(stmt, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if affect > 0 {
		tx.Commit()
		return true, nil
	}

	tx.Rollback()
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

func (u *UserModel) GetPermissions(userID int) ([]models.UserPermission, error) {
	var userPermissions []models.UserPermission

	stmt := `SELECT pu.id AS "userPermissionId", p.id, p.name
			   FROM permission AS p
		  LEFT JOIN permission_user AS pu ON pu.permission_id = p.id
		  LEFT JOIN user AS u ON pu.user_id = u.id
			  WHERE u.id = ?
		   ORDER BY name DESC`

	rows, err := u.DB.Query(stmt, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Permission
		var up models.UserPermission

		err = rows.Scan(&up.UserPermissionID, &p.ID, &p.Name)
		if err != nil {
			return nil, err
		}
		up.Permission = p
		userPermissions = append(userPermissions, up)
	}

	return userPermissions, nil
}

// AddPermission adds a Permission to a User
func (u *UserModel) AddPermission(userID int, p models.Permission) (int, error) {
	if !u.checkValidPermission(p) {
		return 0, models.ErrInvalidPermission
	}

	if !u.checkValidUser(userID) {
		return 0, models.ErrInvalidUser
	}

	stmt := `INSERT INTO permission_user (user_id, permission_id, created)
				  VALUES (?, (SELECT id FROM permission WHERE name = ? LIMIT 1), CURRENT_TIMESTAMP)`

	res, err := u.DB.Exec(stmt, userID, p.Name)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "uk_permission_user_permission_id_user_id") {
				return 0, models.ErrDuplicateUserPermission
			}
		}
		return 0, err
	}

	a, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if a == 0 {
		return 0, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// RemovePermission removes a Permission from a User
func (u *UserModel) RemovePermission(userPermissionID int) (bool, error) {
	stmt := `DELETE FROM permission_user 
	  			   WHERE id = ?`

	res, err := u.DB.Exec(stmt, userPermissionID)

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

// RemoveAllPermissions removes all Permissions from a User
func (u *UserModel) RemoveAllPermissions(userID int) error {
	stmt, err := u.DB.Prepare(`DELETE FROM permission_user WHERE user_id = ?`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID)
	return err
}

// UpdatePermissions replaces all of a Users Permissions with new Permissions
func (u *UserModel) UpdatePermissions(userID int, permissions []*models.UserPermission) error {
	var IDs []int
	for _, p := range permissions {
		IDs = append(IDs, p.Permission.ID)
	}

	tx, _ := u.DB.Begin()
	defer tx.Rollback()

	stmt := `DELETE FROM permission_user WHERE user_id = ? AND permission_id NOT IN (?` + strings.Repeat(", ?", len(IDs)-1) + `)`
	_, err := u.DB.Exec(stmt, userID, IDs)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, id := range IDs {
		stmt := `INSERT INTO permission_user (user_id, permission_id)
			 SELECT ?, ? 
			 FROM DUAL
			 WHERE NOT EXISTS (
				SELECT 1
				FROM permission_user
				WHERE user_id = ? AND permission_id = ?
			 )
			LIMIT 1`

		_, err = u.DB.Exec(stmt, userID, id, userID, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// AddIngredient creates a UserIngredient association. This is used for allowing Users to
// save Ingredient preferences for notifications.
func (u *UserModel) AddIngredient(userID int64, ingredient *models.Ingredient, keyword string) (*models.UserIngredient, error) {
	stmt := `INSERT INTO ingredient_user (ingredient_id, user_id, keyword) VALUES (?, ?, ?)`
	res, err := u.DB.Exec(stmt, ingredient.ID, userID, keyword)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "uk_ingredient_user_ingredient") {
				return nil, models.ErrDuplicateUserIngredient
			}
		}
		return nil, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var userIngredient = &models.UserIngredient{
		UserIngredientID: lastInsertId,
		Ingredient:       ingredient,
	}

	stmt = `SELECT id, created 
			   FROM ingredient_user 
			  WHERE id = ?`

	err = u.DB.QueryRow(stmt, lastInsertId).Scan(&userIngredient.UserIngredientID, &userIngredient.Created)
	if err != nil {
		return nil, err
	}

	return userIngredient, nil
}

// GetIngredients gets all of the UserIngredient associations for the User
func (u *UserModel) GetIngredients(userID int64) ([]*models.UserIngredient, error) {
	stmt := `SELECT iu.id, iu.created, i.id, i.name
			   FROM ingredient_user iu
	      LEFT JOIN ingredient i ON iu.ingredient_id = i.id
              WHERE user_id = ? AND deleted = 0`

	rows, err := u.DB.Query(stmt, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	var userIngredients []*models.UserIngredient

	for rows.Next() {
		i := &models.Ingredient{}
		ui := &models.UserIngredient{}

		err = rows.Scan(&ui.UserIngredientID, &ui.Created, &i.ID, &i.Name)

		if err != nil {
			return nil, err
		}

		ui.Ingredient = i
		userIngredients = append(userIngredients, ui)
	}

	return userIngredients, nil
}

// RemoveIngredient removes the UserIngredient association
func (u *UserModel) RemoveUserIngredient(userIngredientID int64) error {
	stmt := `UPDATE ingredient_user 
				SET deleted = ?
			  WHERE id = ? 
			    AND deleted = 0`

	res, err := u.DB.Exec(stmt, int32(time.Now().Unix()), userIngredientID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return models.ErrNoneAffected
	}

	return nil
}

func (u *UserModel) checkValidPermission(p models.Permission) bool {
	var isValid bool
	stmt := `SELECT IF(COUNT(*), 'true', 'false') 
			   FROM permission 
			  WHERE name = ?`

	err := u.DB.QueryRow(stmt, p.Name).Scan(&isValid)

	if err != nil {
		return false
	}

	return isValid
}

func (u *UserModel) checkValidUser(userID int) bool {
	var isValid bool
	stmt := `SELECT IF(COUNT(*), 'true', 'false')
			   FROM user
			  WHERE id = ? AND status_id = ?`

	err := u.DB.QueryRow(stmt, userID, models.USER_STATUS_VERIFIED).Scan(&isValid)
	if err != nil {
		return false
	}

	return isValid
}
