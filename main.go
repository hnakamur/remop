package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const defaultConfigFilename = "remop.json"

var programLevel = new(slog.LevelVar)

func main() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: programLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			_ = groups
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					// format timestamp with millisecond precision
					a.Value = slog.StringValue(t.Format("2006-01-02T15:04:05.000Z07:00"))
				}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(h))

	if err := run(os.Args); err != nil {
		slog.Error("exit with error", "args", os.Args, "err", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 1 {
		switch args[1] {
		case "tunnel":
			return runTunnelSubcommand(args)
		case "read":
			return runReadSubcommand(args)
		case "version":
			return runVersion(args)
		}
	}

	showUsageAndExit()
	return nil
}

func runTunnelSubcommand(args []string) error {
	fsName := strings.Join(args[:2], " ")
	fs := flag.NewFlagSet(fsName, flag.ExitOnError)
	configFilename := fs.String("c", defaultConfigFilename, "config file path")
	debug := fs.Bool("debug", false, "enable debug log")
	if err := fs.Parse(args[2:]); err != nil {
		return err
	}

	if *debug {
		programLevel.Set(slog.LevelDebug)
	}

	if fs.NArg() != 1 {
		return errors.New("targetHost arg needed")
	}

	cfg, err := readConfigFile(*configFilename)
	if err != nil {
		return err
	}
	return runTunnelServer(cfg, fs.Args()[0])
}

func runReadSubcommand(args []string) error {
	fsName := strings.Join(args[:2], " ")
	fs := flag.NewFlagSet(fsName, flag.ExitOnError)
	socketPath := fs.String("socket", "${REMOP_SOCKET}", "unix socket path")
	debug := fs.Bool("debug", false, "enable debug log")
	if err := fs.Parse(args[2:]); err != nil {
		return err
	}

	if *debug {
		programLevel.Set(slog.LevelDebug)
	}

	if fs.NArg() != 1 {
		return errors.New("itemURL arg needed")
	}

	return doRunReadSubcommand(os.ExpandEnv(*socketPath), fs.Args()[0])
}

func runVersion([]string) error {
	fmt.Println(version())
	return nil
}

func version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return "(devel)"
}

func showUsageAndExit() {
	fmt.Fprintf(os.Stderr, `Usage: %s subcommand options

subcommand:
    tunnel                 start tunnel server
    read                   read the field specified with op://... URL.
    version                print version and exit
`, os.Args[0])
	os.Exit(2)
}
