package config

import "testing"

func TestAppVersion(t *testing.T) {
	saved := Version
	t.Cleanup(func() { Version = saved })

	Version = "v1.2.3"
	if got := appVersion(); got != "v1.2.3" {
		t.Errorf("appVersion() = %q, want the version set at build time", got)
	}

	// Test binaries have no module version: the fallback is used
	Version = ""
	if got := appVersion(); got != "development" {
		t.Errorf("appVersion() = %q, want %q", got, "development")
	}
}
