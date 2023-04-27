package main

import (
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
					n := network.NewNetwork(cCtx.String("name"))
					s := server.NewServer(cCtx.String("port"), n)
					list, err := player.FromFolder(cCtx.String("path"))
					if err != nil {
						return err
					}
					c := n.AddChannel(list)
					list2, _ := player.FromFolder(cCtx.String("path"))
					n.AddChannel(list2)
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
				Action: func(cCtx *cli.Context) error {
					port := cCtx.String("port")
					host := cCtx.String("host")
					c, err := client.Connect(host, port)
					if err != nil {
						fmt.Printf("Couldn't connect to homelab-cable server - is one running at %s?\n", host+":"+port)
						return err
					}
					channel, err := c.CurrentChannel()
					if err != nil {
						return err
					}

					fmt.Printf("%+v\n", channel)

					return nil
				},
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
