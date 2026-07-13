package handlers

import (
    "encoding/json"
    "net/http"
    "os"

    "coinflip/internal/game"
)

type rtpReq struct {
    Value float64 `json:"value"`
}

func adminAuth(r *http.Request) bool {
    return r.Header.Get("X-Admin-Token") == os.Getenv("ADMIN_TOKEN")
}

func SetRTPHandler(w http.ResponseWriter, r *http.Request) {
    if !adminAuth(r) {
        http.Error(w, "forbidden", 403)
        return
    }
    var req rtpReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad json", 400)
        return
    }
    if err := game.SetRTP(r.Context(), req.Value); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    w.Write([]byte(`{"ok":true}`))
}

func GetRTPHandler(w http.ResponseWriter, r *http.Request) {
    if !adminAuth(r) {
        http.Error(w, "forbidden", 403)
        return
    }
    json.NewEncoder(w).Encode(map[string]float64{"rtp": game.GetRTP()})
}