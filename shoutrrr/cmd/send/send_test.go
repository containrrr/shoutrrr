package send_test

import (
	"testing"

	"github.com/containrrr/shoutrrr/shoutrrr/cmd/send"
	"github.com/spf13/cobra"
)

func TestCmd_FlagRegistration(t *testing.T) {
	cmd := send.Cmd

	flags := []struct {
		name      string
		shorthand string
	}{
		{"url", "u"},
		{"message", "m"},
		{"title", "t"},
		{"verbose", "v"},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := cmd.Flags().Lookup(f.name)
			if flag == nil {
				t.Fatalf("expected flag %q to be registered", f.name)
			}
			if flag.Shorthand != f.shorthand {
				t.Fatalf("expected shorthand %q, got %q", f.shorthand, flag.Shorthand)
			}
		})
	}
}

func TestCmd_RequiredFlags(t *testing.T) {
	// Verify cobra still enforces required flags after the upgrade.
	root := &cobra.Command{Use: "root"}
	root.AddCommand(send.Cmd)

	root.SetArgs([]string{"send"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when required flags are missing")
	}
}

func TestCmd_MaxArgs(t *testing.T) {
	err := send.Cmd.Args(send.Cmd, []string{"a", "b", "c"})
	if err == nil {
		t.Fatal("expected error for too many args")
	}
}

func TestCmd_AcceptsTwoArgs(t *testing.T) {
	err := send.Cmd.Args(send.Cmd, []string{"url", "message"})
	if err != nil {
		t.Fatalf("expected two args to be accepted, got error: %v", err)
	}
}

func TestCmd_HasPreRun(t *testing.T) {
	// Verify PreRun is wired up (the LoadFlagsFromAltSources integration).
	// The actual behavior is tested in internal/util/cobra_test.go.
	if send.Cmd.PreRun == nil {
		t.Fatal("expected PreRun to be set for alt source flag loading")
	}
}

func TestCmd_UsesRunE(t *testing.T) {
	// Verify the command uses RunE (not Run) so errors propagate to cobra.
	if send.Cmd.RunE == nil {
		t.Fatal("expected RunE to be set for error propagation")
	}
}
