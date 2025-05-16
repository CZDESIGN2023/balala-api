package dwh

import (
	"go-cs/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

func NewApp(
	conf *conf.Dwh,
	log log.Logger,
) (*DwhApp, func(), error) {

	app, cleanup, err := wireApp(conf, log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	return app, cleanup, err
}
