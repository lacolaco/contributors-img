package config

import (
	"fmt"
	"strings"
	"testing"
)

// StartServer prints the loaded Config at startup, which ends up in Cloud
// Logging, so formatting it must not reveal the token.
func TestConfig_FormatDoesNotLeakToken(t *testing.T) {
	cfg := &Config{Port: "3333", GitHubAuthToken: "super-secret-token"}
	for _, format := range []string{"%v", "%+v", "%s"} {
		if out := fmt.Sprintf(format, cfg); strings.Contains(out, "super-secret-token") {
			t.Fatalf("%s leaks the token: %s", format, out)
		}
	}
}
