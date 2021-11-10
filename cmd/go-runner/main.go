package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"go.uber.org/zap/zapcore"

	_ "github.com/finecloud/apisix-oauth2-plugin/cmd/go-runner/plugins"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/log"
	"github.com/finecloud/apisix-oauth2-plugin/pkg/runner"
)

var (
	InfoOut io.Writer = os.Stdout
)

func newVersionCommand() *cobra.Command {
	var long bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version",
		Run: func(cmd *cobra.Command, _ []string) {
			if long {
				_, _ = fmt.Fprint(InfoOut, longVersion())
			} else {
				_, _ = fmt.Fprintf(InfoOut, "version %s\n", shortVersion())
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&long, "long", false, "show long mode version information")
	return cmd
}

type RunMode enumflag.Flag

const (
	Dev  RunMode = iota // Development
	Prod                // Product
	Prof                // Profile

	ProfileFilePath = "./logs/profile."
	LogFilePath     = "./logs/runner.log"
)

var RunModeIds = map[RunMode][]string{
	Prod: {"prod"},
	Dev:  {"dev"},
	Prof: {"prof"},
}

func openFileToWrite(name string) (*os.File, error) {
	dir := filepath.Dir(name)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func newRunCommand() *cobra.Command {
	var mode RunMode
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run",
		Run: func(cmd *cobra.Command, _ []string) {
			cfg := runner.RunnerConfig{}
			if mode == Prod {
				cfg.LogLevel = zapcore.WarnLevel
				f, err := openFileToWrite(LogFilePath)
				if err != nil {
					log.Fatalf("failed to open log: %s", err)
				}
				cfg.LogOutput = f
			} else if mode == Prof {
				cfg.LogLevel = zapcore.WarnLevel

				cpuProfileFile := ProfileFilePath + "cpu"
				f, err := os.Create(cpuProfileFile)
				if err != nil {
					log.Fatalf("could not create CPU profile: %s", err)
				}
				defer func(f *os.File) {
					_ = f.Close()
				}(f)
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatalf("could not start CPU profile: %s", err)
				}
				defer pprof.StopCPUProfile()

				defer func() {
					memProfileFile := ProfileFilePath + "mem"
					f, err := os.Create(memProfileFile)
					if err != nil {
						log.Fatalf("could not create memory profile: %s", err)
					}
					defer func(f *os.File) {
						_ = f.Close()
					}(f)

					runtime.GC()
					if err := pprof.WriteHeapProfile(f); err != nil {
						log.Fatalf("could not write memory profile: %s", err)
					}
				}()
			}
			runner.Run(cfg)
		},
	}

	cmd.PersistentFlags().VarP(
		enumflag.New(&mode, "mode", RunModeIds, enumflag.EnumCaseInsensitive),
		"mode", "m",
		"the runner's run mode; can be 'prod' or 'dev', default to 'dev'")

	return cmd
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apisix-go-plugin-runner [command]",
		Long:    "The Plugin runner to run Go plugins",
		Version: shortVersion(),
	}

	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newVersionCommand())
	return cmd
}

func main() {
	root := NewCommand()
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
