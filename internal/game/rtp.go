package game

import (
    "context"
    "sync"

    "coinflip/internal/db"
)

var (
    mu         sync.RWMutex
    currentRTP = 97.0
)

func LoadRTP(ctx context.Context) error {
    var v float64
    err := db.Pool.QueryRow(ctx, `SELECT value FROM rtp_config WHERE id=1`).Scan(&v)
    if err != nil {
        return err
    }
    mu.Lock()
    currentRTP = v
    mu.Unlock()
    return nil
}

func GetRTP() float64 {
    mu.RLock()
    defer mu.RUnlock()
    return currentRTP
}

func SetRTP(ctx context.Context, v float64) error {
    if v < 50 || v > 100 {
        return ErrInvalidRTP
    }
    _, err := db.Pool.Exec(ctx,
        `UPDATE rtp_config SET value=$1, updated_at=NOW() WHERE id=1`, v)
    if err != nil {
        return err
    }
    mu.Lock()
    currentRTP = v
    mu.Unlock()
    return nil
}