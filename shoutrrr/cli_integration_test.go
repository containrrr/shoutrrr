package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binary string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "shoutrrr-test-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	binary = filepath.Join(tmp, "shoutrrr")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		panic("failed to build binary: " + string(out))
	}

	os.Exit(m.Run())
}

func mustAbs(rel string) string {
	abs, err := filepath.Abs(rel)
	if err != nil {
		panic(err)
	}
	return abs
}

func run(args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binary, args...)
	// Clear environment to avoid inheriting SHOUTRRR_URL from the test runner
	cmd.Env = append(os.Environ(), "NO_COLOR=1")

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = -1
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func runWithEnv(env []string, args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(), append(env, "NO_COLOR=1")...)

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = -1
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// --- send command ---

func TestSend_LoggerService(t *testing.T) {
	stdout, stderr, code := run("send", "-u", "logger://", "-m", "integration test")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "Notification sent") {
		t.Fatalf("expected 'Notification sent' in stderr, got: %s", stderr)
	}
}

func TestSend_InvalidService(t *testing.T) {
	_, stderr, code := run("send", "-u", "invalid://foo", "-m", "test")
	if code == 0 {
		t.Fatal("expected nonzero exit for invalid service")
	}
	if !strings.Contains(stderr, "unknown service") {
		t.Fatalf("expected 'unknown service' in stderr, got: %s", stderr)
	}
}

func TestSend_MissingFlags(t *testing.T) {
	_, _, code := run("send")
	if code != 64 {
		t.Fatalf("expected exit code 64 (ExUsage), got %d", code)
	}
}

func TestSend_PositionalArgs(t *testing.T) {
	_, stderr, code := run("send", "logger://", "positional message")
	if code != 0 {
		t.Fatalf("expected exit 0 for positional args, got %d\nstderr: %s", code, stderr)
	}
}

func TestSend_EnvURL(t *testing.T) {
	// Provide message via flag, URL via env. Stdin is /dev/null so "-" default would fail.
	_, stderr, code := runWithEnv(
		[]string{"SHOUTRRR_URL=logger://"},
		"send", "-m", "env test",
	)
	if code != 0 {
		t.Fatalf("expected exit 0 with SHOUTRRR_URL, got %d\nstderr: %s", code, stderr)
	}
}

func TestSend_Verbose(t *testing.T) {
	_, stderr, code := run("send", "-u", "logger://", "-m", "verbose test", "-v")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr: %s", code, stderr)
	}
	if !strings.Contains(stderr, "URLs:") {
		t.Fatalf("expected verbose output with 'URLs:' prefix, got: %s", stderr)
	}
}

func TestSend_Title(t *testing.T) {
	_, stderr, code := run("send", "-u", "logger://", "-m", "titled", "-t", "Test Title")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr: %s", code, stderr)
	}
}

func TestSend_TooManyArgs(t *testing.T) {
	_, _, code := run("send", "a", "b", "c")
	if code == 0 {
		t.Fatal("expected nonzero exit for too many args")
	}
}

// --- verify command ---

func TestVerify_LoggerService(t *testing.T) {
	_, _, code := run("verify", "-u", "logger://")
	if code != 0 {
		t.Fatalf("expected exit 0 for logger verify, got %d", code)
	}
}

func TestVerify_DiscordOutputFormat(t *testing.T) {
	stdout, _, code := run("verify", "-u", "discord://mytoken@mywebhookid")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	// Verify the tabular config format contains expected fields and values
	expectedEntries := []struct {
		field string
		value string
	}{
		{"Token", "mytoken"},
		{"WebhookID", "mywebhookid"},
		{"Color", "0x50d9ff"},
		{"SplitLines", "Yes"},
		{"JSON", "No"},
	}
	for _, e := range expectedEntries {
		if !strings.Contains(stdout, e.field) {
			t.Errorf("missing field %q in verify output:\n%s", e.field, stdout)
		}
		if !strings.Contains(stdout, e.value) {
			t.Errorf("missing value %q for field %q in verify output:\n%s", e.value, e.field, stdout)
		}
	}

	// Verify structural markers
	if !strings.Contains(stdout, "<Required>") {
		t.Errorf("missing '<Required>' marker in verify output:\n%s", stdout)
	}
	if !strings.Contains(stdout, "<Default:") {
		t.Errorf("missing '<Default:' marker in verify output:\n%s", stdout)
	}
}

func TestVerify_InvalidService(t *testing.T) {
	_, _, code := run("verify", "-u", "invalid://foo")
	if code == 0 {
		t.Fatal("expected nonzero exit for invalid service verify")
	}
}

func TestVerify_PositionalArg(t *testing.T) {
	_, _, code := run("verify", "logger://")
	if code != 0 {
		t.Fatalf("expected exit 0 for positional arg verify, got %d", code)
	}
}

func TestVerify_MissingURL(t *testing.T) {
	_, _, code := run("verify")
	if code == 0 {
		t.Fatal("expected nonzero exit when URL is missing")
	}
}

// --- docs command ---

