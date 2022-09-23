package handler

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func Test_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   request
		err  bool
	}{
		{
			name: "ok",
			in: request{
				Version:   "2.0",
				Method:    "registerPayout",
				RequestID: 1,
				Params: params{
					UserID: "632d90378b2cb2b83351f130",
					Payout: 10,
				},
			},
		},
		{
			name: "invalid version",
			in: request{
				Version:   "invalid!!!",
				Method:    "registerPayout",
				RequestID: 1,
				Params: params{
					UserID: "632d90378b2cb2b83351f130",
					Payout: 10,
				},
			},
			err: true,
		},
	}
	// go test -coverprofile=coverage.out
	// go tool cover -func=coverage.out
	// go tool cover -html=coverage.out
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.in.Validate()
			check(t, err, test.err)
		})
	}
}

func check(t *testing.T, err error, flag bool) {
	t.Helper()

	if flag && err == nil {
		t.Error("expected error")
		return
	}
	if !flag && err != nil {
		t.Error("unexpected error")
		return
	}
}
