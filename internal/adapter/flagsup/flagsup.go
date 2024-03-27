package flagsup

import (
	"context"

	flagsup "go.tekoapis.com/flagsup/sdk/grpc-go"
)

type ClientAdapter interface {
	IsEnabled(ctx context.Context, flagKey string, fallback bool) bool
	IsEpicMa175Enabled(ctx context.Context) bool
}

type Client struct {
	host   string
	client flagsup.CachedClient
}

func New(host string) ClientAdapter {
	flagSupClient, err := flagsup.NewCachedClient(flagsup.WithServer(host))
	if err != nil {
		flagSupClient = nil
	}
	client := &Client{
		host:   host,
		client: flagSupClient,
	}

	return client
}

func (c *Client) IsEnabled(ctx context.Context, flagKey string, fallback bool) bool {
	if c.client == nil {
		return fallback
	}
	return c.client.GetFlagStatus(ctx, &flagsup.GetFlagStatusRequest{FlagKey: flagKey}, flagsup.Enabled(fallback)).Enabled
}

func (c *Client) IsEpicMa175Enabled(ctx context.Context) bool {
	return c.IsEnabled(ctx, FlagEpicMa175, FlagEpicMa175Default)
}
