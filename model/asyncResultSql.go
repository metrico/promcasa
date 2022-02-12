package model

import "github.com/jmoiron/sqlx"

type AsyncSqlResult struct {
	Rows       *sqlx.Rows
	QueryIndex int
	//sql.Result
	Err error
}
