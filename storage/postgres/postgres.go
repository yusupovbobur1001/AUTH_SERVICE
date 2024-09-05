package postgres

import (
	"auth_service/configs"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	DB  *pgxpool.Pool
	log logger.ILogger
	cfg configs.Config
}

func NewStore(ctx context.Context, log logger.ILogger, cnf configs.Config) (*Store, error) {
	fmt.Println(cnf.PostgresHost, cnf.PostgresPort, cnf.PostgresDB, cnf.PostgresPassword)

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", 
		cnf.PostgresUser,
		cnf.PostgresPassword,
		cnf.PostgresHost,
		cnf.PostgresPort,
		cnf.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("Error parsing Postgres URL", logger.Error(err))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("Error creating pool with config", logger.Error(err))
		return nil, err
	}

	return &Store{
		DB:  pool,
		log: log,
		cfg: cnf,
	}, nil
}

func (s Store) Users() storage.IUserStorage {
	return NewUserRepository(s.DB, s.log)
}
func (s Store) Auths() storage.IAuthStorage {
	return NewAuthRepository(s.DB, s.log)
}

func (s Store) Close() {
	s.DB.Close()
}
