package fpRowGroup

import (
	dbsql "database/sql"
	"entgo.io/ent/dialect/sql"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

type (
	Repo interface {
	}

	repoImpl struct {
		client *ent.Client
		sqlDB  *dbsql.DB
	}
)

const dbEngine = "mysql"

var _ Repo = &repoImpl{} // only for mention that repoImpl implement Repo

// NewRepo ...
func NewRepo(db *dbsql.DB) Repo {
	drv := sql.OpenDB(dbEngine, db)
	//drvDebug := dialect.Debug(drv)
	//client := ent.NewClient(ent.Driver(drvDebug))
	client := ent.NewClient(ent.Driver(drv))
	return &repoImpl{client: client, sqlDB: db}
}

// private function ---------------------------------------------------------------------------------------------
