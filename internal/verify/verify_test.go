package verify

import (
	"testing"

	"bragg-xrd/internal/xrd"
)

func TestRunAllChecks(t *testing.T) {
	checks := RunAll()
	if !AllPass(checks) {
		for _, check := range checks {
			if !check.OK {
				t.Logf("%s: %s", check.Name, check.Message)
			}
		}
		t.Fatal("not all checks pass")
	}
}

func TestCheckEquation(t *testing.T) {
	if !CheckEquation(1.5406, 2.087, 1).OK {
		t.Fatal("equation check failed")
	}
}

func TestCheckSinThetaGate(t *testing.T) {
	if !CheckSinThetaGate(1.5406, 1).OK {
		t.Fatal("sin gate check failed")
	}
}

func TestCheckLambdaOrder(t *testing.T) {
	if !CheckLambdaOrder(1.2, 1.54, 2.5).OK {
		t.Fatal("lambda order check failed")
	}
}

func TestCheckFCCForbidden(t *testing.T) {
	if !CheckFCCForbidden(1.5406, 3.615).OK {
		t.Fatal("FCC forbidden check failed")
	}
}

func TestCheckFCC111Above40(t *testing.T) {
	if !CheckFCC111Above40().OK {
		t.Fatal("FCC 111 check failed")
	}
}

func TestCheckRightAngle(t *testing.T) {
	if !CheckRightAngle(3.0, 1.5, 1).OK {
		t.Fatal("right angle check failed")
	}
}

func TestCheckSpacingFormula(t *testing.T) {
	if !CheckSpacingFormula(3.615, xrd.HKL{H: 1, K: 1, L: 1}, 2.0871212231204974).OK {
		t.Fatal("spacing check failed")
	}
}

func TestCheckStructureFactor(t *testing.T) {
	if !CheckStructureFactor("fcc", xrd.HKL{H: 1, K: 1, L: 1}, 4).OK {
		t.Fatal("structure factor check failed")
	}
}

func TestCheckMaxOrder(t *testing.T) {
	if !CheckMaxOrder(1.5406, 2.5, 3).OK {
		t.Fatal("max order check failed")
	}
}

func TestCheckUnits(t *testing.T) {
	if !CheckUnits().OK {
		t.Fatal("units check failed")
	}
}

func TestCheckScenarioExample(t *testing.T) {
	if !CheckScenarioExample().OK {
		t.Fatal("scenario example check failed")
	}
}
