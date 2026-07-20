package spaceship

import (
	"encoding/json"
	"fmt"

	"github.com/DNSControl/dnscontrol/v4/pkg/providers"
	"github.com/namecheap/go-spaceship-sdk/client"
)

func init() {
	const providerName = "SPACESHIP"
	const providerMaintainer = "@stensonb"
	fns := providers.DspFuncs{
		Initializer:   NewProvider,
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
	client    *client.Client
	apiKey    string
	apiSecret string
	baseURL   string
}

func NewProvider(config map[string]string, metadata json.RawMessage) (providers.DNSServiceProvider, error) {
	apiKey := config["api_key"]
	apiSecret := config["api_secret"]
	baseURL := config["base_url"]

	if apiKey == "" {
		return nil, fmt.Errorf("missing or empty api_key")
	}

	if apiSecret == "" {
		return nil, fmt.Errorf("missing or empty api_secret")
	}

	// set default if not specified
	if baseURL == "" {
		baseURL = "https://spaceship.dev/api/v1"
	}

	client, err := client.NewClient(baseURL, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to build client: %w", err)
	}

	return &spaceshipProvider{
		client:    client,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
	}, nil
}
