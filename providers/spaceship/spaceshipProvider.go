package spaceship

import (
	"encoding/json"
	"fmt"

	"github.com/DNSControl/dnscontrol/v4/pkg/providers"
)

const providerName = "SPACESHIP"

func init() {
	// Register the provider with its activation string
	providers.RegisterDomainServiceProviderType(providerName, newProvider, features)
}

var features = providers.DocumentationNotes{
	providers.CanUseAlias:            providers.Cannot(),
	providers.CanUseCAA:              providers.Can(),
	providers.CanUseSRV:              providers.Can(),
	providers.CanUsePTR:              providers.Cannot(),
	providers.DocDualHost:            providers.Can(),
	providers.DocOfficiallySupported: providers.Cannot(),
}

type spaceshipProvider struct {
	ApiKey    string
	ApiSecret string
}

func newProvider(config map[string]string, metadata json.RawMessage) (providers.DNSServiceProvider, error) {
	api := &spaceshipProvider{}

	api.ApiKey = config["api_key"]
	api.ApiSecret = config["api_secret"]

	if api.ApiKey == "" {
		return nil, fmt.Error("missing or empty api_key")
	}
	if api.ApiSecret == "" {
		return nil, fmt.Error("missing or empty api_secret")
	}

	return api, nil
}
