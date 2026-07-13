package handlers

import (
	"encoding/json"
	"net/http"

	"coinflip/internal/db"
	"coinflip/internal/game"
)

type playReq struct {
	UserID     int64   `json:"user_id"`
	Bet        float64 `json:"bet"`
	Side       string  `json:"side"`
	ClientSeed string  `json:"client_seed"`
	Nonce      int64   `json:"nonce"`
}

// POST /api/play
func PlayHandler(w http.ResponseWriter, r *http.Request) {
	var req playReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	res, err := game.Play(r.Context(), req.UserID, req.Bet, req.Side, req.ClientSeed, req.Nonce)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GET /api/balance?user_id=1
func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("user_id")
	if uid == "" {
		http.Error(w, "user_id required", 400)
		return
	}

	var balance float64
	err := db.Pool.QueryRow(r.Context(),
		`SELECT balance FROM users WHERE id=$1`, uid).Scan(&balance)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

// GET /api/history?user_id=1&limit=20
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("user_id")
	limit := r.URL.Query().Get("limit")
	if uid == "" {
		http.Error(w, "user_id required", 400)
		return
	}
	if limit == "" {
		limit = "20"
	}

	rows, err := db.Pool.Query(r.Context(), `
		SELECT id, bet_amount, side, result, win_amount, coefficient, created_at
		FROM rounds WHERE user_id=$1
		ORDER BY created_at DESC LIMIT $2
	`, uid, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type row struct {
		ID          int64   `json:"id"`
		BetAmount   float64 `json:"bet_amount"`
		Side        string  `json:"side"`
		Result      string  `json:"result"`
		WinAmount   float64 `json:"win_amount"`
		Coefficient float64 `json:"coefficient"`
		CreatedAt   string  `json:"created_at"`
	}

	out := []row{}
	for rows.Next() {
		var it row
		if err := rows.Scan(&it.ID, &it.BetAmount, &it.Side, &it.Result,
			&it.WinAmount, &it.Coefficient, &it.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		out = append(out, it)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}