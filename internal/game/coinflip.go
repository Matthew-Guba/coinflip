package game

import (
    "context"
    "errors"

    "coinflip/internal/db"
)

var (
    ErrInvalidRTP    = errors.New("invalid rtp")
    ErrInvalidSide   = errors.New("invalid side")
    ErrInvalidBet    = errors.New("invalid bet")
    ErrLowBalance    = errors.New("insufficient balance")
)

type PlayResult struct {
    RoundID       int64
    Result        string  // "heads" | "tails"
    Won           bool
    WinAmount     float64
    Coefficient   float64
    NewBalance    float64
    ServerSeed    string
    ServerHash    string
    ClientSeed    string
    Nonce         int64
    RTPUsed       float64
}

// Play — основная логика раунда coinflip
// Стандартный coinflip: 50/50, коэффициент = 2.0
// Мы регулируем шанс победы под RTP: winChance = RTP/100 / 2 * ... 
// Правильнее: если coeff = 2.0, то winChance = RTP/200. Пример: RTP 97% → winChance = 48.5%
func Play(ctx context.Context, userID int64, betAmount float64, side, clientSeed string, nonce int64) (*PlayResult, error) {
    if side != "heads" && side != "tails" {
        return nil, ErrInvalidSide
    }
    if betAmount <= 0 {
        return nil, ErrInvalidBet
    }

    rtp := GetRTP()
    coefficient := 2.0
    winChance := rtp / 100.0 / coefficient  // напр. 0.485

    // Транзакция
    tx, err := db.Pool.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

    var balance float64
    err = tx.QueryRow(ctx,
        `SELECT balance FROM users WHERE id=$1 FOR UPDATE`, userID).Scan(&balance)
    if err != nil {
        return nil, err
    }
    if balance < betAmount {
        return nil, ErrLowBalance
    }

    // списываем ставку
    _, err = tx.Exec(ctx,
        `UPDATE users SET balance = balance - $1 WHERE id=$2`, betAmount, userID)
    if err != nil {
        return nil, err
    }

    // генерируем результат
    serverSeed, err := GenerateServerSeed()
    if err != nil {
        return nil, err
    }
    hash := HashSeed(serverSeed)
    roll := RollFloat(serverSeed, clientSeed, nonce)

    // выигрыш определяется вероятностью
    won := roll < winChance
    var result string
    if won {
        result = side
    } else {
        if side == "heads" {
            result = "tails"
        } else {
            result = "heads"
        }
    }

    winAmount := 0.0
    if won {
        winAmount = betAmount * coefficient
        _, err = tx.Exec(ctx,
            `UPDATE users SET balance = balance + $1 WHERE id=$2`, winAmount, userID)
        if err != nil {
            return nil, err
        }
    }

    var roundID int64
    err = tx.QueryRow(ctx, `
        INSERT INTO rounds
        (user_id, bet_amount, side, result, win_amount, coefficient,
         server_seed, client_seed, nonce, server_seed_hash, rtp_used)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        RETURNING id
    `, userID, betAmount, side, result, winAmount, coefficient,
        serverSeed, clientSeed, nonce, hash, rtp).Scan(&roundID)
    if err != nil {
        return nil, err
    }

    var newBalance float64
    err = tx.QueryRow(ctx,
        `SELECT balance FROM users WHERE id=$1`, userID).Scan(&newBalance)
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, err
    }

    return &PlayResult{
        RoundID:     roundID,
        Result:      result,
        Won:         won,
        WinAmount:   winAmount,
        Coefficient: coefficient,
        NewBalance:  newBalance,
        ServerSeed:  serverSeed,
        ServerHash:  hash,
        ClientSeed:  clientSeed,
        Nonce:       nonce,
        RTPUsed:     rtp,
    }, nil
}