package command

import (
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/govenue/goman"
	"github.com/govenue/notepad"
)

var (
	benchmarkTimes int
	cpuProfileFile string
	memProfileFile string
)

var commandBenchmark = &goman.Command{
	Use:   "benchmark",
	Short: "Benchmark Gean by building a site a number of times.",
	Long: `Gean can build a site many times over and analyze the running process
creating a benchmark.`,
}

func init() {
	initHugoBuilderFlags(commandBenchmark)
	initBenchmarkBuildingFlags(commandBenchmark)

	commandBenchmark.Flags().StringVar(&cpuProfileFile, "cpuprofile", "", "path/filename for the CPU profile file")
	commandBenchmark.Flags().StringVar(&memProfileFile, "memprofile", "", "path/filename for the memory profile file")
	commandBenchmark.Flags().IntVarP(&benchmarkTimes, "count", "n", 13, "number of times to build the site")

	commandBenchmark.RunE = benchmark
}

func benchmark(cmd *goman.Command, args []string) error {
	cfg, err := InitializeConfig(commandBenchmark)
	if err != nil {
		return err
	}

	c, err := newCommandeer(cfg)
	if err != nil {
		return err
	}

	var memProf *os.File
	if memProfileFile != "" {
		memProf, err = os.Create(memProfileFile)
		if err != nil {
			return err
		}
	}

	var cpuProf *os.File
	if cpuProfileFile != "" {
		cpuProf, err = os.Create(cpuProfileFile)
		if err != nil {
			return err
		}
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memAllocated := memStats.TotalAlloc
	mallocs := memStats.Mallocs
	if cpuProf != nil {
		pprof.StartCPUProfile(cpuProf)
	}

	t := time.Now()
	for i := 0; i < benchmarkTimes; i++ {
		if err = c.resetAndBuildSites(false); err != nil {
			return err
		}
	}
	totalTime := time.Since(t)

	if memProf != nil {
		pprof.WriteHeapProfile(memProf)
		memProf.Close()
	}
	if cpuProf != nil {
		pprof.StopCPUProfile()
		cpuProf.Close()
	}

	runtime.ReadMemStats(&memStats)
	totalMemAllocated := memStats.TotalAlloc - memAllocated
	totalMallocs := memStats.Mallocs - mallocs

	notepad.FEEDBACK.Println()
	notepad.FEEDBACK.Printf("Average time per operation: %vms\n", int(1000*totalTime.Seconds()/float64(benchmarkTimes)))
	notepad.FEEDBACK.Printf("Average memory allocated per operation: %vkB\n", totalMemAllocated/uint64(benchmarkTimes)/1024)
	notepad.FEEDBACK.Printf("Average allocations per operation: %v\n", totalMallocs/uint64(benchmarkTimes))

	return nil
}
