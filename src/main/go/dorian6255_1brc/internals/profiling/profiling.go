package profiling

import (
	"os"
	"runtime/pprof"
	"runtime/trace"
)

func profiling() {
	f, err := os.Create("profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()
	// Start tracing
	traceFile, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer traceFile.Close()

	if err := trace.Start(traceFile); err != nil {
		panic(err)
	}
	defer trace.Stop()

}
