package discord

import (
	"context"
	"log"
	"sync"

	"github.com/bonedaddy/unibot/bclient"
	"github.com/bonedaddy/unibot/db"
	"github.com/bwmarrin/discordgo"
)

var (
	rateLimitMsg = "You are being rate limited. Users are allowed 1 blockchain query per command every 60 seconds"
)

// Client wraps bclient and discordgo to provide a discord bot for indexed finance
type Client struct {
	s  *discordgo.Session
	bc *bclient.Client
	db *db.Database

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

// NewClient provides a wrapper around discordgo
func NewClient(ctx context.Context, cfg *Config, bc *bclient.Client, db *db.Database) (*Client, error) {
	wg := &sync.WaitGroup{}

	for _, watcher := range cfg.Watchers {
		wg.Add(1)
		go func(discToken, token0, token1 string) {
			defer wg.Done()
			dg, err := discordgo.New("Bot " + discToken)
			if err != nil {
				log.Println("failed to start watcher: ", err)
				return
			}

			if err := dg.Open(); err != nil {
				log.Println("failed to start watcher: ", err)
				return
			}
		}(watcher.DiscordToken, watcher.Token0Address, watcher.Token1Address)
	}

	client := &Client{bc: bc, wg: wg, db: db}

	log.Println("bot is now running")
	return client, nil
}

// Close terminates the discordgo session
func (c *Client) Close() error {
	c.wg.Wait()
	return c.s.Close()
}
