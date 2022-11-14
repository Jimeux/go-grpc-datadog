package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/Jimeux/go-grpc-datadog/proto/go/pb/second/v1"
	"github.com/Jimeux/go-grpc-datadog/svc/second/internal/config"
	"github.com/Jimeux/go-grpc-datadog/svc/second/internal/db"
	"github.com/Jimeux/go-grpc-datadog/svc/second/internal/o11y"
	"github.com/Jimeux/go-grpc-datadog/svc/second/internal/rpc"
)

func main() {
	_ = godotenv.Load()
	cf := config.Init()
	defer db.Init(cf.DatabaseUser, cf.DatabasePassword, cf.DatabaseHost, cf.DatabasePort, cf.DatabaseName)()
	defer o11y.InitAll(cf.ServiceName, cf.ServiceVersion, cf.Environment, cf.LogPath)()

	grpcServer := buildServer(cf.ServiceName)

	// open connection
	conn, err := net.Listen("tcp", ":"+cf.Port)
	if err != nil {
		o11y.Err(context.Background(), err, "failed to listen on "+cf.Port)
		panic(err)
	}
	defer func() { _ = conn.Close() }()

	// run server
	go func() {
		o11y.Info(context.Background(), "listening on "+cf.Port+"...")
		if err := grpcServer.Serve(conn); err != http.ErrServerClosed {
			o11y.Err(context.Background(), err, "unexpected Serve error on "+cf.Port)
		}
	}()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()
	<-ctx.Done()

	grpcServer.GracefulStop()
}

func buildServer(svcName string) *grpc.Server {
	// configure server
	recoverer := grpc_recovery.WithRecoveryHandler(func(p any) error {
		return status.Errorf(codes.Unknown, "panic encountered: %v", p)
	})
	unaryTrace, streamTrace := o11y.ServerMiddleware(svcName)
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			unaryTrace,
			grpc_recovery.UnaryServerInterceptor(recoverer),
		),
		grpc_middleware.WithStreamServerChain(
			streamTrace,
			grpc_recovery.StreamServerInterceptor(recoverer),
		),
	)
	reflection.Register(grpcServer)
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus(svcName, healthpb.HealthCheckResponse_SERVING)

	// register services
	second.RegisterSecondServiceServer(grpcServer, &rpc.SecondService{})

	return grpcServer
}
