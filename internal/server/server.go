package server

import (
	"context"
	"errors"
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"github.com/ozonmp/lic-license-api/internal/pkg/grpc/interceptor/grpc_logs"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/ozonmp/lic-license-api/internal/api"
	"github.com/ozonmp/lic-license-api/internal/config"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
)

type licenseService interface {
	Get(tx context.Context, subdomainID uint64) (*license.License, error)
	Add(ctx context.Context, service *license.License) (uint64, error)
	List(ctx context.Context, offset uint64, limit uint64) ([]*license.License, error)
	Remove(ctx context.Context, serviceID uint64) (bool, error)
}

// GrpcServer is gRPC server
type GrpcServer struct {
	service   licenseService
	db        *sqlx.DB
	batchSize uint
}

// NewGrpcServer returns gRPC server with supporting of batch listing
func NewGrpcServer(db *sqlx.DB, batchSize uint) *GrpcServer {
	return &GrpcServer{
		db:        db,
		batchSize: batchSize,
	}
}

// Start method runs server
func (s *GrpcServer) Start(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)
	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)

	gatewayServer := createGatewayServer(ctx, grpcAddr, gatewayAddr)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("Gateway server is running on %s", gatewayAddr))
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "Failed running gateway server")
			cancel()
		}
	}()

	metricsServer := createMetricsServer(cfg)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("Metrics server is running on %s", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "Failed running metrics server")
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := createStatusServer(ctx, cfg, isReady)

	go func() {
		statusAdrr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)
		logger.InfoKV(ctx, fmt.Sprintf("Status server is running on %s", statusAdrr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "Failed running status server")
		}
	}()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			logger.DebugKV(ctx, "Failed close listen", "err", err)
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			grpc_logs.MetadataChangingLogsLevelUnaryServerInterceptor(),
			grpc_zap.PayloadUnaryServerInterceptor(
				logger.FromContext(ctx).Desugar(),
				grpc_logs.GetIsEnableDescribeRequestAndResponseDecider(),
			),
		)),
	)

	pb.RegisterLicLicenseApiServiceServer(grpcServer, api.NewLicenseAPI(s.service))
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("GRPC Server is listening on: %s", grpcAddr))
		if err := grpcServer.Serve(l); err != nil {
			logger.FatalKV(ctx, "Failed running gRPC server", "err", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		logger.InfoKV(ctx, "The service is ready to accept requests")
	}()

	if cfg.Project.Debug {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.InfoKV(ctx, fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.InfoKV(ctx, fmt.Sprintf("ctx.Done: %v", done))
	}

	isReady.Store(false)

	if err := gatewayServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "gatewayServer.Shutdown", "err", err)
	} else {
		logger.ErrorKV(ctx, "gatewayServer shut down correctly", "err", err)
	}

	if err := statusServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "statusServer.Shutdown", "err", err)
	} else {
		logger.InfoKV(ctx, "statusServer shut down correctly")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "metricsServer.Shutdown", "err", err)
	} else {
		logger.InfoKV(ctx, "metricsServer shut down correctly")
	}

	grpcServer.GracefulStop()
	logger.InfoKV(ctx, "grpcServer shut down correctly")

	return nil
}
