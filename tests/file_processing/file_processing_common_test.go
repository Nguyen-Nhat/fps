package fileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileprocessing"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"github.com/xo/dburl"
)

func initFileProcessingServerForTesting() *fileprocessing.Server {
	dbConf := config.Load("../..").Database.MySQL
	db, _ := dburl.Open(dbConf.DatabaseURI()) // no handle error, if error test will be terminated
	fileProcessingServer := fileprocessing.InitFileProcessingServer(db)
	return fileProcessingServer
}
