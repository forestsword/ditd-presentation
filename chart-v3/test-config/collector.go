package main

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/converter/expandconverter"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/otelcol"
)

func NewTestConfigProvider(config string) (otelcol.ConfigProvider, error) {
	return otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:       []string{"test:this"},
			Providers:  makeMapProvidersMap(NewTestProvider(config), envprovider.NewWithSettings(confmap.ProviderSettings{})),
			Converters: []confmap.Converter{expandconverter.New(confmap.ConverterSettings{})},
		},
	})
}

func NewTestCollector(configProvider otelcol.ConfigProvider) (*otelcol.Collector, error) {
	info := component.BuildInfo{}
	settings := otelcol.CollectorSettings{
		BuildInfo:      info,
		Factories:      components,
		ConfigProvider: configProvider,
	}
	return otelcol.NewCollector(settings)
}

func makeMapProvidersMap(providers ...confmap.Provider) map[string]confmap.Provider {
	ret := make(map[string]confmap.Provider, len(providers))
	for _, provider := range providers {
		ret[provider.Scheme()] = provider
	}
	return ret
}
