package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bonedaddy/unibot/bclient"
	"github.com/bonedaddy/unibot/db"
	"github.com/bonedaddy/unibot/discord"
	"github.com/bonedaddy/unibot/watcher"
	"github.com/urfave/cli/v2"
)

var (
	bc      *bclient.Client
	Version string
)

func main() {
	app := cli.NewApp()
	app.Name = "unibot"
	app.Usage = "unibot is a discord price watching bot for uniswap"
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "infura.api_key",
			Usage:   "api key for use with infura",
			EnvVars: []string{"INFURA_API_KEY"},
		},
		&cli.StringFlag{
			Name:  "eth.rpc",
			Usage: "specifies the ethereum RPC endpoint",
			Value: bclient.InfuraHTTPURL,
		},
		&cli.StringFlag{
			Name:  "eth.address",
			Usage: "address to lookup in queries",
			Value: "0x5a361A1dfd52538A158e352d21B5b622360a7C13",
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"cfg"},
			Usage:   "path to discord bot configuration file",
			Value:   "config.yml",
		},
		&cli.BoolFlag{
			Name:  "startup.sleep",
			Usage: "whether or not to sleep on startup, useful for giving containers time to initialize",
			Value: false,
		},
		&cli.DurationFlag{
			Name:  "startup.sleep_time",
			Usage: "time.Duration type specifying sleep duration",
			Value: time.Second * 5,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("startup.sleep") {
			time.Sleep(c.Duration("startup.sleep_time"))
		}
		if c.String("config") != "" {
			return nil
		}
		var err error
		if c.String("infura.api_key") != "" {
			var websockets bool
			if strings.Contains(c.String("infura.api_key"), "wss") {
				websockets = true
			}
			bc, err = bclient.NewInfuraClient(c.String("infura.api_key"), websockets)
		} else {
			bc, err = bclient.NewClient(c.String("eth.rpc"))
		}
		return err
	}
	app.Commands = cli.Commands{
		&cli.Command{
			Name:  "discord",
			Usage: "discord bot management",
			Subcommands: cli.Commands{
				&cli.Command{
					Name:    "database",
					Aliases: []string{"db"},
					Usage:   "database management commands",
					Subcommands: cli.Commands{
						&cli.Command{
							Name:        "chain-updater",
							Usage:       "starts the database chain state updater",
							Description: "chain-updater is responsible for persisting all information needed into a database such as price updates",
							Action: func(c *cli.Context) error {
								ctx, cancel := context.WithCancel(c.Context)
								defer cancel()
								cfg, err := discord.LoadConfig(c.String("config"))
								if err != nil {
									return err
								}
								if cfg.InfuraAPIKey != "" {
									bc, err = bclient.NewInfuraClient(cfg.InfuraAPIKey, cfg.InfuraWSEnabled)
								} else {
									bc, err = bclient.NewClient(cfg.ETHRPCEndpoint)
								}
								if err != nil {
									return err
								}
								defer bc.Close()
								database, err := db.New(&db.Opts{
									Type:           cfg.Database.Type,
									Host:           cfg.Database.Host,
									Port:           cfg.Database.Port,
									User:           cfg.Database.User,
									Password:       cfg.Database.Pass,
									DBName:         cfg.Database.DBName,
									SSLModeDisable: cfg.Database.SSLModeDisable,
								})
								if err != nil {
									return err
								}
								defer database.Close()
								if err := database.AutoMigrate(); err != nil {
									return err
								}
								items := watcher.ConfigToWatchItmes(cfg)
								watchService := watcher.New(ctx, database, bc, time.Second*5, items)
								watchService.Start()
								sc := make(chan os.Signal, 1)
								signal.Notify(sc, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt, os.Kill)
								<-sc
								watchService.Stop()
								return nil
							},
						},
					},
				},
				&cli.Command{
					Name:  "gen-config",
					Usage: "generate ndx bot config file",
					Action: func(c *cli.Context) error {
						return discord.NewConfig(c.String("config"))
					},
				},
				&cli.Command{
					Name:  "ndx-bot",
					Usage: "starts NDXBot",
					Action: func(c *cli.Context) error {
						ctx, cancel := context.WithCancel(c.Context)
						defer cancel()
						cfg, err := discord.LoadConfig(c.String("config"))
						if err != nil {
							return err
						}
						if cfg.InfuraAPIKey != "" {
							bc, err = bclient.NewInfuraClient(cfg.InfuraAPIKey, cfg.InfuraWSEnabled)
						} else {
							bc, err = bclient.NewClient(cfg.ETHRPCEndpoint)
						}
						if err != nil {
							return err
						}
						defer bc.Close()
						database, err := db.New(&db.Opts{
							Type:           cfg.Database.Type,
							Host:           cfg.Database.Host,
							Port:           cfg.Database.Port,
							User:           cfg.Database.User,
							Password:       cfg.Database.Pass,
							DBName:         cfg.Database.DBName,
							SSLModeDisable: cfg.Database.SSLModeDisable,
						})
						if err != nil {
							return err
						}
						defer database.Close()
						if err := database.AutoMigrate(); err != nil {

							return err
						}
						wg := &sync.WaitGroup{}
						if c.Bool("update.database") {
							wg.Add(1)
							// launch database price updater loop
							go func() {
								defer wg.Done()
								// dbPriceUpdateLoop(ctx, bc, database)
							}()
						}

						client, err := discord.NewClient(ctx, cfg, bc, database)
						if err != nil {
							return err
						}
						sc := make(chan os.Signal, 1)
						signal.Notify(sc, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt, os.Kill)
						<-sc
						cancel()
						wg.Wait()
						return client.Close()
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "discord.token",
							Usage:   "the discord api token",
							EnvVars: []string{"DISCORD_TOKEN"},
						},
						&cli.BoolFlag{
							Name:  "update.database",
							Usage: "if true launch the db price update routine. if false make sure chain-updater command is running",
							Value: true,
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
