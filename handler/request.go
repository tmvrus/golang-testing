package handler

import (
	"encoding/hex"
	"fmt"
)

type request struct {
	Version   string `json:"jsonrpc"`
	Method    string `json:"method"`
	RequestID int64  `json:"id"`
	Params    params
}
type params struct {
	UserID string  `json:"user_id"`
	Payout float64 `json:"payout"`
}

func (r *request) Validate() error {
	const (
		userIDLen         = 24
		maxPayout float64 = 100
		minPayout float64 = 0.01
	)

	if r.Version != supportAPIVersion {
		return fmt.Errorf("unsuported version %s", r.Version)
	}
	if r.Method == "" {
		return fmt.Errorf("empty method passed")
	}
	if r.Params.Payout < minPayout || r.Params.Payout > maxPayout {
		return fmt.Errorf("invalid payout %f", r.Params.Payout)
	}
	if len(r.Params.UserID) != userIDLen {
		return fmt.Errorf("invalid userID length %q", r.Params.UserID)
	}
	if _, err := hex.DecodeString(r.Params.UserID); err != nil {
		return fmt.Errorf("userID expexted to be hex-format")
	}
	return nil
}
