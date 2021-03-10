//+build wireinject

package worker

import (
	"github.com/google/wire"
	"mcm-api/internal/core"
	"mcm-api/pkg/article"
	"mcm-api/pkg/converter"
	"mcm-api/pkg/media"
)

func InitializeWorker() *worker {
	panic(wire.Build(
		core.InfraSet,
		converter.Set,
		article.Set,
		media.Set,
		newWorker))
}
