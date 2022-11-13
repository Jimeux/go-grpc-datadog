package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/Jimeux/go-grpc-datadog/proto/go/pb/first/v1"
	"github.com/Jimeux/go-grpc-datadog/proto/go/pb/second/v1"
	"github.com/Jimeux/go-grpc-datadog/svc/first/internal/config"
	"github.com/Jimeux/go-grpc-datadog/svc/first/internal/o11y"
	"github.com/Jimeux/go-grpc-datadog/svc/first/internal/rpc"
)

func main() {
	_ = godotenv.Load()
	cf := config.Init()
	defer o11y.InitAll(cf.ServiceName, cf.Environment, cf.ServiceVersion, cf.LogPath)()

	secondSvcClient, closeSecondSvc := buildSecondServiceClient(cf.SecondServiceName, cf.ServerServiceHost)
	defer closeSecondSvc()

	grpcServer := buildServer(cf.ServiceName, secondSvcClient)

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

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*20))
	defer cancel()
	grpcServer.GracefulStop()
}

func buildServer(svcName string, secondSvcClient second.SecondServiceClient) *grpc.Server {
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
	first.RegisterFirstServiceServer(grpcServer, rpc.NewClientService(secondSvcClient))

	return grpcServer
}

func buildSecondServiceClient(svcName, host string) (second.SecondServiceClient, func()) {
	cliUnaryTrace, cliStreamTrace := o11y.ClientMiddleware(svcName)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(cliUnaryTrace),
		grpc.WithStreamInterceptor(cliStreamTrace),
		grpc.WithDefaultCallOptions(grpc.MaxRetryRPCBufferSize(200)),
		// TODO 2022/11/13 @Jimeux
		// https://github.com/grpc/grpc/blob/master/doc/service_config.md
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
			  "name": [{"service": ""}],
			  "waitForReady": true,
			  "retryPolicy": {
				  "MaxAttempts": 5,
				  "InitialBackoff": ".01s",
				  "MaxBackoff": "1s",
				  "BackoffMultiplier": 1.0,
				  "RetryableStatusCodes": [ "UNAVAILABLE" ]
			  }
			}]}`),
	)
	if err != nil {
		o11y.Err(context.Background(), err, "failed to connect to SecondService on "+host)
		panic(err)
	}
	return second.NewSecondServiceClient(conn), func() { _ = conn.Close() }
}
