package patterns

import (
	"testing"

	"bragg-xrd/internal/xrd"
)

func TestBuildPattern(t *testing.T) {
	pattern, err := Build(1.5406, 3.615, "fcc", 4)
	if err != nil {
		t.Fatal(err)
	}
	if !IsNonEmpty(pattern) {
		t.Fatal("pattern empty")
	}
}

func TestNormalize(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	pattern = Normalize(pattern)
	if MaxIntensity(pattern) != 100 {
		t.Fatalf("max=%g", MaxIntensity(pattern))
	}
}

func TestNoForbidden(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	if !IsClean(pattern, "fcc") {
		t.Fatal("FCC pattern contains forbidden peaks")
	}
}

func TestHas111(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	if !HasPeak(pattern, xrd.HKL{H: 1, K: 1, L: 1}) {
		t.Fatal("missing 111")
	}
}

func TestSortAndMonotonic(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	pattern = Sort(pattern)
	if !IsMonotonic(pattern) {
		t.Fatal("pattern not sorted")
	}
}

func TestStrongest(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	if len(Strongest(pattern, 3).Peaks) != 3 {
		t.Fatal("strongest count wrong")
	}
}

func TestSummaryAndText(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	if Summary(pattern) == "" || Text(pattern) == "" || CSV(pattern) == "" {
		t.Fatal("empty output helpers")
	}
}

func TestPeakAtIndex(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	peak, err := PeakAtIndex(pattern, 0)
	if err != nil {
		t.Fatal(err)
	}
	if peak.TwoTheta <= 0 {
		t.Fatal("peak angle non-positive")
	}
}

func TestAngleStats(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	if MeanAngle(pattern) <= 0 || MedianAngle(pattern) <= 0 {
		t.Fatal("angle stats invalid")
	}
	if StddevAngle(pattern) < 0 {
		t.Fatal("negative stddev")
	}
}

func TestStrongestPeak(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	peak := StrongestPeak(pattern)
	if peak.Intensity <= 0 {
		t.Fatal("strongest peak non-positive intensity")
	}
}

func TestFamilyLabels(t *testing.T) {
	pattern, _ := Build(1.5406, 3.615, "fcc", 4)
	if len(FamilyLabels(pattern)) == 0 {
		t.Fatal("no family labels")
	}
}
