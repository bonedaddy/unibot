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
	ctx, cancel := context.WithCancel(ctx)

	dg, err := discordgo.New("Bot " + cfg.MainDiscordToken)
	if err != nil {
		cancel()
		return nil, err
	}

	if err := dg.Open(); err != nil {
		cancel()
		return nil, err
	}
	if err := dg.UpdateListeningStatus("!ndx help"); err != nil {
		log.Println("failed to udpate streaming status: ", err)
	}

	client := &Client{s: dg, bc: bc, ctx: ctx, cancel: cancel, wg: wg, db: db}

	log.Println("bot is now running")
	return client, nil
}

// Close terminates the discordgo session
func (c *Client) Close() error {
	c.cancel()
	c.wg.Wait()
	return c.s.Close()
}
