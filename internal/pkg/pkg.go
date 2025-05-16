package pkg

import (
	"go-cs/internal/pkg/biz_id"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// http_api.NewRegistryInterface,
	// server3.NewServer3Api,
	// openapi.NewOpenAPIService,
	// ipc.NewIpcClient,
	biz_id.NewBusinessIdService,
)
