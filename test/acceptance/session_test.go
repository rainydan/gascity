//go:build acceptance_a

// Session command acceptance tests.
//
// These exercise gc session subcommands as a black box. Session
// management is fundamental to the agent lifecycle. Most mutating
// operations need a running controller, so Tier A tests focus on
// list, prune, and error paths for each subcommand.
package acceptance_test

import (
	"strings"
	"testing"

	helpers "github.com/gastownhall/gascity/test/acceptance/helpers"
)

// --- gc session (bare command) ---

func TestSession_NoSubcommand_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("session")
	if err == nil {
		t.Fatal("expected error for bare 'gc session', got success")
	}
	if !strings.Contains(out, "missing subcommand") {
		t.Errorf("expected 'missing subcommand' message, got:\n%s", out)
	}
}

func TestSession_UnknownSubcommand_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("session", "explode")
	if err == nil {
		t.Fatal("expected error for unknown subcommand, got success")
	}
	if !strings.Contains(out, "unknown subcommand") {
		t.Errorf("expected 'unknown subcommand' message, got:\n%s", out)
	}
}

// --- gc session list ---

func TestSessionList_EmptyCity_ShowsNoSessions(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("session", "list")
	if err != nil {
		t.Fatalf("gc session list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "No sessions found") {
		t.Errorf("expected 'No sessions found' on fresh city, got:\n%s", out)
	}
}

func TestSessionList_JSON_ReturnsValidOutput(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("session", "list", "--json")
	if err != nil {
		t.Fatalf("gc session list --json: %v\n%s", err, out)
	}
	trimmed := strings.TrimSpace(out)
	if trimmed != "[]" && trimmed != "null" && !strings.HasPrefix(trimmed, "[") {
		t.Errorf("expected JSON array on fresh city, got:\n%s", out)
	}
}

// --- gc session prune ---

func TestSessionPrune_EmptyCity_ShowsNone(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	out, err := c.GC("session", "prune")
	if err != nil {
		t.Fatalf("gc session prune: %v\n%s", err, out)
	}
	if !strings.Contains(out, "No sessions to prune") {
		t.Errorf("expected 'No sessions to prune' on fresh city, got:\n%s", out)
	}
}

// --- gc session new (error paths) ---

func TestSessionNew_NonexistentTemplate_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	_, err := c.GC("session", "new", "nonexistent-template-xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent template, got success")
	}
}

func TestSessionNew_MissingTemplate_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	// cobra.ExactArgs(1) or custom validation handles this.
	_, err := c.GC("session", "new")
	if err == nil {
		t.Fatal("expected error for missing template, got success")
	}
}

// --- gc session close (error paths) ---

func TestSessionClose_NonexistentSession_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	_, err := c.GC("session", "close", "nonexistent-session-xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent session, got success")
	}
}

// --- gc session rename (error paths) ---

func TestSessionRename_NonexistentSession_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	_, err := c.GC("session", "rename", "nonexistent-session", "new-title")
	if err == nil {
		t.Fatal("expected error for nonexistent session, got success")
	}
}

// --- gc session peek (error paths) ---

func TestSessionPeek_NonexistentSession_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	_, err := c.GC("session", "peek", "nonexistent-session-xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent session, got success")
	}
}

// --- gc session kill (error paths) ---

func TestSessionKill_NonexistentSession_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	_, err := c.GC("session", "kill", "nonexistent-session-xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent session, got success")
	}
}

// --- gc session wake (error paths) ---

func TestSessionWake_NonexistentSession_ReturnsError(t *testing.T) {
	c := helpers.NewCity(t, testEnv)
	c.Init("claude")

	_, err := c.GC("session", "wake", "nonexistent-session-xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent session, got success")
	}
}
