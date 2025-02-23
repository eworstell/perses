// Copyright 2021 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dashboard

import (
	"fmt"
	"testing"

	"github.com/perses/perses/pkg/model/api/v1/common"
	"github.com/perses/perses/pkg/model/api/v1/variable"
	"github.com/stretchr/testify/assert"
)

func TestBuildVariableDependencies(t *testing.T) {
	testSuite := []struct {
		title     string
		variables []Variable
		result    map[string][]string
	}{
		{
			title:     "no variable, not dep",
			variables: nil,
			result:    map[string][]string{},
		},
		{
			title: "constant variable, no dep",
			variables: []Variable{
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "myVariable",
					},
				},
			},
			result: map[string][]string{},
		},
		{
			title: "query variable with no variable used",
			variables: []Variable{
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "vector(1)",
								},
							},
						},
						Name: "myVariable",
					},
				},
			},
			result: map[string][]string{},
		},
		{
			title: "query variable with variable used",
			variables: []Variable{
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "sum by($doe) (rate($foo{label='$bar'}))",
								},
							},
						},
						Name: "myVariable",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "test",
								},
							},
						},
						Name: "foo",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "vector($foo)",
								},
							},
						},
						Name: "bar",
					},
				},
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "doe",
					},
				},
			},
			result: map[string][]string{
				"myVariable": {
					"doe", "foo", "bar",
				},
				"bar": {
					"foo",
				},
			},
		},
		{
			title: "query variable label_values with variable used",
			variables: []Variable{
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusLabelValuesVariable",
								Spec: map[string]interface{}{
									"label_name": "$foo",
									"matchers": []interface{}{
										"$foo{$bar='test'}",
									},
								},
							},
						},
						Name: "myVariable",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "test",
								},
							},
						},
						Name: "foo",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "vector($foo)",
								},
							},
						},
						Name: "bar",
					},
				},
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "doe",
					},
				},
			},
			result: map[string][]string{
				"myVariable": {
					"foo", "bar",
				},
				"bar": {
					"foo",
				},
			},
		},
		{
			title: "multiple usage of the same variable",
			variables: []Variable{
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "sum by($doe, $bar) (rate($foo{label='$bar'}))",
								},
							},
						},
						Name: "myVariable",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "test",
								},
							},
						},
						Name: "foo",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "vector($foo)",
								},
							},
						},
						Name: "bar",
					},
				},
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "doe",
					},
				},
			},
			result: map[string][]string{
				"myVariable": {
					"doe", "bar", "foo",
				},
				"bar": {
					"foo",
				},
			},
		},
		{
			title: "variable with only number is ignored",
			variables: []Variable{
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "filter_platform",
					},
				},
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "PaaS",
					},
				},
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "filter_kube_sts",
					},
				},
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "extlabels_prometheus_namespace",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr":       "group by(prometheus) (label_replace(kube_statefulset_labels{$filter_platform,stack=~\"$PaaS\",$filter_kube_sts,stack=~\"$PaaS\",namespace=~\"$extlabels_prometheus_namespace\"},\"prometheus\",\"$1\",\"label_app_kubernetes_io_instance\",\"([^-]+)-?.*\"))",
									"label_name": "prometheus",
								},
							},
						},
						Name: "foo",
					},
				},
			},
			result: map[string][]string{
				"foo": {
					"filter_platform", "PaaS", "filter_kube_sts", "extlabels_prometheus_namespace",
				},
			},
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			result, err := buildVariableDependencies(test.variables)
			assert.NoError(t, err)
			assert.Equal(t, len(test.result), len(result))
			for k, v := range test.result {
				assert.ElementsMatch(t, v, result[k])
			}
		})
	}
}

