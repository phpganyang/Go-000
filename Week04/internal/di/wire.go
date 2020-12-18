package di

import (
	"Week04/internal/dao"
	"Week04/internal/server/http"
	"Week04/internal/service"
	"github.com/google/wire"
)

func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.New, http.New, NewApp))
}
