/*
Package providers contains definitions for OAUTH services.
*/
package providers

type ProviderStore struct {
	providers map[string]Provider
	baseUrl   string
}

// NewProviderStore() creates a provider store containing the given providers
func NewProviderStore(baseUrl string, provider ...func(baseUrl string) Provider) *ProviderStore {
	ps := &ProviderStore{
		providers: make(map[string]Provider, 0),
		baseUrl:   baseUrl,
	}

	for _, p_func := range provider {
		p := p_func(baseUrl)
		ps.providers[p.GetName()] = p
	}

	return ps
}

func (ps *ProviderStore) GetProvider(name string) (Provider, bool) {
	provider, ok := ps.providers[name]
	return provider, ok
}
