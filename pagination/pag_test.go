package pagination

import (
	"database/sql"
	"errors"
	"testing"
	
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID   int
	Name string
}

func TestNewPagination(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	assert.Equal(t, 10, p.GetPageSize())
	assert.Equal(t, 10, p.PaginationAbstract.pageSize)
	assert.Equal(t, 1, p.PaginationAbstract.pageNum)
	assert.Equal(t, db, p.PaginationAbstract.conn)
}

func TestSetPageSize(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetPageSize(20)
	assert.Equal(t, 20, p.GetPageSize())
}

func TestSetPage(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetPage(3)
	assert.Equal(t, 3, p.PaginationAbstract.pageNum)
}

func TestSetRawQuery(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	query := "SELECT * FROM test"
	params := []interface{}{1, "test"}
	p.SetRawQuery(query, params...)
	assert.Equal(t, query, p.PaginationAbstract.rawQuery)
	assert.Equal(t, params, p.PaginationAbstract.params)
}

func TestSetOrder(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetOrder("id DESC", "name ASC")
	assert.Equal(t, " order by id DESC,name ASC", p.PaginationAbstract.orderBy)
}

func TestCalculateOffset(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	assert.Equal(t, 0, p.calculateOffset())
	
	p.SetPage(2)
	assert.Equal(t, 10, p.calculateOffset())
	
	p.SetPage(3)
	assert.Equal(t, 20, p.calculateOffset())
}

func TestGetQuery(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetRawQuery("SELECT * FROM test")
	p.SetOrder("id DESC")
	
	expected := "SELECT * FROM test order by id DESC limit 10 offset 0"
	assert.Equal(t, expected, p.GetQuery())
}

func TestNextPage(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetTotalCount(30)
	
	// Test next page exists
	nextPage := p.nextPage()
	assert.NotNil(t, nextPage)
	assert.Equal(t, 2, *nextPage)
	
	// Test no next page
	p.SetPage(3)
	nextPage = p.nextPage()
	assert.Nil(t, nextPage)
}

func TestPreviousPage(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	
	// Test no previous page on page 1
	prevPage := p.previousPage()
	assert.Nil(t, prevPage)
	
	// Test previous page exists
	p.SetPage(2)
	prevPage = p.previousPage()
	assert.NotNil(t, prevPage)
	assert.Equal(t, 1, *prevPage)
}

func TestSetCountByQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetRawQuery("SELECT * FROM test")
	
	// Mock the count query
	rows := sqlmock.NewRows([]string{"count"}).AddRow(25)
	mock.ExpectQuery("select count\\(\\*\\) from \\(SELECT \\* FROM test\\)").WillReturnRows(rows)
	
	err = p.SetCountByQuery()
	assert.NoError(t, err)
	assert.Equal(t, 25, p.PaginationAbstract.totalCount)
	
	// Test error case
	mock.ExpectQuery("select count\\(\\*\\) from \\(SELECT \\* FROM test\\)").WillReturnError(errors.New("count error"))
	err = p.SetCountByQuery()
	assert.Error(t, err)
}

func TestRunSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetRawQuery("SELECT id, name FROM test")
	p.SetTotalCount(10)
	
	// Mock the query
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Test 1").
		AddRow(2, "Test 2")
	mock.ExpectQuery("SELECT id, name FROM test limit 10 offset 0").WillReturnRows(rows)
	
	scanFunc := func(e *[]TestModel, rows *sql.Rows) error {
		for rows.Next() {
			var m TestModel
			if err := rows.Scan(&m.ID, &m.Name); err != nil {
				return err
			}
			*e = append(*e, m)
		}
		return nil
	}
	
	err = p.runSQL(scanFunc)
	assert.NoError(t, err)
	assert.Len(t, p.PaginationAbstract.data, 2)
	assert.Equal(t, 1, p.PaginationAbstract.data[0].ID)
	assert.Equal(t, "Test 1", p.PaginationAbstract.data[0].Name)
	assert.Equal(t, 2, p.PaginationAbstract.data[1].ID)
	assert.Equal(t, "Test 2", p.PaginationAbstract.data[1].Name)
}

func TestJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	
	p := New[TestModel](10, 1, db)
	p.SetRawQuery("SELECT id, name FROM test")
	
	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(25)
	mock.ExpectQuery("select count\\(\\*\\) from \\(SELECT id, name FROM test\\)").WillReturnRows(countRows)
	
	// Mock data query
	dataRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Test 1").
		AddRow(2, "Test 2")
	mock.ExpectQuery("SELECT id, name FROM test limit 10 offset 0").WillReturnRows(dataRows)
	
	scanFunc := func(e *[]TestModel, rows *sql.Rows) error {
		for rows.Next() {
			var m TestModel
			if err := rows.Scan(&m.ID, &m.Name); err != nil {
				return err
			}
			*e = append(*e, m)
		}
		return nil
	}
	
	envelope, err := p.JSON(scanFunc)
	assert.NoError(t, err)
	assert.NotNil(t, envelope)
	
	data, ok := envelope["data"].([]TestModel)
	assert.True(t, ok)
	assert.Len(t, data, 2)
	
	count, ok := envelope["count"].(int)
	assert.True(t, ok)
	assert.Equal(t, 25, count)
	
	nextPage, ok := envelope["next_page"].(*int)
	assert.True(t, ok)
	assert.NotNil(t, nextPage)
	assert.Equal(t, 2, *nextPage)
	
	prevPage, ok := envelope["previous_page"].(*int)
	assert.True(t, ok)
	assert.Nil(t, prevPage)
}

func TestValidate(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	t.Run("valid page number", func(t *testing.T) {
		p := New[TestModel](10, 1, db)
		p.SetTotalCount(15)
		err := p.validate()
		assert.NoError(t, err)
	})
	
	t.Run("invalid page number", func(t *testing.T) {
		p := New[TestModel](10, 0, db)
		p.SetTotalCount(15)
		err := p.validate()
		assert.Error(t, err)
		assert.Equal(t, "pagination pageNum cannot be less than 1", err.Error())
	})
	
	t.Run("page number exceeds total count", func(t *testing.T) {
		p := New[TestModel](10, 3, db)
		p.SetTotalCount(15)
		err := p.validate()
		assert.Error(t, err)
		assert.Equal(t, "pagination pageNum and totalCount must be greater than pageSize", err.Error())
	})
}
