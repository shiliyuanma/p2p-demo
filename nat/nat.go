package nat

import (
	"fmt"

	"github.com/ccding/go-stun/stun"
	"github.com/urfave/cli/v2"
)

var (
	Cmd = &cli.Command{
		Name:  "nat",
		Usage: "nat type check",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "stun-addr",
				Aliases: []string{"addr"},
				Value:   "stun.stunprotocol.org:3478",
				Usage:   "stun server address (eg:stun.stunprotocol.org:3478)",
			},
		},
		Action: func(ctx *cli.Context) error {
			c := stun.NewClient()
			c.SetServerAddr(ctx.String("stun-addr"))
			nat, host, err := c.Discover()
			if err != nil || host == nil {
				return err
			}
			fmt.Printf("nat type: %s \npublic address: %s\n", nat.String(), host.String())
			return nil
		},
	}
)
