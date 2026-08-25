package xrd

import (
	"math"
	"testing"
)

func TestBraggEquation(t *testing.T) {
	result, err := BraggAngle(1.5406, 2.087, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Possible {
		t.Fatal("peak should be possible")
	}
	residual := CheckEquation(result, 0)
	if residual > 1e-9 {
		t.Fatalf("residual=%g", residual)
	}
}

func TestSinThetaGate(t *testing.T) {
	result, err := BraggAngle(1.5406, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if result.Possible {
		t.Fatal("impossible order marked possible")
	}
	if result.TwoTheta != 0 {
		t.Fatalf("2theta=%g, want 0", result.TwoTheta)
	}
}

func TestLambdaLargerRaisesTwoTheta(t *testing.T) {
	if !LambdaLargerTwoThetaHigher(1.2, 1.54, 2.5) {
		t.Fatal("larger lambda did not raise 2theta")
	}
}

func TestLargerACellLowersTwoTheta(t *testing.T) {
	if !LargerDTwoThetaLower(1.5406, 2.0, 2.5) {
		t.Fatal("larger d did not lower 2theta")
	}
}

func TestLatticeSpacingFormula(t *testing.T) {
	d, err := LatticeSpacing(3.615, HKL{H: 1, K: 1, L: 1})
	if err != nil {
		t.Fatal(err)
	}
	want := 3.615 / math.Sqrt(3)
	if math.Abs(d-want) > 1e-12 {
		t.Fatalf("d=%g, want %g", d, want)
	}
}

func TestFCCForbidden(t *testing.T) {
	forbidden, err := IsForbidden("fcc", HKL{H: 1, K: 0, L: 0})
	if err != nil {
		t.Fatal(err)
	}
	if !forbidden {
		t.Fatal("100 should be forbidden in FCC")
	}
	allowed, _ := Allowed("fcc", HKL{H: 1, K: 1, L: 1})
	if !allowed {
		t.Fatal("111 should be allowed in FCC")
	}
}

func TestBCCForbidden(t *testing.T) {
	forbidden, _ := IsForbidden("bcc", HKL{H: 1, K: 0, L: 0})
	if !forbidden {
		t.Fatal("100 should be forbidden in BCC")
	}
	allowed, _ := Allowed("bcc", HKL{H: 1, K: 1, L: 0})
	if !allowed {
		t.Fatal("110 should be allowed in BCC")
	}
}

func TestFCC111Above40(t *testing.T) {
	ok, err := IsFCCExampleAbove40()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("FCC 111 should be above 40 degrees")
	}
}

func TestAllOrders(t *testing.T) {
	results, err := AllOrders(1.5406, 2.5)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) < 2 {
		t.Fatalf("orders=%d", len(results))
	}
	for i := 1; i < len(results); i++ {
		if results[i].Theta <= results[i-1].Theta {
			t.Fatal("theta not increasing")
		}
	}
}

func TestRightAngleWhenNλ2d(t *testing.T) {
	result, err := BraggAngle(3.0, 1.5, 1)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(result.Theta-90) > 1e-9 {
		t.Fatalf("theta=%g, want 90", result.Theta)
	}
}

func TestValidationRejectsBadInputs(t *testing.T) {
	if _, err := BraggAngle(0, 1, 1); err == nil {
		t.Fatal("accepted zero wavelength")
	}
	if _, err := BraggAngle(1, -1, 1); err == nil {
		t.Fatal("accepted negative d")
	}
	if _, err := BraggAngle(1, 1, 0); err == nil {
		t.Fatal("accepted n=0")
	}
}

func TestPowderNoForbiddenFCC(t *testing.T) {
	result, err := Powder(1.5406, 3.615, "fcc", 4)
	if err != nil {
		t.Fatal(err)
	}
	if ContainsForbidden(result, "fcc") {
		t.Fatal("FCC powder contains forbidden reflections")
	}
}

func TestPowderContains111(t *testing.T) {
	result, err := Powder(1.5406, 3.615, "fcc", 4)
	if err != nil {
		t.Fatal(err)
	}
	if !HasHKL(result, HKL{H: 1, K: 1, L: 1}) {
		t.Fatal("FCC powder missing 111")
	}
}

func TestSingleCrystalSubset(t *testing.T) {
	peaks, err := SingleCrystal(1.5406, 3.615, "fcc", HKL{H: 1, K: 1, L: 1}, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(peaks) == 0 {
		t.Fatal("single crystal should have peaks")
	}
}

func TestStructureFactor(t *testing.T) {
	if f, _ := StructureFactor("fcc", HKL{H: 1, K: 0, L: 0}); f != 0 {
		t.Fatalf("FCC 100 F=%d", f)
	}
	if f, _ := StructureFactor("fcc", HKL{H: 1, K: 1, L: 1}); f != 4 {
		t.Fatalf("FCC 111 F=%d", f)
	}
}

func TestMillerHelpers(t *testing.T) {
	if NormalizeHKL(HKL{H: 2, K: 2, L: 2}) != (HKL{H: 1, K: 1, L: 1}) {
		t.Fatal("normalize failed")
	}
	if Multiplicity(HKL{H: 1, K: 1, L: 1}) != 8 {
		t.Fatalf("multiplicity=%d", Multiplicity(HKL{H: 1, K: 1, L: 1}))
	}
}

func TestUnits(t *testing.T) {
	if AngstromsToMeters(1) != 1e-10 {
		t.Fatal("A to m conversion wrong")
	}
	if MetersToAngstroms(1e-10) != 1 {
		t.Fatal("m to A conversion wrong")
	}
}
