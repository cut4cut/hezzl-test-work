package app

import (
	"fmt"
	"net"
	"time"

	"github.com/cut4cut/hezzl-test-work/config"
	"github.com/cut4cut/hezzl-test-work/internal/controller/rpc"
	"github.com/cut4cut/hezzl-test-work/internal/usecase"
	"github.com/cut4cut/hezzl-test-work/internal/usecase/repo"
	kf "github.com/cut4cut/hezzl-test-work/pkg/kafka"
	"github.com/cut4cut/hezzl-test-work/pkg/logger"
	"github.com/cut4cut/hezzl-test-work/pkg/postgres"
	rd "github.com/cut4cut/hezzl-test-work/pkg/redis"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) {
	// Wait other containers befor init app
	time.Sleep(15 * time.Second)

	l := logger.New(cfg.Log.Level)

	// Postgres
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Redis
	rdc, err := rd.New(cfg.Redis.Port, cfg.Redis.Expiration)
	if err != nil {
		l.Fatal(fmt.Errorf("app - run - redis.New: %w", err))
	}

	// Kafka
	prd, err := kf.New()
	if err != nil {
		l.Fatal(fmt.Errorf("app - run - kafka.New: %w", err))
	}
	defer func() {
		if err := prd.Close(); err != nil {
			l.Fatal(fmt.Errorf("app - run - prd.Close(): %w", err))
		}
	}()

	// Repo
	r := repo.New(pg)
	userUseCase := usecase.New(r, rdc, prd)

	// RPC - Net
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.RPC.Port))
	if err != nil {
		l.Fatal(fmt.Errorf("app - run - net.Listen: %w", err))
	}

	// RPC - Serv
	s := grpc.NewServer()
	rpc.Register(s, userUseCase, l)
	l.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		l.Fatal(fmt.Errorf("app - run - s.Serve: %w", err))
	}
}
