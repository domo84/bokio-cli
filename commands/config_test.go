package commands

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// The key table drives validation, completion and help text. If help stops listing a
// key, users lose the only place the keys are documented.
func TestConfigHelpListsEveryKey(t *testing.T) {
	long := newConfigCmd().Long

	for _, key := range configKeys {
		if !strings.Contains(long, key.Name) {
			t.Errorf("config help does not mention key %q", key.Name)
		}
		if !strings.Contains(long, key.Env) {
			t.Errorf("config help does not mention env var %q", key.Env)
		}
	}

	if !strings.Contains(long, configFilePath()) {
		t.Errorf("config help does not mention the config file path %q", configFilePath())
	}
}

func TestConfigSubcommandHelpListsKeys(t *testing.T) {
	for _, sub := range newConfigCmd().Commands() {
		if sub.Name() == "list" {
			continue // list documents sources, not the keys themselves
		}
		for _, key := range configKeys {
			if !strings.Contains(sub.Long, key.Name) {
				t.Errorf("%q help does not mention key %q", sub.Name(), key.Name)
			}
		}
	}
}

func TestIsValidKey(t *testing.T) {
	for _, name := range configKeyNames() {
		if !isValidKey(name) {
			t.Errorf("isValidKey(%q) = false, want true", name)
		}
	}

	// auth_mode was removed: nothing read it, so accepting it misled users.
	for _, name := range []string{"auth_mode", "", "companyid", "COMPANY_ID"} {
		if isValidKey(name) {
			t.Errorf("isValidKey(%q) = true, want false", name)
		}
	}
}

func TestUnknownKeyErrorListsValidKeys(t *testing.T) {
	err := unknownKeyError("auth_mode")
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, name := range configKeyNames() {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not list %q", err, name)
		}
	}
}

func TestConfigDisplayValueMasksSecrets(t *testing.T) {
	t.Cleanup(viper.Reset)

	secret := configKey{Name: "client_secret", Secret: true}
	plain := configKey{Name: "company_id"}

	viper.Set(secret.Name, "super-secret-value")
	viper.Set(plain.Name, "3fa85f64-5717-4562-b3fc-2c963f66afa6")

	if got := configDisplayValue(secret); got != "(set)" {
		t.Errorf("secret value = %q, want %q", got, "(set)")
	}
	if got := configDisplayValue(plain); got != "3fa85f64-5717-4562-b3fc-2c963f66afa6" {
		t.Errorf("plain value = %q", got)
	}

	viper.Set(secret.Name, "")
	if got := configDisplayValue(secret); got != "(not set)" {
		t.Errorf("unset secret = %q, want %q", got, "(not set)")
	}
}

func TestConfigValueSource(t *testing.T) {
	t.Cleanup(viper.Reset)

	key := configKey{Name: "company_id", Env: "BOKIO_COMPANY_ID"}

	if got := configValueSource(key); got != "default" {
		t.Errorf("source with nothing set = %q, want %q", got, "default")
	}

	t.Setenv(key.Env, "from-env")
	if got := configValueSource(key); got != "env" {
		t.Errorf("source with env set = %q, want %q", got, "env")
	}
}
