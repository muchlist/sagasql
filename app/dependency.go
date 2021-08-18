package app

import (
	"github.com/muchlist/sagasql/dao"
	"github.com/muchlist/sagasql/handler"
	"github.com/muchlist/sagasql/service"
	"github.com/muchlist/sagasql/utils/mcrypt"
	"github.com/muchlist/sagasql/utils/mjwt"
)

var (
	// Utils
	cryptoUtils = mcrypt.NewCrypto()
	jwt         = mjwt.NewJwt()

	// User Domain
	userDao     = dao.NewUserDao()
	userService = service.NewUserService(userDao, cryptoUtils, jwt)
	userHandler = handler.NewUserHandler(userService)
)
