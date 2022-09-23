package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tmvrus/golang-testing/service"
	"github.com/tmvrus/golang-testing/storage"
)

const (
	supportAPIVersion    = "2.0"
	registerPayoutMethod = "registerPayout"
)

type Handler struct {
	auth    service.Authorizator
	storage storage.Payout
}

func NewHandler(a service.Authorizator, s storage.Payout) *Handler {
	return &Handler{
		auth:    a,
		storage: s,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}
	if c := r.Header.Get("Content-Type"); c != "application/json" {
		http.Error(w, "Expected content type application/json, got "+c, http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body read error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req request
	if err := json.Unmarshal(body, &req); err != nil {
		writeResponse(w, newErrResponse(0, apiError{
			Code:    codeInvalidJSON,
			Message: "invalid json",
			Data:    fmt.Sprintf("got invalid json: %q", string(body)),
		}))
		return
	}
	if err := req.Validate(); err != nil {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeInvalidMethodParams,
			Message: "invalid request",
			Data:    fmt.Sprintf("got invalid request: %s", err.Error()),
		}))
		return
	}

	if req.Method != registerPayoutMethod {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeMethodNotFound,
			Message: "method not found",
			Data:    fmt.Sprintf("requested rpc-method %s is not found", req.Method),
		}))
		return
	}

	ctx := r.Context()
	token := r.Header.Get("API-Token")
	if token == "" {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeNotNotAuthorized,
			Message: "empty api token",
		}))
		return
	}
	ok, err := h.auth.Authorized(ctx, token)
	if err != nil {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeFailedAuthorization,
			Message: "authorization process has ben failed",
			Data:    token,
		}))
		return
	}
	if !ok {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeNotNotAuthorized,
			Message: "user not authorized",
			Data:    token,
		}))
		return
	}
	if err := h.storage.Register(ctx, req.Params.UserID, req.RequestID, req.Params.Payout); err != nil {
		writeResponse(w, newErrResponse(req.RequestID, apiError{
			Code:    codeUnknownError,
			Message: "unknown error",
			Data:    fmt.Sprintf("%v", req),
		}))
		return
	}

	writeResponse(w, newResponse(req.RequestID, []byte(`ok`)))
}