func TestBuildVariableDependenciesError(t *testing.T) {
	testSuite := []struct {
		title     string
		variables []Variable
		err       error
	}{
		{
			title: "variable used but not defined",
			variables: []Variable{
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "sum by($doe, $bar) (rate($foo{label='$bar'}))",
								},
							},
						},
						Name: "myVariable",
					},
				},
			},
			err: fmt.Errorf("variable %q is used in the variable %q but not defined", "doe", "myVariable"),
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			_, err := buildVariableDependencies(test.variables)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestGraph_BuildOrder(t *testing.T) {
	testSuite := []struct {
		title        string
		variables    []string
		dependencies map[string][]string
		result       []VariableGroup
	}{
		{
			title:     "single variable",
			variables: []string{"myVariable"},
			result:    []VariableGroup{{Variables: []string{"myVariable"}}},
		},
		{
			title:     "independent variable",
			variables: []string{"a", "d", "e"},
			result:    []VariableGroup{{Variables: []string{"a", "d", "e"}}},
		},
		{
			title:     "a depend on d depend on e",
			variables: []string{"a", "d", "e"},
			dependencies: map[string][]string{
				"a": {"d"},
				"d": {"e"},
			},
			result: []VariableGroup{
				{Variables: []string{"e"}},
				{Variables: []string{"d"}},
				{Variables: []string{"a"}},
			},
		},
		{
			title:     "complete dep graph",
			variables: []string{"f", "d", "c", "b", "g", "a", "h", "e"},
			dependencies: map[string][]string{
				"e": {"a", "b"},
				"a": {"c", "f", "b"},
				"h": {"b"},
				"g": {"d"},
				"c": {"f"},
				"b": {"f"},
			},
			result: []VariableGroup{
				{Variables: []string{"f", "d"}},
				{Variables: []string{"c", "b", "g"}},
				{Variables: []string{"a", "h"}},
				{Variables: []string{"e"}},
			},
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			g := newGraph(test.variables, test.dependencies)
			result, err := g.buildOrder()
			assert.NoError(t, err)
			assert.Equal(t, len(test.result), len(result))
			for i := 0; i < len(result); i++ {
				assert.ElementsMatch(t, test.result[i].Variables, result[i].Variables)
			}
		})
	}
}

func TestGraph_BuildOrderError(t *testing.T) {
	testSuite := []struct {
		title        string
		variables    []string
		dependencies map[string][]string
	}{
		{
			title:     "simple circular dep",
			variables: []string{"a", "b"},
			dependencies: map[string][]string{
				"a": {"b"},
				"b": {"a"},
			},
		},
		{
			title:     "circular dep on the same node",
			variables: []string{"a"},
			dependencies: map[string][]string{
				"a": {"a"},
			},
		},
		{
			title:     "circular dep with transition",
			variables: []string{"f", "d", "c", "b", "g", "a", "h", "e"},
			dependencies: map[string][]string{
				"e": {"a", "b"},
				"a": {"c", "f", "b"},
				"h": {"b"},
				"g": {"d", "c"},
				"c": {"f"},
				"b": {"f"},
				"d": {"d"},
			},
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			g := newGraph(test.variables, test.dependencies)
			_, err := g.buildOrder()
			assert.Equal(t, fmt.Errorf("circular dependency detected"), err)
		})
	}
}

func TestBuildOrder(t *testing.T) {
	testSuite := []struct {
		title     string
		variables []Variable
		result    []VariableGroup
	}{
		{
			title: "no variable",
		},
		{
			title: "constant variable, no dep",
			variables: []Variable{
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "myVariable",
					},
				},
			},
			result: []VariableGroup{{Variables: []string{"myVariable"}}},
		},
		{
			title: "multiple usage of same variable",
			variables: []Variable{
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "sum by($doe, $bar) (rate($foo{label='$bar'}))",
								},
							},
						},
						Name: "myVariable",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "test",
								},
							},
						},
						Name: "foo",
					},
				},
				{
					Kind: variable.KindList,
					Spec: &ListVariableSpec{
						ListSpec: variable.ListSpec{
							Plugin: common.Plugin{
								Kind: "PrometheusPromQLVariable",
								Spec: map[string]interface{}{
									"expr": "vector($foo)",
								},
							},
						},
						Name: "bar",
					},
				},
				{
					Kind: variable.KindText,
					Spec: &TextVariableSpec{
						TextSpec: variable.TextSpec{
							Value: "myConstant",
						},
						Name: "doe",
					},
				},
			},
			result: []VariableGroup{
				{Variables: []string{"doe", "foo"}},
				{Variables: []string{"bar"}},
				{Variables: []string{"myVariable"}},
			},
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			groups, err := BuildVariableOrder(test.variables)
			assert.NoError(t, err)
			assert.Equal(t, len(test.result), len(groups))
			for i := 0; i < len(groups); i++ {
				assert.ElementsMatch(t, test.result[i].Variables, groups[i].Variables)
			}
		})
	}
}
