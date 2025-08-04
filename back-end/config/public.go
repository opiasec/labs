package config

type PublicConfig struct {
	PasswordAuthEnabled  bool           `json:"passwordAuthEnabled"`
	ProviderLoginEnabled bool           `json:"providerLoginEnabled"`
	ProviderConfig       ProviderConfig `json:"providerConfig,omitempty"`
}

type ProviderConfig struct {
	Name             string `json:"name"`
	AuthorizationURL string `json:"authorizationUrl"`
	TokenURL         string `json:"tokenUrl"`
	Audience         string `json:"audience"`
	ClientID         string `json:"clientId"`
	Scope            string `json:"scope"`
}
