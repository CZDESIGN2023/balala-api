package test

import (
	"go-cs/internal/conf"
	"go-cs/internal/dwh/data"

	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

type MockDwhCtx struct {
	dwhData          *data.DwhData
	JobVariablesRepo *data.JobVariablesRepo
}

func NewMockDwhCtx() *MockDwhCtx {

	dwhConf := MockDwhConfig()
	dwhData := data.NewDwhData(dwhConf)
	jobVariablesRepo := data.NewJobVariablesRepo(dwhData)

	return &MockDwhCtx{
		dwhData:          data.NewDwhData(dwhConf),
		JobVariablesRepo: jobVariablesRepo,
	}
}

var mockDwhConfig *conf.Dwh

func MockDwhConfig() *conf.Dwh {

	if mockDwhConfig != nil {
		return mockDwhConfig
	}

	c := kconfig.New(
		kconfig.WithSource(
			file.NewSource("/Users/cyt/project/ed-project-manage/server_go/configs/config.mock.test.yaml"),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var b conf.Bootstrap

	if err := c.Scan(&b); err != nil {
		panic(err)
	}

	mockDwhConfig = b.Dwh
	return mockDwhConfig
}
