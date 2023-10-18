package client

import (
	"context"
	"net/http"
	"strings"

	graphql "github.com/hasura/go-graphql-client"
)

// GQLClient is an interface for a graphql client
type GQLClient interface {
	Query(tx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error
}

// Client is a client for the metadata api
type Client struct {
	gqlCli     GQLClient
	httpClient *http.Client
}

// Option is a function that modifies a client
type Option func(*Client)

// New creates a new metadata api client
func New(url string, opts ...Option) *Client {
	c := &Client{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.gqlCli = graphql.NewClient(url, c.httpClient)

	return c
}

// WithHTTPClient functional option to set the http client
func WithHTTPClient(cli *http.Client) Option {
	return func(c *Client) {
		c.httpClient = cli
	}
}

// StatusUpdate mutates the requested nodeID with json status
func (c *Client) StatusUpdate(ctx context.Context, input *StatusUpdateInput) (*StatusUpdate, error) {
	vars := map[string]interface{}{
		"input": *input,
	}

	r := new(StatusUpdate)
	if err := c.gqlCli.Mutate(ctx, r, vars); err != nil {
		return nil, translateGQLErr(err)
	}

	return r, nil
}

func translateGQLErr(err error) error {
	switch {
	case strings.Contains(err.Error(), "invalid or expired jwt"):
		return ErrUnauthorized
	case strings.Contains(err.Error(), "subject doesn't have access"):
		return ErrPermissionDenied
	}

	return err
}
