package handler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/tmvrus/golang-testing/handler"
)

//go:generate mockgen -destination=./mock/storage.go -package=mock -source=../storage/storage.go

func Test_Handler(t *testing.T) {
	t.Parallel()

	//	t.Skip("skip for now")

	tests := []struct {
		name string
		in   []byte
		want string
	}{
		{
			name: "check validation",
			in:   []byte(`{"jsonrpc":"XXX","method":"registerPayout","id":1,"Params":{"user_id":"632d90378b2cb2b83351f130","payout":10}}`),
			want: `{"error":{"code":-32602, "data":"got invalid request: unsuported version XXX", "message":"invalid request"}, "id":1, "jsonrpc":"2.0","result":null}`,
		},
		//{
		//	name: "check auth",
		//	in:   []byte(`{"jsonrpc":"2.0","method":"registerPayout","id":1,"Params":{"user_id":"632d90378b2cb2b83351f130","payout":10}}`),
		//	want: `{"error":{"code":-32002, "data":"TOKEN", "message":"user not authorized"}, "id":1, "jsonrpc":"2.0", "result":null}`,
		//},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(test.in))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("API-Token", "TOKEN")
			w := httptest.NewRecorder()

			handler.NewHandler(nil, nil).ServeHTTP(w, req)

			response, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)
			require.JSONEq(t, test.want, string(response))

		})
	}
}
