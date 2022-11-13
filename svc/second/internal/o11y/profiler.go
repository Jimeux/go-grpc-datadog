package o11y

import (
	"log"

	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func InitProfiler(svcName, env, ver string) func() {
	err := profiler.Start(
		profiler.WithService(svcName),
		profiler.WithEnv(env),
		profiler.WithVersion(ver),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			// The profiles below are disabled by
			// default to keep overhead low, but
			// can be enabled as needed.
			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	return profiler.Stop
}
