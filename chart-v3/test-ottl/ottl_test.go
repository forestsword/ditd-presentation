package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"gopkg.in/yaml.v2"
)

type Values struct {
	OTTLStatements OTTLStatements `yaml:"ottlStatements"`
}

type OTTLStatements struct {
	ContextSpan []Scenario `yaml:"contextSpan"`
}

type Scenario struct {
	Name        string
	Function    string
	Tests       []ScenarioTest
	target      ottl.GetSetter[pcommon.Value]
	Regex       string
	Replacement string
	optFunction ottl.Optional[ottl.FunctionGetter[pcommon.Value]]
	want        func(pcommon.Value)
}

type ScenarioTest struct {
	Input  string
	Expect string
}

func TestReplacements(t *testing.T) {
	stdFunctions := ottlfuncs.StandardFuncs[pcommon.Value]()
	target := &ottl.StandardGetSetter[pcommon.Value]{
		Getter: func(ctx context.Context, tCtx pcommon.Value) (any, error) {
			return tCtx.Str(), nil
		},
		Setter: func(ctx context.Context, tCtx pcommon.Value, val any) error {
			tCtx.SetStr(val.(string))
			return nil
		},
	}
	content, err := os.ReadFile("../values.yaml")
	if err != nil {
		t.Fatal("cannot read values file")
	}
	var values Values
	err = yaml.Unmarshal(content, &values)
	if err != nil {
		t.Fatal("cannot unmarshal values")
	}

	for _, scenario := range values.OTTLStatements.ContextSpan {
		for _, test := range scenario.Tests {
			t.Run(fmt.Sprintf("%s/%s", scenario.Name, test.Expect), func(t *testing.T) {
				scenarioValue := pcommon.NewValueStr(test.Input)
				replacePatternFunctionFactory := stdFunctions[scenario.Function]
				assert.Equal(t, scenario.Function, replacePatternFunctionFactory.Name())
				args := replacePatternFunctionFactory.CreateDefaultArguments()
				argsReal, ok := args.(*ottlfuncs.ReplacePatternArguments[pcommon.Value])
				if !ok {
					t.Fatal("not ok")
				}
				argsReal.Target = target
				argsReal.RegexPattern = scenario.Regex
				argsReal.Replacement = ottl.StandardStringGetter[pcommon.Value]{
					Getter: func(context.Context, pcommon.Value) (any, error) {
						return scenario.Replacement, nil
					},
				}
				argsReal.Function = scenario.optFunction
				fctx := ottl.FunctionContext{
					Set: componenttest.NewNopTelemetrySettings(),
				}

				exprFunc, err := replacePatternFunctionFactory.CreateFunction(fctx, argsReal)
				assert.NoError(t, err)

				result, err := exprFunc(nil, scenarioValue)
				assert.NoError(t, err)
				assert.Nil(t, result)

				expected := pcommon.NewValueStr(test.Expect)
				assert.Equal(t, expected, scenarioValue)
			})
		}
	}
}
