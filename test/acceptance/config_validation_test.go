//go:build acceptance_a

// Config validation acceptance tests.
//
// These exercise the real gc binary's config subcommands (show, explain)
// as a black box against various valid and invalid configurations. No
// supervisor is needed — these commands only read and validate config files.
package acceptance_test

import (
	"strings"
	"testing"

	helpers "github.com/gastownhall/gascity/test/acceptance/helpers"
)

// --- gc config show ---

// TestConfigShow_ValidConfig_DumpsTOML verifies that gc config show on
// a valid city outputs parseable TOML containing the workspace name.
func TestConfigShow_ValidConfig_DumpsTOML(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("config", "show")
	if err != nil {
		t.Fatalf("gc config show failed: %v\n%s", err, out)
	}

	// Should contain workspace name from init template.
	if !strings.Contains(out, "[workspace]") {
		t.Errorf("expected [workspace] section in TOML output, got:\n%s", out)
	}
}

// TestConfigShow_Validate_ValidConfig_Succeeds verifies that --validate
// exits 0 and prints "Config valid." for a well-formed config.
func TestConfigShow_Validate_ValidConfig_Succeeds(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("config", "show", "--validate")
	if err != nil {
		t.Fatalf("gc config show --validate failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Config valid.") {
		t.Errorf("expected 'Config valid.' message, got:\n%s", out)
	}
}

// TestConfigShow_Validate_DuplicateAgents_ReturnsError verifies that
// --validate catches duplicate agent names and returns a non-zero exit.
func TestConfigShow_Validate_DuplicateAgents_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	// Overwrite with config that has duplicate agents.
	c.WriteConfig(`[workspace]
name = "duptest"

[[agent]]
name = "worker"

[[agent]]
name = "worker"
`)

	out, err := c.GC("config", "show", "--validate")
	if err == nil {
		t.Fatal("expected error for duplicate agents, got success")
	}
	if !strings.Contains(out, "duplicate") {
		t.Errorf("expected 'duplicate' in error, got:\n%s", out)
	}
}

// TestConfigShow_InvalidTOML_ReturnsError verifies that gc config show
// on syntactically invalid TOML returns a parse error.
func TestConfigShow_InvalidTOML_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	// Overwrite with broken TOML.
	c.WriteConfig(`[workspace
name = "broken"
this is not valid toml !!!
`)

	out, err := c.GC("config", "show")
	if err == nil {
		t.Fatal("expected error for invalid TOML, got success")
	}
	_ = out // Error message format varies by TOML parser.
}

// TestConfigShow_Provenance_ShowsSources verifies that --provenance
// outputs source file information.
func TestConfigShow_Provenance_ShowsSources(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("config", "show", "--provenance")
	if err != nil {
		t.Fatalf("gc config show --provenance failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Sources") {
		t.Errorf("expected 'Sources' in provenance output, got:\n%s", out)
	}
	if !strings.Contains(out, "city.toml") {
		t.Errorf("expected 'city.toml' in provenance output, got:\n%s", out)
	}
}

// --- gc config explain ---

// TestConfigExplain_ShowsAgents verifies that gc config explain lists
// configured agents with their resolved fields.
func TestConfigExplain_ShowsAgents(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	// Write a config with a known agent.
	c.WriteConfig(`[workspace]
name = "explain-test"

[[agent]]
name = "coder"
start_command = "echo hello"
`)

	out, err := c.GC("config", "explain")
	if err != nil {
		t.Fatalf("gc config explain failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "coder") {
		t.Errorf("expected 'coder' agent in output, got:\n%s", out)
	}
}

// TestConfigExplain_AgentFilter_MatchesOne verifies that --agent filters
// to a single agent.
func TestConfigExplain_AgentFilter_MatchesOne(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	c.WriteConfig(`[workspace]
name = "filter-test"

[[agent]]
name = "alpha"
start_command = "echo a"

[[agent]]
name = "beta"
start_command = "echo b"
`)

	out, err := c.GC("config", "explain", "--agent", "alpha")
	if err != nil {
		t.Fatalf("gc config explain --agent alpha: %v\n%s", err, out)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected 'alpha' in output, got:\n%s", out)
	}
	if strings.Contains(out, "beta") {
		t.Errorf("should NOT contain 'beta' when filtered to alpha:\n%s", out)
	}
}

// TestConfigExplain_FilterNoMatch_ReturnsError verifies that filtering
// to a nonexistent agent returns an error mentioning the filters.
func TestConfigExplain_FilterNoMatch_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	c.WriteConfig(`[workspace]
name = "nomatch"

[[agent]]
name = "worker"
`)

	out, err := c.GC("config", "explain", "--agent", "nosuchagent")
	if err == nil {
		t.Fatal("expected error for unmatched filter, got success")
	}
	if !strings.Contains(out, "no agents match") {
		t.Errorf("expected 'no agents match' in error, got:\n%s", out)
	}
}

// --- gc config (bare command) ---

// TestConfig_NoSubcommand_ShowsHelp verifies that gc config with no
// subcommand shows help text (not an error).
func TestConfig_NoSubcommand_ShowsHelp(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("config")
	if err != nil {
		t.Fatalf("gc config should show help, not error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "show") || !strings.Contains(out, "explain") {
		t.Errorf("help text should mention subcommands 'show' and 'explain', got:\n%s", out)
	}
}
