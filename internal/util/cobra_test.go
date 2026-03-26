package util_test

import (
	"testing"

	"github.com/containrrr/shoutrrr/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().StringP("url", "u", "", "")
	cmd.Flags().StringP("message", "m", "", "")
	return cmd
}

func TestLoadFlagsFromAltSources_PositionalURL(t *testing.T) {
	cmd := newTestCmd()
	util.LoadFlagsFromAltSources(cmd, []string{"slack://token@channel"})

	got, _ := cmd.Flags().GetString("url")
	if got != "slack://token@channel" {
		t.Fatalf("expected url flag to be set from args[0], got %q", got)
	}

	msg, _ := cmd.Flags().GetString("message")
	if msg != "" {
		t.Fatalf("expected message to be empty when only one arg given, got %q", msg)
	}
}

func TestLoadFlagsFromAltSources_PositionalURLAndMessage(t *testing.T) {
	cmd := newTestCmd()
	util.LoadFlagsFromAltSources(cmd, []string{"slack://token@channel", "hello world"})

	got, _ := cmd.Flags().GetString("url")
	if got != "slack://token@channel" {
		t.Fatalf("expected url from args[0], got %q", got)
	}

	msg, _ := cmd.Flags().GetString("message")
	if msg != "hello world" {
		t.Fatalf("expected message from args[1], got %q", msg)
	}
}

func TestLoadFlagsFromAltSources_EnvURL(t *testing.T) {
	v := viper.GetViper()
	v.Set("SHOUTRRR_URL", "slack://env-token@channel")
	defer v.Set("SHOUTRRR_URL", "")

	cmd := newTestCmd()
	util.LoadFlagsFromAltSources(cmd, []string{})

	got, _ := cmd.Flags().GetString("url")
	if got != "slack://env-token@channel" {
		t.Fatalf("expected url from env, got %q", got)
	}

	msg, _ := cmd.Flags().GetString("message")
	if msg != "-" {
		t.Fatalf("expected message to default to stdin (-) when url from env, got %q", msg)
	}
}

func TestLoadFlagsFromAltSources_EnvURLWithExistingMessage(t *testing.T) {
	v := viper.GetViper()
	v.Set("SHOUTRRR_URL", "slack://env-token@channel")
	defer v.Set("SHOUTRRR_URL", "")

	cmd := newTestCmd()
	_ = cmd.Flags().Set("message", "existing")
	util.LoadFlagsFromAltSources(cmd, []string{})

	msg, _ := cmd.Flags().GetString("message")
	if msg != "existing" {
		t.Fatalf("expected existing message to be preserved, got %q", msg)
	}
}

func TestLoadFlagsFromAltSources_NoArgsNoEnv(t *testing.T) {
	v := viper.GetViper()
	v.Set("SHOUTRRR_URL", "")

	cmd := newTestCmd()
	util.LoadFlagsFromAltSources(cmd, []string{})

	got, _ := cmd.Flags().GetString("url")
	if got != "" {
		t.Fatalf("expected url to remain empty, got %q", got)
	}
}

func TestLoadFlagsFromAltSources_PositionalArgsTakePrecedenceOverEnv(t *testing.T) {
	v := viper.GetViper()
	v.Set("SHOUTRRR_URL", "slack://env-url")
	defer v.Set("SHOUTRRR_URL", "")

	cmd := newTestCmd()
	util.LoadFlagsFromAltSources(cmd, []string{"slack://arg-url"})

	got, _ := cmd.Flags().GetString("url")
	if got != "slack://arg-url" {
		t.Fatalf("expected positional arg to take precedence over env, got %q", got)
	}
}
