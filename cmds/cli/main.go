package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/clabland/go-homelab-cable/client"
	"github.com/clabland/go-homelab-cable/network"
	"github.com/clabland/go-homelab-cable/player"
	"github.com/clabland/go-homelab-cable/server"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start a homelab cable server",
				Action: func(cCtx *cli.Context) error {
					n := network.NewNetwork(cCtx.String("network_name"), cCtx.String("network_owner"))
					s := server.NewServer(cCtx.String("port"), n)
					list, err := player.FromFolder(cCtx.String("path"))
					if err != nil {
						return err
					}
					c := n.AddChannel(list)
					err = n.SetChannelLive(c.ID)
					if err != nil {
						return err
					}
					s.Serve()
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Value: "3004",
						Usage: "port to run on",
					},
					&cli.StringFlag{
						Name:  "network_name",
						Value: "Homelab Cable",
						Usage: "the name of your homelab cable network",
					},
					&cli.StringFlag{
						Name:  "network_owner",
						Value: "clabretro",
						Usage: "the owner of your homelab cable network",
					},
					&cli.StringFlag{
						Name:     "path",
						Value:    "",
						Usage:    "path to media folder",
						Required: true,
					},
				},
			},
			{
				Name:  "client",
				Usage: "start a homelab cable client to interact with an already-running server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Value: "3004",
						Usage: "server port to connect to",
					},
					&cli.StringFlag{
						Name:  "host",
						Value: "http://localhost",
						Usage: "host the server is running on",
					},
					&cli.BoolFlag{
						Name:  "json",
						Value: false,
						Usage: "output command results in JSON",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "live",
						Usage: "view information about the live channel",
						Action: func(cCtx *cli.Context) error {
							c, err := connect(cCtx)
							if err != nil {
								return err
							}
							channel, err := c.CurrentChannel()
							if err != nil {
								return err
							}

							if c.JSONOut {
								chanBytes, err := json.Marshal(channel)
								if err != nil {
									return err
								}
								fmt.Printf("%s", chanBytes)
								return nil
							}

							fmt.Printf("%s", channel)

							return nil
						},
					},
					{
						Name:  "play_next",
						Usage: "play the next piece of media for the live channel",
						Action: func(cCtx *cli.Context) error {
							c, err := connect(cCtx)
							if err != nil {
								return err
							}
							channel, err := c.LiveNext()
							if err != nil {
								return err
							}

							if c.JSONOut {
								chanBytes, err := json.Marshal(channel)
								if err != nil {
									return err
								}
								fmt.Printf("%s", chanBytes)
								return nil
							}

							fmt.Printf("%s", channel)

							return nil
						},
					},
				},
			},
			{
				Name:  "path_test",
				Usage: "list the media files a given --path would play",
				Action: func(cCtx *cli.Context) error {
					path := cCtx.String("path")
					list, err := player.FromFolder(path)
					if err != nil {
						return err
					}
					fmt.Println(list.All())
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Value:    "",
						Usage:    "path to media folder",
						Required: true,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func connect(ctx *cli.Context) (*client.Client, error) {
	port := ctx.String("port")
	host := ctx.String("host")
	jsonOut := ctx.Bool("json")
	c, err := client.Connect(host, port)
	c.JSONOut = jsonOut
	if err != nil {
		fmt.Printf("Couldn't connect to homelab-cable server - is one running at %s?\n", host+":"+port)
		return nil, err
	}
	return c, nil
}
