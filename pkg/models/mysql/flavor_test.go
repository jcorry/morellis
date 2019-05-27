package mysql

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jcorry/morellis/pkg/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestFlavorModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	cols := []string{"id", "name", "description", "created", "id", "name"}
	created := time.Now()
	tests := []struct {
		name    string
		id      int64
		rows    *sqlmock.Rows
		sqlErr  error
		wantErr error
	}{
		{
			"Success",
			1,
			sqlmock.NewRows(cols).AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 12, "vanilla").AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 13, "cream"),
			nil,
			nil,
		},
		{
			"No rows",
			1,
			nil,
			sql.ErrNoRows,
			models.ErrNoRecord,
		},
		{
			"Err rows",
			1,
			sqlmock.NewRows(cols).AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 12, "vanilla").AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 13, "cream").RowError(1, fmt.Errorf("row error")),
			nil,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := `^SELECT f.id, f.name, f.description, f.created, i.id, i.name
			   FROM flavor AS f
		  LEFT JOIN flavor_ingredient AS fi ON f.id = fi.flavor_id
		  LEFT JOIN ingredient AS i ON i.id = fi.ingredient_id
			  WHERE f.id = (.+)$`

			if tt.wantErr == nil {
				mock.ExpectQuery(query).WithArgs(tt.id).WillReturnRows(tt.rows)
			} else {
				mock.ExpectQuery(query).WithArgs(tt.id).WillReturnError(tt.sqlErr)
			}

			f := FlavorModel{DB: db}

			_, err := f.Get(int(tt.id))
			if tt.wantErr == sql.ErrNoRows {
				if err != models.ErrNoRecord {
					t.Errorf("Got unexpected error: %s", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFlavorModel_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	cols := []string{"id", "name", "description", "created", "id", "name"}
	created := time.Now()

	tests := []struct {
		name     string
		limit    int
		offset   int
		order    string
		wantRows *sqlmock.Rows
		wantErr  error
	}{
		{
			"Success",
			10,
			0,
			"",
			sqlmock.NewRows(cols).AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 12, "vanilla").AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 13, "cream"),
			nil,
		},
		{
			"Err rows",
			10,
			0,
			"",
			sqlmock.NewRows(cols).AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 12, "vanilla").AddRow(1, "Vanilla", "Smooth, creamy vanilla", created, 13, "cream").RowError(1, fmt.Errorf("row error")),
			fmt.Errorf("row error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := fmt.Sprintf(`^SELECT f.id, f.name, f.description, f.created, i.id, i.name
			   FROM flavor AS f
		       JOIN flavor_ingredient AS fi ON f.id = fi.flavor_id
		  LEFT JOIN ingredient AS i ON i.id = fi.ingredient_id
		   ORDER BY %s
			  LIMIT (.+), (.+)$`, "f.name")
			if tt.wantErr == nil {
				mock.ExpectQuery(query).WithArgs(tt.offset, tt.limit).WillReturnRows(tt.wantRows)
			} else {
				mock.ExpectQuery(query).WithArgs(tt.offset, tt.limit).WillReturnError(tt.wantErr)
			}

			f := FlavorModel{DB: db}

			_, err := f.List(tt.limit, tt.offset, tt.order)
			if err != tt.wantErr {
				t.Errorf("Got unexpected error, want %s; Got %s", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFlavorModel_Insert_ShouldCommit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	flavorID := int64(100)

	flavor := &models.Flavor{
		Name: "Vanilla",
		Ingredients: []models.Ingredient{
			{
				Name: "vanilla",
			},
			{
				Name: "sugar",
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec(`^INSERT INTO flavor \(name, description, created\) VALUES \((.+), (.+), (.+)\)$`).
		WillReturnResult(sqlmock.NewResult(flavorID, 1))

	getByNameQuery := `^SELECT id, name FROM ingredient WHERE LOWER\(name\) = (.+)$`
	insertIngredientQuery := `^INSERT INTO flavor_ingredient \(flavor_id, ingredient_id\) VALUES \((.+), (.+)\)$`
	for idx, i := range flavor.Ingredients {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(idx, i.Name)
		mock.ExpectQuery(getByNameQuery).WithArgs(i.Name).WillReturnRows(rows)

		mock.ExpectExec(insertIngredientQuery).
			WithArgs(flavorID, idx).
			WillReturnResult(sqlmock.NewResult(int64(100+idx), 1))
	}

	mock.ExpectCommit()

	f := FlavorModel{DB: db}

	_, err = f.Insert(flavor)

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestFlavorModel_Insert_ShouldRollbackOnFlavorInsertFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	insertError := fmt.Errorf("insert fail")

	flavor := &models.Flavor{
		Name: "Vanilla",
		Ingredients: []models.Ingredient{
			{
				Name: "vanilla",
			},
			{
				Name: "sugar",
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec(`^INSERT INTO flavor \(name, description, created\) VALUES \((.+), (.+), (.+)\)$`).
		WillReturnError(insertError)

	mock.ExpectRollback()

	f := FlavorModel{DB: db}

	_, err = f.Insert(flavor)

	if err != insertError {
		t.Errorf("Unexpected err: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestFlavorModel_Insert_ShouldRollbackOnIngredientInsertFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	flavorID := int64(100)

	flavor := &models.Flavor{
		Name: "Vanilla",
		Ingredients: []models.Ingredient{
			{
				Name: "vanilla",
			},
			{
				Name: "sugar",
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec(`^INSERT INTO flavor \(name, description, created\) VALUES \((.+), (.+), (.+)\)$`).
		WillReturnResult(sqlmock.NewResult(flavorID, 1))

	getByNameQuery := `^SELECT id, name FROM ingredient WHERE LOWER\(name\) = (.+)$`
	insertIngredientQuery := `^INSERT INTO flavor_ingredient \(flavor_id, ingredient_id\) VALUES \((.+), (.+)\)$`
	for idx, i := range flavor.Ingredients {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(idx, i.Name)
		mock.ExpectQuery(getByNameQuery).WithArgs(i.Name).WillReturnRows(rows)

		mock.ExpectExec(insertIngredientQuery).
			WithArgs(flavorID, idx).
			WillReturnResult(sqlmock.NewResult(int64(100+idx), 1))
	}

	mock.ExpectCommit()

	f := FlavorModel{DB: db}

	_, err = f.Insert(flavor)

	if err != nil {
		t.Errorf("Unexpected err: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestFlavorModel_Count(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	cols := []string{`COUNT(id)`}

	tests := []struct {
		name      string
		wantCount int
		wantRow   *sqlmock.Rows
		wantErr   error
	}{
		{
			"Success",
			12,
			sqlmock.NewRows(cols).AddRow(12),
			nil,
		},
		{
			"Error",
			0,
			nil,
			fmt.Errorf("Scan err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := `^SELECT COUNT\(id\) FROM flavor WHERE 1$`
			if tt.wantErr == nil {
				mock.ExpectQuery(query).WillReturnRows(tt.wantRow)
			} else {
				mock.ExpectQuery(query).WillReturnError(tt.wantErr)
			}

			f := FlavorModel{DB: db}

			count := f.Count()
			if count != tt.wantCount {
				t.Errorf("Unexpected count want %d; Got %d", tt.wantCount, count)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
