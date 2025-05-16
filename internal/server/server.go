package server

import (
	"github.com/google/wire"
	"go-cs/internal/server/auth/server3auth"
	"go-cs/internal/server/file"
	temp_file_task "go-cs/internal/server/file/task"
	"go-cs/internal/server/job"
	"go-cs/internal/server/kafka"
	"go-cs/internal/server/river"
	"go-cs/internal/server/websock"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	websock.NewServer,
	file.NewServer,
	temp_file_task.New,
	NewGRPCServer,
	NewHTTPServer,
	kafka.NewKafkaBroker,
	job.NewExample,
	job.NewCron,
	NewEtcdClient,
	server3auth.InitAuthServer,
	river.New,
)
