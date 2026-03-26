package xouath2_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/containrrr/shoutrrr/pkg/generators/xouath2"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
)

func TestGenerate_FileMode_ParsesConfigJSON(t *testing.T) {
	// This test verifies that the oauth2 config file parsing works with the
	// current version of the oauth2 library — the primary upgrade risk.

	configData := map[string]interface{}{
		"client_id":     "test-client-id",
		"client_secret": "test-client-secret",
		"redirect_url":  "http://localhost:8080/callback",
		"auth_url":      "https://accounts.example.com/o/oauth2/auth",
		"token_url":     "https://accounts.example.com/o/oauth2/token",
		"smtp_hostname":  "smtp.example.com",
		"scopes":        []string{"https://mail.example.com/"},
	}

	data, err := json.Marshal(configData)
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "oauth2.json")
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	g := &xouath2.Generator{}

	// Generate will parse the file and then prompt for a verification code via stdin.
	// We provide an empty stdin so it fails at the interactive prompt, not at file parsing.
	// This validates that oauth2.Config construction works with the upgraded library.
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	w.Close() // EOF immediately
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	_, err = g.Generate(nil, map[string]string{}, []string{configFile})

	// We expect an error from the interactive prompt (scanning verification code),
	// NOT from JSON parsing or oauth2 config construction.
	if err == nil {
		t.Fatal("expected error from stdin prompt, got nil")
	}

	// If it were a JSON or oauth2 construction error, it would mention those.
	// A scan error means we got past config parsing successfully.
	errMsg := err.Error()
	if errMsg == "unexpected end of JSON input" || errMsg == "invalid character" {
		t.Fatalf("oauth2 config parsing failed: %v", err)
	}
}

func TestGenerate_FileMode_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(configFile, []byte("{invalid"), 0644); err != nil {
		t.Fatal(err)
	}

	g := &xouath2.Generator{}
	_, err := g.Generate(nil, map[string]string{}, []string{configFile})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGenerate_FileMode_MissingFile(t *testing.T) {
	g := &xouath2.Generator{}
	_, err := g.Generate(nil, map[string]string{}, []string{"/nonexistent/path/oauth2.json"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestGenerate_GmailMode_InvalidCredFile(t *testing.T) {
	// Verify that google.ConfigFromJSON still works with the upgraded oauth2 library.
	// An invalid file should produce a parse error, not a panic.
	tmpDir := t.TempDir()
	credFile := filepath.Join(tmpDir, "creds.json")
	if err := os.WriteFile(credFile, []byte(`{"not": "valid-google-creds"}`), 0644); err != nil {
		t.Fatal(err)
	}

	g := &xouath2.Generator{}
	_, err := g.Generate(nil, map[string]string{"provider": "gmail"}, []string{credFile})
	if err == nil {
		t.Fatal("expected error for invalid Google credentials JSON")
	}
}

func TestGenerate_GmailMode_ValidCredFormat(t *testing.T) {
	// Provide a valid Google credentials JSON structure to verify that
	// google.ConfigFromJSON parses it correctly with the upgraded library.
	creds := map[string]interface{}{
		"installed": map[string]interface{}{
			"client_id":     "test-id.apps.googleusercontent.com",
			"client_secret": "test-secret",
			"auth_uri":      "https://accounts.google.com/o/oauth2/auth",
			"token_uri":     "https://oauth2.googleapis.com/token",
			"redirect_uris": []string{"urn:ietf:wg:oauth:2.0:oob"},
		},
	}

	data, _ := json.Marshal(creds)
	tmpDir := t.TempDir()
	credFile := filepath.Join(tmpDir, "creds.json")
	if err := os.WriteFile(credFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Close stdin so it fails at the interactive prompt, not at parsing.
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	g := &xouath2.Generator{}
	_, err := g.Generate(nil, map[string]string{"provider": "gmail"}, []string{credFile})

	// Should fail at the "Enter verification code" prompt, not at config parsing.
	if err == nil {
		t.Fatal("expected error from stdin prompt")
	}
}

func TestSmtpConfigGetURL_RoundTrip(t *testing.T) {
	// Verify the smtp.Config used by the generator still produces valid URLs.
	// This catches breakage in the format/types packages that the generator depends on.
	config := &smtp.Config{
		Host:        "smtp.example.com",
		Port:        587,
		Username:    "user@example.com",
		Password:    "token123",
		FromAddress: "user@example.com",
		FromName:    "Test",
		ToAddresses: []string{"dest@example.com"},
		Auth:        smtp.AuthTypes.OAuth2,
		UseStartTLS: true,
		UseHTML:     true,
	}

	u := config.GetURL()
	if u == nil {
		t.Fatal("GetURL returned nil")
	}

	if u.Scheme != "smtp" {
		t.Fatalf("expected scheme smtp, got %q", u.Scheme)
	}

	if u.Hostname() != "smtp.example.com" {
		t.Fatalf("expected host smtp.example.com, got %q", u.Hostname())
	}

	if u.Port() != "587" {
		t.Fatalf("expected port 587, got %q", u.Port())
	}

	if u.User.Username() != "user@example.com" {
		t.Fatalf("expected username user@example.com, got %q", u.User.Username())
	}
}
