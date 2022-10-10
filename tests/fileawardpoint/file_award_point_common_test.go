package awardpoint

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/fileawardpoint"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"github.com/xo/dburl"
)

func initAwardPointServerForTesting() *fileawardpoint.Server {
	dbConf := config.Load("../..").Database.MySQL
	db, _ := dburl.Open(dbConf.DatabaseURI()) // no handle error, if error test will be terminated
	fileAwardPointServer := fileawardpoint.InitFileAwardPointServer(db)
	return fileAwardPointServer
}
