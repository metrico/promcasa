package service

import (
	"github.com/jmoiron/sqlx"
)

// Service : here you tell us what Salutation is
type ServiceData struct {
	Session []*sqlx.DB
}

//ServiceConfig
type ServiceConfig struct {
	Session *sqlx.DB
}

//ServiceConfigDatabases
type ServiceConfigDatabases struct {
	Session map[string]*sqlx.DB
}
