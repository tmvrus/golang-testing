package handler_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/tmvrus/golang-testing/handler"
	"github.com/tmvrus/golang-testing/handler/mock"
)

//go:generate mockgen -destination=./mock/storage.go -package=mock -source=../storage/storage.go
//go:generate mockgen -destination=./mock/auth.go -package=mock -source=../service/auth.go

func Test_Handler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		in        []byte
		want      string
		setupMock func(auth *mock.MockAuthorizator, stor *mock.MockPayout)
	}{
		{
			name: "check validation",
			in:   []byte(`{"jsonrpc":"XXX","method":"registerPayout","id":1,"Params":{"user_id":"632d90378b2cb2b83351f130","payout":10}}`),
			want: `{"error":{"code":-32602, "data":"got invalid request: unsuported version XXX", "message":"invalid request"}, "id":1, "jsonrpc":"2.0","result":null}`,
		},
		{
			name: "check auth",
			in:   []byte(`{"jsonrpc":"2.0","method":"registerPayout","id":1,"Params":{"user_id":"632d90378b2cb2b83351f130","payout":10}}`),
			want: `{"error":{"code":-32002, "data":"TOKEN", "message":"user not authorized"}, "id":1, "jsonrpc":"2.0", "result":null}`,
			setupMock: func(auth *mock.MockAuthorizator, stor *mock.MockPayout) {
				m := map[string]string{
					"token":   "TOKEN",
					"user_id": "632d90378b2cb2b83351f130",
				}
				auth.
					EXPECT().
					Authorized(gomock.Any(), m).
					Return(false, nil)
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctr := gomock.NewController(t)
			t.Cleanup(ctr.Finish)
			authMock := mock.NewMockAuthorizator(ctr)

			if test.setupMock != nil {
				test.setupMock(authMock, nil)
			}
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(test.in))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("API-Token", "TOKEN")
			w := httptest.NewRecorder()

			smock := storMock(func() {})
			handler.NewHandler(authMock, smock).ServeHTTP(w, req)

			response, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)
			require.JSONEqf(t, test.want, string(response), string(response))

		})
	}
}

type storMock func()

func (storMock) Register(ctx context.Context, userID string, reqID int64, payout float64) error {
	return nil
}
func (storMock) Count(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}
