package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"

	otelV1Alpha1 "github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name       = "avgolemono"
	valuesFile = "../values.yaml"
	namespace  = "soups"
	chartPath  = ".."
)

const (
	labelNameDefaultEgressPolicy  = "use-default-egress-policy"
	labelValueDefaultEgressPolicy = "true"

	labelNameOtelOperator  = "otelOperator"
	labelValueOtelOperator = "true"

	LabelOSStable         = corev1.LabelOSStable
	LabelOSStableLinux    = "linux"
	LabelNodeRole         = "node.kubernetes.io/role"
	LabelNodeRoleWorker   = "worker"
	LabelEthosWorkload    = "ethos.corp.adobe.com/ethos-workload"
	LabelEthosWorkloadARM = "arm64"
)

const (
	grpcPort int32 = 4317
	httpPort int32 = 4318
)

type Test struct {
	Name       string
	ValuesFile string
	// SetValues is an option for render
	// see https://pkg.go.dev/github.com/gruntwork-io/terratest/modules/helm#Options
	SetValues     map[string]string
	InputTemplate string
	Namespace     string
	TestFunc      func(*testing.T, string)
}

var tests = []Test{
	{
		Name:          "validate",
		ValuesFile:    valuesFile,
		InputTemplate: "templates/collector.yaml",
		SetValues:     setValuesWithReplicas,
		Namespace:     namespace,
		TestFunc: func(t *testing.T, output string) {
			var obj otelV1Alpha1.OpenTelemetryCollector
			helm.UnmarshalK8SYaml(t, output, &obj)
			VerifyMetadata(t, obj.ObjectMeta)
			assert.Equal(t, "OpenTelemetryCollector", obj.Kind)
			assert.Equal(t, otelV1Alpha1.ModeDeployment, obj.Spec.Mode)
			assert.Equal(t, "", obj.Spec.Image)
			// We need to convert the values string to an int64
			// to be later converted in the assertion to int32
			expectedReplicas, err := strconv.ParseInt(setValuesWithReplicas["replicas"], 10, 32)
			assert.NoError(t, err)
			assert.Equal(t, int32(expectedReplicas), *obj.Spec.Replicas)
		},
	},
}

func TestTemplates(t *testing.T) {
	for _, test := range tests {
		testName := test.InputTemplate
		if test.Name != "" {
			testName = fmt.Sprintf("%s/%s", testName, test.Name)
		}
		t.Run(testName, func(t *testing.T) {
			options := &helm.Options{
				ValuesFiles: []string{test.ValuesFile},
			}
			if test.SetValues != nil {
				options.SetValues = test.SetValues
			}

			args := []string{"--namespace", test.Namespace}
			output := helm.RenderTemplate(t, options, chartPath, name, []string{test.InputTemplate}, args...)
			test.TestFunc(t, output)
		})
	}
}

func VerifyMetadata(t *testing.T, m metav1.ObjectMeta) {
	assert.Equal(t, namespace, m.Namespace)
}
