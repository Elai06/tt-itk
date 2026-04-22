package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Client struct {
	Conn DB
}

func New(ctx context.Context, dbUrl string) (*Client, error) {
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Client{Conn: conn}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.Conn.Close(ctx)
}