package spaceship

import (
	"encoding/json"
	"fmt"

	"github.com/DNSControl/dnscontrol/v4/pkg/providers"
)

func init() {
	const providerName = "SPACESHIP"
	const providerMaintainer = "@stensonb"
	fns := providers.DspFuncs{
		Initializer: newProvider,
		RecordAuditor: AuditRecords,
	}
	// Register the provider with its activation string
	providers.RegisterDomainServiceProviderType(providerName, fns, features)
	providers.RegisterMaintainer(providerName, providerMaintainer)
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
		return nil, fmt.Errorf("missing or empty api_key")
	}
	if api.ApiSecret == "" {
		return nil, fmt.Errorf("missing or empty api_secret")
	}

	return api, nil
}
