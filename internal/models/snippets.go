package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// insert a new snippet to the db
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// sql statement we wan't to execute
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES(?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// use the Exec() method on the connection pool to execute the statment
	// parameters: 1. sql statement, 2. arguments for placeholder parameters in the query (title, content, expires)
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// get the ID of the newly inserted record in the snippets table (DB)
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// our function returns an int and an error, the id is of type int64, so we need to convert it
	return int(id), nil
}

// return a specific snippet based on its id

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
