package config

import (
	"os"
	"testing"
)

func TestReload(t *testing.T) {
	env := map[string]string{
		"BITFLYER_API_KEY":    "bitflyer_api_key_value",
		"BITFLYER_API_SECRET": "bitflyer_api_secret_value",
	}
	teardown := replaceGetenv(env)
	defer teardown()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BitFlyerAPIKey != env["BITFLYER_API_KEY"] {
		t.Fatalf("unexpected cfg.BitFlyerAPIKey: %s != %s", cfg.BitFlyerAPIKey, env["BITFLYER_API_KEY"])
	}

	env["BITFLYER_API_KEY"] = "bitflyer_api_key_value_replaced"
	if err := cfg.Reload(); err != nil {
		t.Fatalf("unexpected error on reload: %v", err)
	}

	if cfg.BitFlyerAPIKey != env["BITFLYER_API_KEY"] {
		t.Fatalf("unexpected cfg.BitFlyerAPIKey: %s != %s", cfg.BitFlyerAPIKey, env["BITFLYER_API_KEY"])
	}
}

func replaceGetenv(mapping map[string]string) func() {
	getenv = func(key string) string {
		if v, ok := mapping[key]; ok {
			return v
		}

		return ""
	}

	return func() {
		getenv = os.Getenv
	}
}
