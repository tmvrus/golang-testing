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

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			err := test.in.Validate()
			if test.err && err == nil {
				t.Error("expected error")
				return
			}
			if !test.err && err != nil {
				t.Error("unexpected error")
				return
			}
		})
	}
}
