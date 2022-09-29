package user

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/user"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"github.com/xo/dburl"
)

func initUserServerForTesting() *user.UserServer {
	dbConf := config.Load().Database.MySQL
	db, _ := dburl.Open(dbConf.DatabaseURI()) // no handle error, if error test will be terminated
	userServer := user.InitUserServer(db)
	return userServer
}
