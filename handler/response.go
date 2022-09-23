package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	codeInvalidJSON         int64 = -32600
	codeMethodNotFound            = -32601
	codeInvalidMethodParams       = -32602
	codeFailedAuthorization       = -32001
	codeNotNotAuthorized          = -32002
	codeUnknownError              = -32003
)

type response struct {
	Version string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *apiError       `json:"error,omitempty"`
}

type apiError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func newResponse(id int64, result json.RawMessage) response {
	return response{
		Version: supportAPIVersion,
		ID:      id,
		Result:  result,
	}
}

func newErrResponse(id int64, err apiError) response {
	return response{
		Version: supportAPIVersion,
		ID:      id,
		Error:   &err,
	}
}

func writeResponse(w http.ResponseWriter, res response) {
	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(response)))
	if _, err := w.Write(response); err != nil {
		log.Printf("faile to write response: %s", err.Error())
	}
}
