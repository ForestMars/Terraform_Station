package factory_test

import (
	"testing"

	"github.com/ForestMars/TerraformStation"
	"github.com/ForestMars/TerraformStation/factory"
)

var tfCommandCases = []struct {
	name     string
	input    *TerraformStation.TFCommandInput
	expected *TerraformStation.TFCommandResult
}{
	{
		name: "terraform init",
		input: &TerraformStation.TFCommandInput{
			Command: "init,
		},
		expected: &TerraformStation.TFCommandResult{
			Result: "result",
		},
	},
	{
		name: "terraform plan,
		input: &TerraformStation.TFCommandInput{
			Command: "plan",
		},
		expected: &TerraformStation.TFCommandResult{
			Result: "result",
		},
	},
}

func TestTFComannd(t *testing.T) {
	serv := factory.New(nil, &TerraformStation.Config{Option1: "opt1"}).OptionalParam1("param1").OptionalParam2("param2")
	for _, c := range tfCommandCases {
		result, err := serv.TFCommand(c.input)
		if err != nil {
			panic(err)
		}

		diff := testingutils.PrettyJsonDiff(c.expected, result)
		if len(diff) > 0 {
			t.Error(c.name, diff)
		}
	}
}
