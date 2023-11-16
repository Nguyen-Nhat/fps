package common

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/xo/dburl"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/enttest"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

func CreateEntClientFromDB(db *sql.DB) *ent.Client {
	drv := entsql.OpenDB("mysql", db)
	return ent.NewClient(ent.Driver(drv))
}

func PrepareDatabase(ctx context.Context) *sql.DB {
	// 1. Set `RUN_PROFILE=TEST` for testing
	_ = os.Setenv(config.EnvKeyRunProfile, config.ProfileTest)

	// 2. Init DB Connection, DB Client
	dbConf := config.Load("../../..").Database.MySQL
	db, _ := dburl.Open(dbConf.DatabaseURI()) // no handle error, if error test will be terminated
	entClient := CreateEntClientFromDB(db)

	// 3. Mock data
	clearDataDbAndInsertMockData(ctx, db, entClient)

	return db
}

func PrepareDatabaseSqlite(ctx context.Context, t *testing.T) (*sql.DB, *ent.Client) {
	// 1. Set `RUN_PROFILE=TEST` for testing
	_ = os.Setenv(config.EnvKeyRunProfile, config.ProfileTest)

	// 2. Init DB Connection, DB Client
	db, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	assert.NoError(t, err)
	drv := entsql.OpenDB("sqlite3", db)
	// Migration schema
	entClient := enttest.NewClient(t, enttest.WithOptions(ent.Driver(drv)))

	// 3. Mock data
	clearDataDbAndInsertMockData(ctx, db, entClient)

	return db, entClient
}

func clearDataDbAndInsertMockData(ctx context.Context, db *sql.DB, entClient *ent.Client) {
	// 3. Drop tables in DB
	_, _ = db.ExecContext(ctx, "DROP TABLE users")
	_, _ = db.ExecContext(ctx, "DROP TABLE processing_file")
	_, _ = db.ExecContext(ctx, "DROP TABLE processing_file_row")
	_, _ = db.ExecContext(ctx, "DROP TABLE fps_client")
	_, _ = db.ExecContext(ctx, "DROP TABLE config_mapping")
	_, _ = db.ExecContext(ctx, "DROP TABLE config_task")

	// 4. Migration DB Schema
	if err := entClient.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed Creating Schema Resources: %v", err)
	}

	// 5. Mocking data to database
	mockProcessingFile(ctx, entClient)
	mockProcessingFileRow(ctx, entClient)
	mockXXX(ctx, entClient) // will mock data for other models

	fmt.Println() // new line in console => don't care about it
}

func mockProcessingFile(ctx context.Context, dbClient *ent.Client) {
	logger.Infof("Mock Processing File ...")
	_, err := fileprocessing.SaveAll(ctx, dbClient, processingFiles, false)
	if err != nil {
		logger.Errorf("Mock Processing File ... Failed: %v", err)
		panic(err)
	}
	logger.Infof("Mock Processing File ... Finished")
}

func mockProcessingFileRow(ctx context.Context, dbClient *ent.Client) {
	logger.Infof("Mock Processing File Row...")
	_, err := fileprocessingrow.SaveAll(ctx, dbClient, processingFileRows, false)
	if err != nil {
		logger.Errorf("Mock Processing File Row ... Failed: %v", err)
		panic(err)
	}
	logger.Infof("Mock Processing File Row... Finished")
}

func mockXXX(_ context.Context, _ *ent.Client) {
	logger.Infof("Mock XXX ...")
	// doSth ...
	logger.Infof("Mock XXX ... Finished")
}
