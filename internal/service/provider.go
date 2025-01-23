package service

import (
	"github.com/google/wire"
	"github.com/lantonster/askme/internal/service/uploads"
)

var ProviderSetService = wire.NewSet(
	uploads.NewUploadsService,
)
