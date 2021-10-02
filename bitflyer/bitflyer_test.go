package bitflyer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetBalance(t *testing.T) {
	const responseBody = `[
  {
    "currency_code": "JPY",
    "amount": 1024078,
    "available": 508000
  },
  {
    "currency_code": "BTC",
    "amount": 10.24,
    "available": 4.12
  },
  {
    "currency_code": "ETH",
    "amount": 20.48,
    "available": 16.38
  }
]`
	const (
		apikey    = "longlongaccesskey"
		apisecret = "longlonglongapisecret"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("request method is unexpected: %s != %s", r.Method, http.MethodGet)
		}

		if got := r.Header.Get("ACCESS-KEY"); got != apikey {
			t.Fatalf("ACCESS-KEY header is unexpected: %s != %s", got, apikey)
		}

		// TODO: test sign

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer srv.Close()

	c := New(apikey, "apisecret")
	c.httpEndpoint = srv.URL

	got, err := c.GetBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	want := []Balance{
		{
			CurrencyCode: "JPY",
			Amount:       1024078,
			Available:    508000,
		},
		{
			CurrencyCode: "BTC",
			Amount:       10.24,
			Available:    4.12,
		},
		{
			CurrencyCode: "ETH",
			Amount:       20.48,
			Available:    16.38,
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatal(diff)
	}
}
