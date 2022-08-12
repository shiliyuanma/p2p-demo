package hole

import "github.com/urfave/cli/v2"

var (
	Cmd = &cli.Command{
		Name:  "hole",
		Usage: "tradition hole punch",
		Subcommands: []*cli.Command{
			relayCmd,
			peerCmd,
		},
	}

	relayCmd = &cli.Command{
		Name: "relay",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 6002,
			},
		},
		Action: func(ctx *cli.Context) error {
			return startRelay(ctx.Int("port"))
		},
	}

	peerCmd = &cli.Command{
		Name: "peer",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "relay",
				Usage:    "relay address(127.0.0.1:6002)",
				Value:    "127.0.0.1:6002",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "role",
				Value: "p2pv",
				Usage: "p2pv/p2pp",
			},
			&cli.StringFlag{
				Name:  "pwd",
				Value: "123456",
			},
		},
		Action: func(ctx *cli.Context) error {
			relay := ctx.String("relay")
			pwd := ctx.String("pwd")
			role := WORK_P2P_PROVIDER
			if ctx.String("role") != role {
				role = WORK_P2P_VISITOR
			}
			return startPeer(relay, role, pwd)
		},
	}
)
