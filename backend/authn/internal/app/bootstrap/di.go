package bootstrap

import (
	"context"
	"log/slog"

	postgresadapter2 "github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres"
	"github.com/faber-numeris/beholder/authn/internal/core/services"
	"github.com/faber-numeris/beholder/authn/internal/infrastructure/config"
	infrapostgres "github.com/faber-numeris/beholder/authn/internal/infrastructure/postgres"
	"github.com/faber-numeris/beholder/authn/internal/platform/util"
	outboundport2 "github.com/faber-numeris/beholder/authn/internal/ports/outbound"
)

func ProvideHashingService() outboundport2.HashingService {
	return services.NewHashingService()
}

type repositoryComposition struct {
	outboundport2.UserRepository
	outboundport2.UserConfirmationRepository
}

func (r repositoryComposition) Ping() bool {
	return r.UserRepository.Ping()
}

func ProvideRepository(dbConfig config.IDatabaseConfig) outboundport2.Repository {

	pool := util.Must(infrapostgres.Connect(context.Background(), dbConfig))
	_ = pool

	db := infrapostgres.GetDB()
	if db == nil {
		slog.Warn("Database not connected, repositories will return errors until connection is established")
	}
	return repositoryComposition{
		postgresadapter2.NewUserRepository(db),
		postgresadapter2.NewUserConfirmationRepository(db),
	}
}
