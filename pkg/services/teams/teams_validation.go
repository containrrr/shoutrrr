package teams

import (
	"fmt"
	"regexp"
)

const (
	UUID4Length        = 36
	HashLength         = 32
	WebhookDomain      = ".webhook.office.com"
	ExpectedComponents = 7 // 1 match + 6 captures
	Path               = "webhookb2"
	ProviderName       = "IncomingWebhook"
)

// parseAndVerifyWebhookURL extracts and validates webhook components from a URL.
func parseAndVerifyWebhookURL(webhookURL string) ([5]string, error) {
	pattern, err := regexp.Compile(
		`https://([^.]+)` + WebhookDomain + `/` + Path + `/([0-9a-f-]{36})@([0-9a-f-]{36})/` + ProviderName + `/([0-9a-f]{32})/([0-9a-f-]{36})/([^/]+)`,
	)
	if err != nil {
		return [5]string{}, err
	}
	groups := pattern.FindStringSubmatch(webhookURL)
	if len(groups) != ExpectedComponents {
		return [5]string{}, fmt.Errorf(
			"invalid webhook URL format: expected %d components, got %d",
			ExpectedComponents,
			len(groups),
		)
	}
	parts := [5]string{groups[2], groups[3], groups[4], groups[5], groups[6]}
	if err := verifyWebhookParts(parts); err != nil {
		return [5]string{}, err
	}
	return parts, nil
}

// verifyWebhookParts ensures webhook components meet format requirements.
// It checks lengths of Group, Tenant, AltID, and GroupOwner, and ensures ExtraID is present.
// Returns an error if any component is invalid.
func verifyWebhookParts(parts [5]string) error {
	type partSpec struct {
		name     string
		length   int
		index    int
		optional bool
	}
	specs := []partSpec{
		{name: "group ID", length: UUID4Length, index: 0, optional: true},
		{name: "tenant ID", length: UUID4Length, index: 1, optional: true},
		{name: "altID", length: HashLength, index: 2, optional: true},
		{name: "groupOwner", length: UUID4Length, index: 3, optional: true},
	}

	for _, spec := range specs {
		if len(parts[spec.index]) != spec.length && parts[spec.index] != "" {
			return fmt.Errorf(
				"%s must be %d characters, got %d",
				spec.name,
				spec.length,
				len(parts[spec.index]),
			)
		}
	}

	if parts[4] == "" {
		return fmt.Errorf("extraID is required")
	}
	return nil
}

// BuildWebhookURL constructs a Teams webhook URL from components.
func BuildWebhookURL(host, group, tenant, altID, groupOwner, extraID string) string {
	return fmt.Sprintf("https://%s/%s/%s@%s/%s/%s/%s/%s",
		host, Path, group, tenant, ProviderName, altID, groupOwner, extraID)
}
