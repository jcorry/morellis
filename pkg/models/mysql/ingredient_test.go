package mysql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/jcorry/morellis/pkg/models"
)

func TestIngredientModel_GetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	cols := []string{"id", "name"}

	tests := []struct {
		name       string
		searchTerm string
		wantErr    error
		wantRows   *sqlmock.Rows
	}{
		{
			"Match exists", "Chocolate", nil, sqlmock.NewRows(cols).AddRow(1, "Chocolate"),
		},
		{
			"Match doesn't exist", "Vanilla", sql.ErrNoRows, nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == nil {
				mock.ExpectQuery(`^SELECT id, name FROM ingredient WHERE LOWER\(name\) = (.+)?`).WithArgs(tt.searchTerm).WillReturnRows(tt.wantRows)
			} else {
				mock.ExpectQuery(`^SELECT id, name FROM ingredient WHERE LOWER\(name\) = (.+)?`).WithArgs(tt.searchTerm).WillReturnError(tt.wantErr)
			}

			ing := IngredientModel{DB: db}

			ingredient, err := ing.GetByName(tt.searchTerm)
			if tt.wantErr == sql.ErrNoRows {
				if err != models.ErrNoRecord {
					t.Errorf("Got unexpexted error: %s", err)
				}
			}

			if ingredient != nil && ingredient.Name != tt.searchTerm {
				t.Errorf("Got unexpected Ingredient")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

		})
	}
}

func TestIngredientModel_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	cols := []string{"id", "name"}

	tests := []struct {
		name     string
		limit    int
		offset   int
		order    string
		search   []string
		wantErr  error
		wantRows *sqlmock.Rows
	}{
		{
			"Single term; match exists; sort by name",
			10,
			0,
			"name",
			[]string{"van"},
			nil,
			sqlmock.NewRows(cols).AddRow(1, "Vanilla"),
		},
		{
			"Multi term; match exists; sort by name",
			10,
			0,
			"name",
			[]string{"van", "coc"},
			nil,
			sqlmock.NewRows(cols).AddRow(1, "Vanilla").AddRow(13, "Coconut"),
		},
		{
			"Single term; match exists; sort by name",
			10,
			0,
			"name",
			[]string{"spiders"},
			sql.ErrNoRows,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querySql := `SELECT id, name FROM ingredient WHERE LOWER\(name\) LIKE (.+)`
			for i := 1; i <= len(tt.search)-1; i++ {
				querySql += ` OR LOWER\(name\) LIKE (.+)`
			}
			querySql += fmt.Sprintf(` ORDER BY %s`, tt.order)
			querySql += ` LIMIT (.+) OFFSET (.+)`
			querySql += `$`

			args := make([]driver.Value, len(tt.search))
			for i, term := range tt.search {
				args[i] = fmt.Sprintf(`%%%s%%`, term)
			}

			args = append(args, tt.limit, tt.offset)

			if tt.wantErr == nil {
				mock.ExpectQuery(querySql).WithArgs(args...).WillReturnRows(tt.wantRows)
			} else {
				mock.ExpectQuery(querySql).WithArgs(args...).WillReturnError(tt.wantErr)
			}

			ing := IngredientModel{DB: db}

			_, err := ing.Search(tt.limit, tt.offset, tt.order, tt.search)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected err: %s; got nil", tt.wantErr)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestIngredientModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub DB connection", err)
	}
	defer db.Close()

	tests := []struct {
		name           string
		ingredientName string
		wantID         int64
		wantErr        error
	}{
		{
			"No error insert",
			"Chocolate",
			17,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			querySql := `^INSERT INTO ingredient \(name, created\) VALUES \((.+), (.+)\)$`
			insertIngredient := models.Ingredient{Name: tt.ingredientName}

			args := []driver.Value{tt.ingredientName, sqlmock.AnyArg()}

			mock.ExpectExec(querySql).WithArgs(args...).WillReturnResult(sqlmock.NewResult(tt.wantID, 1))

			ingredientDao := IngredientModel{DB: db}

			ingredient, err := ingredientDao.Insert(&insertIngredient)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
			}

			if ingredient.ID != tt.wantID {
				t.Errorf("Mismatched IDs, want %d, got %d", tt.wantID, ingredient.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}

}
