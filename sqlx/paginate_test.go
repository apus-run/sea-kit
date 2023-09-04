package sqlx

import (
	"testing"
)

type Test struct {
	Name     string `json:"name"     db:"name"`
	LastName string `json:"lastName" db:"last_name"`
}

func TestPaginateSimple(t *testing.T) {
	pagin := Instance(Test{})
	query, queryCount := pagin.
		Query("SELECT t.* FROM test t").
		Sort([]string{"name", "last_name"}).
		Desc([]string{"true", "false"}).
		Page(3).
		RowsPerPage(50).
		SearchBy("vinicius", "t.id").
		Select()

	t.Log(*query)
	t.Log(*queryCount)
}

func TestPaginate(t *testing.T) {
	queryString := "SELECT t.* FROM test t WHERE 1=1 and ((t.id::TEXT ilike '%vinicius%') ) ORDER BY name DESC, last_name ASC  LIMIT 50 OFFSET 100;"
	queryCountString := "SELECT COUNT(t.id) FROM test t WHERE 1=1 and ((t.id::TEXT ilike '%vinicius%') ) "

	pagin := Instance(Test{})
	query, queryCount := pagin.
		Query("SELECT t.* FROM test t").
		Sort([]string{"name", "last_name"}).
		Desc([]string{"true", "false"}).
		Page(3).
		RowsPerPage(50).
		SearchBy("vinicius", "t.id").
		Select()

	// t.Log(*query)
	// t.Log(*queryCount)

	if queryString != *query {
		t.Errorf("Wrong query")
		return
	}

	if queryCountString != *queryCount {
		t.Errorf("Wrong query count")
	}
}

func TestPaginateWithArgs(t *testing.T) {
	queryString := "SELECT t.* FROM test t WHERE t.name = 'jhon' and ((t.id::TEXT ilike '%vinicius%') ) ORDER BY name DESC, last_name ASC  LIMIT 50 OFFSET 100;"
	queryCountString := "SELECT COUNT(t.id) FROM test t WHERE t.name = 'jhon' and ((t.id::TEXT ilike '%vinicius%') ) "

	pagin := Instance(Test{})

	pagin.Query("SELECT t.* FROM test t").
		Sort([]string{"name", "last_name"}).
		Desc([]string{"true", "false"}).
		Page(3).
		RowsPerPage(50)

	if true {
		pagin.WhereArgs("t.name = 'jhon'")
		pagin.SearchBy("vinicius", "t.id")
		query, queryCount := pagin.Select()
		t.Log(*query)
		t.Log(*queryCount)

		if queryString != *query {
			t.Errorf("Wrong query")
		}

		if queryCountString != *queryCount {
			t.Errorf("Wrong query count")
		}
	}

}
