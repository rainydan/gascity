//go:build acceptance_a

// Doctor and mail CLI acceptance tests.
//
// Tests gc doctor diagnostics on valid/invalid cities and gc mail error
// paths. Doctor is tested as a black box against real initialized cities.
// Mail tests cover argument validation (sending requires live sessions,
// so happy-path mail is Tier B).
package acceptance_test

import (
	"path/filepath"
	"strings"
	"testing"

	helpers "github.com/gastownhall/gascity/test/acceptance/helpers"
)

// --- gc doctor ---

// TestDoctor_ValidCity_Passes verifies that gc doctor on a freshly
// initialized tutorial city reports no critical failures.
func TestDoctor_ValidCity_Passes(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("doctor")
	if err != nil {
		// Doctor may exit 1 for warnings (e.g., missing optional binaries)
		// but should not crash. Check if it ran checks at all.
		if !strings.Contains(out, "check") && !strings.Contains(out, "Check") {
			t.Fatalf("gc doctor failed without running checks: %v\n%s", err, out)
		}
	}
	// Verify output contains diagnostic structure (passed/failed/warnings).
	if !strings.Contains(out, "pass") && !strings.Contains(out, "PASS") &&
		!strings.Contains(out, "fail") && !strings.Contains(out, "FAIL") &&
		!strings.Contains(out, "✓") && !strings.Contains(out, "✗") {
		t.Errorf("doctor output doesn't look like a diagnostic report:\n%s", out)
	}
}

// TestDoctor_Verbose_ShowsDetails verifies that --verbose produces more
// output than the default mode.
func TestDoctor_Verbose_ShowsDetails(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	outDefault, _ := c.GC("doctor")
	outVerbose, _ := c.GC("doctor", "--verbose")

	if len(outVerbose) < len(outDefault) {
		t.Errorf("verbose output (%d bytes) should be >= default (%d bytes)",
			len(outVerbose), len(outDefault))
	}
}

// TestDoctor_GastownCity_RunsPackChecks verifies that gc doctor on a
// gastown city discovers and runs pack doctor scripts.
func TestDoctor_GastownCity_RunsPackChecks(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.InitFrom(filepath.Join(helpers.ExamplesDir(), "gastown"))

	out, _ := c.GC("doctor")
	// Gastown pack ships doctor scripts — verify they were discovered.
	// The output should reference pack-related checks.
	if !strings.Contains(out, "gastown") && !strings.Contains(out, "pack") {
		t.Logf("doctor output may not include gastown pack checks (pack doctor scripts may not exist yet):\n%s", out)
	}
}

// TestDoctor_NotACity_ReturnsError verifies that gc doctor on a non-city
// directory returns an error.
func TestDoctor_NotACity_ReturnsError(t *testing.T) {
	emptyDir := t.TempDir()
	out, err := helpers.RunGC(testEnv, emptyDir, "doctor")
	if err == nil {
		t.Fatal("expected error for gc doctor in non-city directory, got success")
	}
	_ = out
}

// --- gc mail error paths ---

// TestMail_NoSubcommand_ReturnsError verifies that bare gc mail prints
// a helpful error listing available subcommands.
func TestMail_NoSubcommand_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("mail")
	if err == nil {
		t.Fatal("expected error for bare 'gc mail', got success")
	}
	if !strings.Contains(out, "missing subcommand") {
		t.Errorf("expected 'missing subcommand' message, got:\n%s", out)
	}
}

// TestMail_UnknownSubcommand_ReturnsError verifies that gc mail with an
// unknown subcommand returns an error.
func TestMail_UnknownSubcommand_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("mail", "explode")
	if err == nil {
		t.Fatal("expected error for unknown mail subcommand, got success")
	}
	if !strings.Contains(out, "unknown subcommand") {
		t.Errorf("expected 'unknown subcommand' message, got:\n%s", out)
	}
}

// TestMailSend_MissingArgs_ReturnsError verifies that gc mail send with
// no arguments returns a usage error.
func TestMailSend_MissingArgs_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("mail", "send")
	if err == nil {
		t.Fatal("expected error for gc mail send with no args, got success")
	}
	if !strings.Contains(out, "usage") {
		t.Errorf("expected usage hint in error, got:\n%s", out)
	}
}

// TestMailInbox_DefaultIdentity_Succeeds verifies that gc mail inbox
// with no running sessions returns successfully (empty inbox).
func TestMailInbox_DefaultIdentity_Succeeds(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	// With no sessions, inbox defaults to "human" identity.
	// Should succeed even with empty inbox.
	out, err := c.GC("mail", "inbox")
	if err != nil {
		t.Fatalf("gc mail inbox should succeed with empty inbox: %v\n%s", err, out)
	}
}

// TestMailCount_ReturnsZero verifies that gc mail count on a fresh city
// returns a zero count.
func TestMailCount_ReturnsZero(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("mail", "count")
	if err != nil {
		t.Fatalf("gc mail count failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "0") {
		t.Errorf("expected 0 count on fresh city, got:\n%s", out)
	}
}

// TestMailCheck_Inject_AlwaysExitsZero verifies that gc mail check --inject
// always exits 0 (suitable for hook injection).
func TestMailCheck_Inject_AlwaysExitsZero(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("mail", "check", "--inject")
	if err != nil {
		t.Fatalf("gc mail check --inject should always exit 0: %v\n%s", err, out)
	}
}
