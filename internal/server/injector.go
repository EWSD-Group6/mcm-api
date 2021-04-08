//+build wireinject

package server

import (
	"github.com/google/wire"
	"mcm-api/internal/core"
	"mcm-api/pkg/article"
	"mcm-api/pkg/authz"
	"mcm-api/pkg/comment"
	"mcm-api/pkg/contributesession"
	"mcm-api/pkg/contribution"
	"mcm-api/pkg/faculty"
	"mcm-api/pkg/media"
	"mcm-api/pkg/startup"
	"mcm-api/pkg/statistic"
	"mcm-api/pkg/systemdata"
	"mcm-api/pkg/user"
)

func InitializeServer() *Server {
	panic(wire.Build(
		core.InfraSet,
		user.Set,
		authz.Set,
		startup.Set,
		faculty.Set,
		media.Set,
		contributesession.Set,
		contribution.Set,
		article.Set,
		comment.Set,
		systemdata.Set,
		statistic.Set,
		core.HandlerSet,
		newServer,
	))
}