func TestDocs_LoggerMarkdown(t *testing.T) {
	stdout, _, code := run("docs", "-f", "markdown", "logger")
	if code != 0 {
		t.Fatalf("expected exit 0 for docs logger, got %d", code)
	}
	if !strings.Contains(stdout, "### URL Fields") {
		t.Fatalf("expected '### URL Fields' section header, got: %s", stdout)
	}
	if !strings.Contains(stdout, "### Query/Param Props") {
		t.Fatalf("expected '### Query/Param Props' section header, got: %s", stdout)
	}
}

func TestDocs_LoggerConsole(t *testing.T) {
	_, _, code := run("docs", "-f", "console", "logger")
	if code != 0 {
		t.Fatalf("expected exit 0 for docs console, got %d", code)
	}
}

func TestDocs_InvalidFormat(t *testing.T) {
	_, _, code := run("docs", "-f", "invalid", "logger")
	if code == 0 {
		t.Fatal("expected nonzero exit for invalid format")
	}
}

func TestDocs_MissingService(t *testing.T) {
	_, _, code := run("docs")
	if code == 0 {
		t.Fatal("expected nonzero exit when no service specified")
	}
}

func TestDocs_SlackMarkdownFormat(t *testing.T) {
	stdout, _, code := run("docs", "-f", "markdown", "slack")
	if code != 0 {
		t.Fatalf("expected exit 0 for docs slack, got %d", code)
	}

	// Verify structural sections are present
	if !strings.Contains(stdout, "### URL Fields") {
		t.Fatalf("missing '### URL Fields' section:\n%s", stdout)
	}
	if !strings.Contains(stdout, "### Query/Param Props") {
		t.Fatalf("missing '### Query/Param Props' section:\n%s", stdout)
	}

	// Verify known fields are documented with expected formatting
	expectedFields := []string{
		"__Token__",
		"__Channel__",
		"__BotName__",
		"__Color__",
		"__Icon__",
		"__Title__",
		"__ThreadTS__",
	}
	for _, field := range expectedFields {
		if !strings.Contains(stdout, field) {
			t.Errorf("missing field %s in slack docs output", field)
		}
	}

	// Verify URL part format uses service-url code blocks
	if !strings.Contains(stdout, `<code class="service-url">`) {
		t.Errorf("missing service-url code block formatting:\n%s", stdout)
	}
	if !strings.Contains(stdout, "slack://") {
		t.Errorf("missing slack:// URL scheme in docs:\n%s", stdout)
	}

	// Verify props description blurb
	if !strings.Contains(stdout, "?key=value&key=value") {
		t.Errorf("missing query param usage hint:\n%s", stdout)
	}
}

func TestDocs_DiscordMarkdownFormat(t *testing.T) {
	stdout, _, code := run("docs", "-f", "markdown", "discord")
	if code != 0 {
		t.Fatalf("expected exit 0 for docs discord, got %d", code)
	}

	expectedFields := []string{
		"__Token__",
		"__WebhookID__",
		"__Avatar__",
		"__Color__",
		"__JSON__",
		"__SplitLines__",
		"__ThreadID__",
		"__Username__",
	}
	for _, field := range expectedFields {
		if !strings.Contains(stdout, field) {
			t.Errorf("missing field %s in discord docs output", field)
		}
	}

	// Verify default values are rendered
	if !strings.Contains(stdout, "0x50D9ff") {
		t.Errorf("missing default Color value '0x50D9ff' in discord docs:\n%s", stdout)
	}
	if !strings.Contains(stdout, "discord://") {
		t.Errorf("missing discord:// URL scheme in docs:\n%s", stdout)
	}
}

// --- generate command ---

func TestGenerate_MissingService(t *testing.T) {
	_, _, code := run("generate")
	if code == 0 {
		t.Fatal("expected nonzero exit when no service specified")
	}
}

func TestGenerate_InvalidService(t *testing.T) {
	_, _, code := run("generate", "-s", "nonexistent")
	if code == 0 {
		t.Fatal("expected nonzero exit for nonexistent service")
	}
}

// --- root command ---

func TestRoot_Version(t *testing.T) {
	stdout, _, code := run("--version")
	if code != 0 {
		t.Fatalf("expected exit 0 for --version, got %d", code)
	}
	if !strings.Contains(stdout, "shoutrrr") {
		t.Fatalf("expected version output to contain 'shoutrrr', got: %s", stdout)
	}
}

func TestRoot_Help(t *testing.T) {
	stdout, _, code := run("--help")
	if code != 0 {
		t.Fatalf("expected exit 0 for --help, got %d", code)
	}
	// All subcommands should be listed
	for _, sub := range []string{"send", "verify", "docs", "generate"} {
		if !strings.Contains(stdout, sub) {
			t.Fatalf("expected help to list %q subcommand, got: %s", sub, stdout)
		}
	}
}

func TestRoot_UnknownCommand(t *testing.T) {
	_, _, code := run("nonexistent-command")
	if code == 0 {
		t.Fatal("expected nonzero exit for unknown command")
	}
}
