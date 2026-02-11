package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

// RootFlags are global flags available to all commands.
type RootFlags struct {
	JSON    bool   `help:"Output JSON to stdout" short:"j"`
	Plain   bool   `help:"Output plain TSV (for scripting)"`
	NoColor bool   `help:"Disable colors" env:"NO_COLOR"`
	Device  string `help:"Device name or IP" short:"d" env:"LAMETRIC_DEVICE"`
	Verbose bool   `help:"Enable verbose logging" short:"v"`
}

// CLI is the top-level Kong CLI struct.
type CLI struct {
	RootFlags `embed:""`

	Notify    NotifyCmd    `cmd:"" name:"notify" help:"Send a notification"`
	Dismiss   DismissCmd   `cmd:"" name:"dismiss" help:"Dismiss a notification"`
	Version   VersionCmd   `cmd:"" name:"version" help:"Show version information"`
	Setup     SetupCmd     `cmd:"" help:"Setup wizard for a new device"`
	Discover  DiscoverCmd  `cmd:"" help:"Find LaMetric devices on the network"`
	Auth      AuthCmd      `cmd:"" help:"Manage API keys"`
	Device    DeviceCmd    `cmd:"" help:"Show device information"`
	Display   DisplayCmd   `cmd:"" name:"display" help:"Display settings"`
	Audio     AudioCmd     `cmd:"" name:"audio" help:"Audio settings"`
	Bluetooth BluetoothCmd `cmd:"" name:"bluetooth" help:"Bluetooth settings"`
	WiFi      WiFiCmd      `cmd:"" name:"wifi" help:"WiFi status"`
	Sounds    SoundsCmd    `cmd:"" name:"sounds" help:"List available sounds"`
	Icons     IconsCmd     `cmd:"" name:"icons" help:"List available icon aliases"`
	App       AppCmd       `cmd:"" name:"app" help:"App control"`
	Stream    StreamCmd    `cmd:"" name:"stream" help:"Streaming controls"`
}

type exitPanic struct{ code int }

// Execute runs the CLI with the given arguments.
func Execute(args []string) (err error) {
	parser, err := newParser()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if ep, ok := r.(exitPanic); ok {
				if ep.code == 0 {
					err = nil
					return
				}

				err = &ExitError{Code: ep.code, Err: errors.New("exited")}
				return
			}

			panic(r)
		}
	}()

	if len(args) == 0 {
		args = []string{"--help"}
	}

	kctx, err := parser.Parse(args)
	if err != nil {
		parsedErr := wrapParseError(err)
		_, _ = fmt.Fprintln(os.Stderr, parsedErr)
		return parsedErr
	}

	err = kctx.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}

func wrapParseError(err error) error {
	if err == nil {
		return nil
	}

	var parseErr *kong.ParseError
	if errors.As(err, &parseErr) {
		return &ExitError{Code: CodeUsage, Err: parseErr}
	}

	return err
}

func newParser() (*kong.Kong, error) {
	vars := kong.Vars{
		"version": VersionString(),
	}

	cli := &CLI{}
	parser, err := kong.New(
		cli,
		kong.Name("lametric"),
		kong.Description("LaMetric CLI - control your LaMetric TIME/SKY from the command line"),
		kong.Vars(vars),
		kong.Writers(os.Stdout, os.Stderr),
		kong.Exit(func(code int) { panic(exitPanic{code: code}) }),
		kong.Bind(&cli.RootFlags),
	)
	if err != nil {
		return nil, err
	}

	return parser, nil
}
