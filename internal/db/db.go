package db

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(ctx context.Context, url string) error {
    p, err := pgxpool.New(ctx, url)
    if err != nil {
        return err
    }
    Pool = p
    return Pool.Ping(ctx)
}