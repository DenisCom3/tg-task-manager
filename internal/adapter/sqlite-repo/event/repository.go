package event

import "database/sql"

type Repository struct {
	conn *sql.DB
}

func newRepo(conn *sql.DB) Repository {
	return Repository{
		conn: conn,
	}
}