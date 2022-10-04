package common

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileawardpoint"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/xo/dburl"
	"log"
	"os"
)

func PrepareDatabase(ctx context.Context) *sql.DB {
	// 1. Set `RUN_PROFILE=TEST` for testing
	_ = os.Setenv(config.EnvKeyRunProfile, config.ProfileTest)

	// 2. Init DB Connection, DB Client
	dbConf := config.Load("../..").Database.MySQL
	db, _ := dburl.Open(dbConf.DatabaseURI()) // no handle error, if error test will be terminated
	drv := entsql.OpenDB("mysql", db)
	entClient := ent.NewClient(ent.Driver(drv))

	// 3. Drop tables in DB
	_, _ = db.ExecContext(ctx, "DROP TABLE users")
	_, _ = db.ExecContext(ctx, "DROP TABLE file_award_point")
	_, _ = db.ExecContext(ctx, "DROP TABLE member_transaction")

	// 4. Migration DB Schema
	if err := entClient.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed Creating Schema Resources: %v", err)
	}

	// 5. Mocking data to database
	mockFileAwardPoint(ctx, entClient)
	mockXXX(ctx, entClient) // will mock data for other models

	fmt.Println() // new line in console => don't care about it

	return db
}

func mockFileAwardPoint(ctx context.Context, dbClient *ent.Client) {
	logger.Infof("Mock File Award Point ...")
	_, err := fileawardpoint.SaveAll(ctx, dbClient, fileAwardPoints, false)
	if err != nil {
		logger.Errorf("Mock File Award Point ... Failed: %v", err)
	}
	logger.Infof("Mock File Award Point ... Finished")
}

func mockXXX(_ context.Context, _ *ent.Client) {
	logger.Infof("Mock XXX ...")
	// doSth ...
	logger.Infof("Mock XXX ... Finished")
}
