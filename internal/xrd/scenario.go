package xrd

import (
	"encoding/json"
	"fmt"
	"os"
)

type Scenario struct {
	Name     string  `json:"name"`
	Lambda   float64 `json:"lambda"`
	A        float64 `json:"a,omitempty"`
	D        float64 `json:"d,omitempty"`
	Lattice  string  `json:"lattice"`
	HKL      HKL     `json:"hkl,omitempty"`
	MaxIndex int     `json:"max_index,omitempty"`
}

func LoadScenario(path string) (Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	var scenario Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return Scenario{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := ValidateWavelength(scenario.Lambda); err != nil {
		return Scenario{}, err
	}
	if scenario.D <= 0 && scenario.A <= 0 {
		return Scenario{}, fmt.Errorf("scenario needs d or a")
	}
	return scenario, nil
}

func RunScenario(scenario Scenario) (BraggResult, error) {
	d := scenario.D
	if d <= 0 {
		var err error
		d, err = LatticeSpacing(scenario.A, scenario.HKL)
		if err != nil {
			return BraggResult{}, err
		}
	}
	return BraggAngle(scenario.Lambda, d, 1)
}

func RunPowderScenario(scenario Scenario) (PowderResult, error) {
	if scenario.MaxIndex <= 0 {
		scenario.MaxIndex = 4
	}
	return Powder(scenario.Lambda, scenario.A, scenario.Lattice, scenario.MaxIndex)
}

func SaveScenario(path string, scenario Scenario) error {
	data, err := json.MarshalIndent(scenario, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func FCCCopper() Scenario {
	return Scenario{
		Name: "fcc-cu", Lambda: 1.5406, A: 3.615,
		Lattice: "fcc", HKL: HKL{H: 1, K: 1, L: 1}, MaxIndex: 4,
	}
}

func Examples() []Scenario {
	return []Scenario{
		FCCCopper(),
		{Name: "bcc-fe", Lambda: 1.5406, A: 2.866, Lattice: "bcc", HKL: HKL{H: 1, K: 1, L: 0}, MaxIndex: 3},
		{Name: "nacl", Lambda: 1.5406, A: 5.64, Lattice: "primitive", HKL: HKL{H: 2, K: 0, L: 0}, MaxIndex: 3},
	}
}

func ScenarioPaths() []string {
	return []string{"example/fcc-cu-cu.json"}
}

func ExampleLambda() float64 {
	return 1.5406
}

func ExampleA() float64 {
	return 3.615
}

func ExampleHKL() HKL {
	return HKL{H: 1, K: 1, L: 1}
}

func ExampleTwoTheta() (float64, error) {
	lambda, d, err := PrepareBraggInputs(FCCCopper())
	if err != nil {
		return 0, err
	}
	return TwoTheta(lambda, d, 1)
}

func IsFCCExampleAbove40() (bool, error) {
	twoTheta, err := ExampleTwoTheta()
	if err != nil {
		return false, err
	}
	return twoTheta > 40, nil
}

func Describe(scenario Scenario) string {
	return fmt.Sprintf("%s lambda=%.4f a=%.4f lattice=%s", scenario.Name, scenario.Lambda, scenario.A, scenario.Lattice)
}

func RunAllExamples() ([]BraggResult, error) {
	results := make([]BraggResult, 0, len(Examples()))
	for _, scenario := range Examples() {
		result, err := RunScenario(scenario)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func ExampleText() string {
	return "Cu K-alpha lambda=1.5406 A, Cu FCC a=3.615 A, (111)"
}
