package cli

import (
	"testing"

	"bragg-xrd/internal/xrd"
)

func TestParseHKL(t *testing.T) {
	hkl, err := parseHKL("111")
	if err != nil {
		t.Fatal(err)
	}
	if hkl != (xrd.HKL{H: 1, K: 1, L: 1}) {
		t.Fatalf("hkl = %+v", hkl)
	}
}

func TestValidateWavelength(t *testing.T) {
	if err := validateWavelength(1.5406); err != nil {
		t.Fatal(err)
	}
	if err := validateWavelength(0); err == nil {
		t.Fatal("accepted zero wavelength")
	}
}

func TestValidateLattice(t *testing.T) {
	if err := validateLattice("fcc"); err != nil {
		t.Fatal(err)
	}
	if err := validateLattice("hex"); err == nil {
		t.Fatal("accepted hex lattice")
	}
}

func TestLoadExample(t *testing.T) {
	scenario, err := loadExample("../../example/fcc-cu-cu.json")
	if err != nil {
		t.Fatal(err)
	}
	if scenario.Lambda != 1.5406 || scenario.A != 3.615 {
		t.Fatalf("scenario = %+v", scenario)
	}
}

func TestExamplePaths(t *testing.T) {
	if len(ExamplePaths()) != 1 {
		t.Fatal("expected one example path")
	}
}

func TestEvaluateExample(t *testing.T) {
	scenario, _ := loadExample("../../example/fcc-cu-cu.json")
	result, err := xrd.RunScenario(scenario)
	if err != nil {
		t.Fatal(err)
	}
	if result.TwoTheta < 40 {
		t.Fatalf("2theta=%g", result.TwoTheta)
	}
}
