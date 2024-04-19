package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

const (
	name       = "avgolemono"
	valuesFile = "../values.yaml"
	namespace  = "soups"
	chartPath  = ".."
)

type OpenTelemetrySpec struct {
	Config string `json:"config,omitempty"`
}

type OpenTelemetryCollector struct {
	Spec OpenTelemetrySpec `json:"spec,omitempty"`
}

type Test struct {
	Name          string
	ValuesFile    string
	InputTemplate string
	SetValues     map[string]string
	Namespace     string
	TestFunc      func(*testing.T, *otelcol.Config)
	EnvVars       map[string]string
}

var tests = []Test{
	{
		Name:          "validate",
		ValuesFile:    valuesFile,
		InputTemplate: "templates/collector.yaml",
		SetValues:     setValuesDefault,
		Namespace:     namespace,
		EnvVars: map[string]string{
			"OTELCOL_HTTP_ENDPOINT": "192.168.1.1:4318",
		},
		TestFunc: func(t *testing.T, config *otelcol.Config) {
			otlpReceiver := config.Receivers[component.MustNewID("otlp")].(*otlpreceiver.Config)
			assert.Equal(t, "0.0.0.0:4317", otlpReceiver.Protocols.GRPC.NetAddr.Endpoint)
			assert.Equal(t, "192.168.1.1:4318", otlpReceiver.Protocols.HTTP.ServerConfig.Endpoint)
		},
	},
}

func TestConfig(t *testing.T) {
	for _, test := range tests {
		testName := test.InputTemplate
		if test.Name != "" {
			testName = fmt.Sprintf("%s/%s", testName, test.Name)
		}
		t.Run(testName, func(t *testing.T) {
			for key, value := range test.EnvVars {
				t.Setenv(key, value)
			}
			options := &helm.Options{
				ValuesFiles: []string{test.ValuesFile},
			}
			if test.SetValues != nil {
				options.SetValues = test.SetValues
			}

			args := []string{"--namespace", test.Namespace}
			output := helm.RenderTemplate(t, options, chartPath, name, []string{test.InputTemplate}, args...)
			var obj OpenTelemetryCollector
			helm.UnmarshalK8SYaml(t, output, &obj)

			configProvider, err := NewTestConfigProvider(obj.Spec.Config)
			require.NoError(t, err)

			// See https://github.com/open-telemetry/opentelemetry-collector/blob/main/otelcol/command_validate.go
			collector, err := NewTestCollector(configProvider)
			require.NoError(t, err)

			err = collector.DryRun(context.Background()) // Basic validation
			require.NoError(t, err)

			factories, err := components()
			require.NoError(t, err)

			resolvedConfig, err := configProvider.Get(context.Background(), factories)
			require.NoError(t, err)

			test.TestFunc(t, resolvedConfig) // Deep validation
		})
	}
}
