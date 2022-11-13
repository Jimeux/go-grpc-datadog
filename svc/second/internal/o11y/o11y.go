package o11y

func InitAll(name, env, ver, logPath string) func() {
	stopLog := InitLogger(name, env, ver, logPath)
	stopTrace := InitTracer(name, env, ver)
	stopProfile := InitProfiler(name, env, ver)

	return func() {
		stopLog()
		stopTrace()
		stopProfile()
	}
}
