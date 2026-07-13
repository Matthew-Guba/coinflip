package game

import (
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

// GenerateServerSeed — криптографически стойкий сид сервера
func GenerateServerSeed() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

// HashSeed — хэш, который отдаём игроку ДО игры
func HashSeed(seed string) string {
    h := sha256.Sum256([]byte(seed))
    return hex.EncodeToString(h[:])
}

// RollFloat — детерминированное число [0,1) на основе (serverSeed, clientSeed, nonce)
func RollFloat(serverSeed, clientSeed string, nonce int64) float64 {
    msg := fmt.Sprintf("%s:%d", clientSeed, nonce)
    mac := hmac.New(sha256.New, []byte(serverSeed))
    mac.Write([]byte(msg))
    sum := mac.Sum(nil)

    // берём первые 8 байт → uint64 → нормализуем
    var v uint64
    for i := 0; i < 8; i++ {
        v = (v << 8) | uint64(sum[i])
    }
    return float64(v) / float64(^uint64(0))
}