package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"shiliyuanma/p2p-demo/hole"
	"shiliyuanma/p2p-demo/nat"
	"shiliyuanma/p2p-demo/x"

	"github.com/urfave/cli/v2"
)

var (
	app *cli.App
)

func init() {
	app = &cli.App{
		Name:    "p2p-demo",
		Version: "v1.0.0",
		After: func(ctx *cli.Context) error {
			if x.Hold {
				ch := make(chan os.Signal, 1)
				signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
				<-ch
			}
			x.Close()
			return nil
		},
		ExitErrHandler: func(ctx *cli.Context, err error) {
			fmt.Println("exit with error:", err)
			x.Close()
		},
		Commands: []*cli.Command{
			nat.Cmd,
			hole.Cmd,
		},
	}
}

func main() {
	//app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
