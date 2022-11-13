package o11y

import (
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func InitTracer(svcName, env, ver string) func() {
	rules := []tracer.SamplingRule{
		tracer.RateRule(0.1),
	}
	tracer.Start(
		tracer.WithSamplingRules(rules),
		tracer.WithRuntimeMetrics(),
		tracer.WithService(svcName),
		tracer.WithEnv(env),
		tracer.WithUniversalVersion(ver),
	)
	return tracer.Stop
}

func ServerMiddleware(svcName string) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	ui := grpctrace.UnaryServerInterceptor(
		grpctrace.WithServiceName(svcName),
		grpctrace.WithStreamMessages(false),
		grpctrace.WithIgnoredMethods("/grpc.health.v1.Health/Check"),
	)
	si := grpctrace.StreamServerInterceptor(
		grpctrace.WithServiceName(svcName),
		grpctrace.WithStreamMessages(false),
		grpctrace.WithIgnoredMethods(
			"/grpc.health.v1.Health/Watch",
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		),
	)
	return ui, si
}

func ClientMiddleware(svcName string) (grpc.UnaryClientInterceptor, grpc.StreamClientInterceptor) {
	ui := grpctrace.UnaryClientInterceptor(
		grpctrace.WithServiceName(svcName),
		grpctrace.WithStreamMessages(false),
	)
	si := grpctrace.StreamClientInterceptor(
		grpctrace.WithServiceName(svcName),
		grpctrace.WithStreamCalls(false),
	)
	return ui, si
}
