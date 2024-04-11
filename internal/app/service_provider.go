package app

import (
	"context"
	"log"
	"os"
	"time"

	accessService "github.com/mistandok/auth/internal/service/access"

	"github.com/mistandok/auth/internal/api/access"
	endpointAccessRepository "github.com/mistandok/auth/internal/repository/endpoint_access"

	"github.com/gomodule/redigo/redis"
	"github.com/mistandok/auth/internal/api/auth"
	authService "github.com/mistandok/auth/internal/service/auth"
	jwtService "github.com/mistandok/auth/internal/service/jwt"
	"github.com/mistandok/auth/internal/utils/password"

	"github.com/mistandok/platform_common/pkg/closer"
	"github.com/mistandok/platform_common/pkg/db"
	"github.com/mistandok/platform_common/pkg/db/pg"
	"github.com/mistandok/platform_common/pkg/memory_db"
	"github.com/mistandok/platform_common/pkg/memory_db/rs"

	"github.com/mistandok/auth/internal/api/user"
	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/config/env"
	"github.com/mistandok/auth/internal/repository"
	userRepository "github.com/mistandok/auth/internal/repository/user"
	whiteList "github.com/mistandok/auth/internal/repository/white_list"
	"github.com/mistandok/auth/internal/service"
	userService "github.com/mistandok/auth/internal/service/user"

	"github.com/rs/zerolog"
)

type serviceProvider struct {
	pgConfig       *config.PgConfig
	grpcConfig     *config.GRPCConfig
	httpConfig     *config.HTTPConfig
	swaggerConfig  *config.SwaggerConfig
	passwordConfig *config.PasswordConfig
	jwtConfig      *config.JWTConfig
	redisConfig    *config.RedisConfig
	logger         *zerolog.Logger
	passManager    *password.Manager

	dbClient      db.Client
	txManager     db.TxManager
	redisDbClient memory_db.Client

	userRepo           repository.UserRepository
	whiteListRepo      repository.WhiteListRepository
	endpointAccessRepo repository.EndpointAccessRepository

	chatService   service.UserService
	authService   service.AuthService
	jwtService    service.JWTService
	accessService service.AccessService

	userImpl   *user.Implementation
	authImpl   *auth.Implementation
	accessImpl *access.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PgConfig ..
func (s *serviceProvider) PgConfig() *config.PgConfig {
	if s.pgConfig == nil {
		cfgSearcher := env.NewPgCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig ..
func (s *serviceProvider) GRPCConfig() *config.GRPCConfig {
	if s.grpcConfig == nil {
		cfgSearcher := env.NewGRPCCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить pg config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// HTTPConfig ..
func (s *serviceProvider) HTTPConfig() *config.HTTPConfig {
	if s.httpConfig == nil {
		cfgSearcher := env.NewHTTPCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

// SwaggerConfig ..
func (s *serviceProvider) SwaggerConfig() *config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfgSearcher := env.NewSwaggerConfigSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// PasswordConfig ..
func (s *serviceProvider) PasswordConfig() *config.PasswordConfig {
	if s.passwordConfig == nil {
		cfgSearcher := env.NewPasswordConfigSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить password config: %s", err.Error())
		}

		s.passwordConfig = cfg
	}

	return s.passwordConfig
}

// JWTConfig ..
func (s *serviceProvider) JWTConfig() *config.JWTConfig {
	if s.jwtConfig == nil {
		cfgSearcher := env.NewJWTConfigSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить jwt config: %s", err.Error())
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

// RedisConfig ..
func (s *serviceProvider) RedisConfig() *config.RedisConfig {
	if s.redisConfig == nil {
		cfgSearcher := env.NewRedisCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить white list redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

// Logger ..
func (s *serviceProvider) Logger() *zerolog.Logger {
	if s.logger == nil {
		cfgSearcher := env.NewLogCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить pg config: %s", err.Error())
		}

		s.logger = setupZeroLog(cfg)
	}

	return s.logger
}

// PassManager ..
func (s *serviceProvider) PassManager() *password.Manager {
	if s.passManager == nil {
		s.passManager = password.NewManager(s.PasswordConfig())
	}

	return s.passManager
}

// DBClient ..
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PgConfig().DSN(), s.Logger())
		if err != nil {
			log.Fatalf("ошибка при создании клиента DB: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("нет связи с БД: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager ..
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = pg.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) RedisDBClient(_ context.Context) memory_db.Client {
	if s.redisDbClient == nil {
		redisConfig := s.RedisConfig()
		redisPool := &redis.Pool{
			MaxIdle:     redisConfig.MaxIdle,
			IdleTimeout: redisConfig.IdleTimeout,
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, "tcp", redisConfig.Address())
			},
			TestOnBorrowContext: func(ctx context.Context, conn redis.Conn, lastUsed time.Time) error {
				if time.Since(lastUsed) < time.Minute {
					return nil
				}
				_, err := conn.Do("PING")
				return err
			},
		}
		s.redisDbClient = rs.New(redisPool)

		closer.Add(s.redisDbClient.Close)
	}

	return s.redisDbClient
}

// UserRepository ..
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepository.NewRepo(s.Logger(), s.DBClient(ctx))
	}

	return s.userRepo
}

// EndpointAccessRepository ..
func (s *serviceProvider) EndpointAccessRepository(ctx context.Context) repository.EndpointAccessRepository {
	if s.endpointAccessRepo == nil {
		s.endpointAccessRepo = endpointAccessRepository.NewRepo(s.Logger(), s.DBClient(ctx))
	}

	return s.endpointAccessRepo
}

func (s *serviceProvider) WhiteListRepository(ctx context.Context) repository.WhiteListRepository {
	if s.whiteListRepo == nil {
		s.whiteListRepo = whiteList.NewWhiteListRepo(s.RedisDBClient(ctx))
	}

	return s.whiteListRepo
}

// UserService ..
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.chatService == nil {
		s.chatService = userService.NewService(
			s.Logger(),
			s.UserRepository(ctx),
			s.PassManager(),
		)
	}

	return s.chatService
}

func (s *serviceProvider) JWTService(ctx context.Context) service.JWTService {
	if s.jwtService == nil {
		s.jwtService = jwtService.NewService(s.Logger(), s.JWTConfig(), s.WhiteListRepository(ctx))
	}

	return s.jwtService
}

// UserImpl ..
func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

// AuthService ..
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.Logger(),
			s.UserRepository(ctx),
			s.JWTService(ctx),
			s.PassManager(),
		)
	}

	return s.authService
}

// AccessService ..
func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(s.Logger(), s.EndpointAccessRepository(ctx), s.JWTService(ctx))
	}

	return s.accessService
}

// AuthImpl ..
func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func setupZeroLog(logConfig *config.LogConfig) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: logConfig.TimeFormat}
	logger := zerolog.New(output).With().Timestamp().Logger()
	logger = logger.Level(logConfig.LogLevel)
	zerolog.TimeFieldFormat = logConfig.TimeFormat

	return &logger
}

// AccessImpl ..
func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}
