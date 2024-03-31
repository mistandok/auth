package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mistandok/auth/internal/api/auth"
	authService "github.com/mistandok/auth/internal/service/auth"
	jwtService "github.com/mistandok/auth/internal/service/jwt"
	"github.com/mistandok/auth/internal/utils/password"

	"github.com/mistandok/platform_common/pkg/closer"
	"github.com/mistandok/platform_common/pkg/db"
	"github.com/mistandok/platform_common/pkg/db/pg"

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
	pgConfig        *config.PgConfig
	grpcConfig      *config.GRPCConfig
	httpConfig      *config.HTTPConfig
	swaggerConfig   *config.SwaggerConfig
	passwordConfig  *config.PasswordConfig
	jwtConfig       *config.JWTConfig
	whiteListConfig *config.WhiteListRedisConfig
	logger          *zerolog.Logger
	passManager     *password.Manager

	dbClient      db.Client
	txManager     db.TxManager
	whiteListPool *redis.Pool

	userRepo      repository.UserRepository
	whiteListRepo repository.WhiteListRepository

	chatService service.UserService
	authService service.AuthService
	jwtService  service.JWTService

	userImpl *user.Implementation
	authImpl *auth.Implementation
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

// WhiteListConfig ..
func (s *serviceProvider) WhiteListConfig() *config.WhiteListRedisConfig {
	if s.whiteListConfig == nil {
		cfgSearcher := env.NewWhiteListRedisCfgSearcher()
		cfg, err := cfgSearcher.Get()
		if err != nil {
			log.Fatalf("не удалось получить white list redis config: %s", err.Error())
		}

		s.whiteListConfig = cfg
	}

	return s.whiteListConfig
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

func (s *serviceProvider) WhiteListPool(_ context.Context) *redis.Pool {
	if s.whiteListPool == nil {
		s.whiteListPool = &redis.Pool{
			MaxIdle:     5,
			IdleTimeout: 60 * time.Second,
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, "tcp", s.WhiteListConfig().Address())
			},
			TestOnBorrowContext: func(ctx context.Context, conn redis.Conn, lastUsed time.Time) error {
				if time.Since(lastUsed) < time.Minute {
					return nil
				}
				_, err := conn.Do("PING")
				return err
			},
		}
	}

	return s.whiteListPool
}

// UserRepository ..
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepository.NewRepo(s.Logger(), s.DBClient(ctx))
	}

	return s.userRepo
}

func (s *serviceProvider) WhiteListRepository(ctx context.Context) repository.WhiteListRepository {
	if s.whiteListRepo == nil {
		s.whiteListRepo = whiteList.NewWhiteListRepo(s.WhiteListPool(ctx))
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
