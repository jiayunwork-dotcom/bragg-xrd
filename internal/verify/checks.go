package verify

import (
	"fmt"
	"math"

	"bragg-xrd/internal/xrd"
)

type Check struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func CheckEquation(lambda, d float64, n int) Check {
	result, err := xrd.BraggAngle(lambda, d, n)
	if err != nil {
		return Check{Name: "equation", OK: false, Message: err.Error()}
	}
	if !result.Possible {
		return Check{Name: "equation", OK: false, Message: "not possible"}
	}
	residual := xrd.CheckEquation(result, 0)
	ok := residual <= 1e-9
	return Check{
		Name:    "equation",
		OK:      ok,
		Message: fmt.Sprintf("residual=%g 2theta=%.6f", residual, result.TwoTheta),
	}
}

func CheckSinThetaGate(lambda, d float64) Check {
	result, err := xrd.BraggAngle(lambda, d, 5)
	if err != nil {
		return Check{Name: "sin-gate", OK: false, Message: err.Error()}
	}
	ok := !result.Possible && result.TwoTheta == 0
	return Check{
		Name:    "sin-gate",
		OK:      ok,
		Message: fmt.Sprintf("sin=%.4f possible=%v", result.SinTheta, result.Possible),
	}
}

func CheckLambdaOrder(lambda1, lambda2, d float64) Check {
	ok := xrd.LambdaLargerTwoThetaHigher(lambda1, lambda2, d)
	return Check{
		Name:    "lambda-order",
		OK:      ok,
		Message: fmt.Sprintf("lambda %.4f vs %.4f", lambda1, lambda2),
	}
}

func CheckAOrder(lambda, a1, a2 float64, hkl xrd.HKL) Check {
	d1, err := xrd.LatticeSpacing(a1, hkl)
	if err != nil {
		return Check{Name: "a-order", OK: false, Message: err.Error()}
	}
	d2, err := xrd.LatticeSpacing(a2, hkl)
	if err != nil {
		return Check{Name: "a-order", OK: false, Message: err.Error()}
	}
	t1, _ := xrd.TwoTheta(lambda, d1, 1)
	t2, _ := xrd.TwoTheta(lambda, d2, 1)
	ok := t2 < t1
	return Check{
		Name:    "a-order",
		OK:      ok,
		Message: fmt.Sprintf("a %.4f -> %.4f, a larger gives 2theta %.4f < %.4f", a1, a2, t2, t1),
	}
}

func CheckFCCForbidden(lambda, a float64) Check {
	result, err := xrd.Powder(lambda, a, "fcc", 4)
	if err != nil {
		return Check{Name: "fcc-forbidden", OK: false, Message: err.Error()}
	}
	ok := !xrd.ContainsForbidden(result, "fcc")
	return Check{
		Name:    "fcc-forbidden",
		OK:      ok,
		Message: fmt.Sprintf("peaks=%d", result.Count),
	}
}

func CheckBCCForbidden(lambda, a float64) Check {
	result, err := xrd.Powder(lambda, a, "bcc", 4)
	if err != nil {
		return Check{Name: "bcc-forbidden", OK: false, Message: err.Error()}
	}
	ok := !xrd.ContainsForbidden(result, "bcc")
	return Check{
		Name:    "bcc-forbidden",
		OK:      ok,
		Message: fmt.Sprintf("peaks=%d", result.Count),
	}
}

func CheckFCC111Above40() Check {
	twoTheta, err := xrd.ExampleTwoTheta()
	if err != nil {
		return Check{Name: "fcc-111-40", OK: false, Message: err.Error()}
	}
	ok := twoTheta > 40
	return Check{
		Name:    "fcc-111-40",
		OK:      ok,
		Message: fmt.Sprintf("2theta=%.4f", twoTheta),
	}
}

func CheckRightAngle(lambda, d float64, n int) Check {
	result, err := xrd.BraggAngle(lambda, d, n)
	if err != nil {
		return Check{Name: "right-angle", OK: false, Message: err.Error()}
	}
	ok := result.Possible && math.Abs(result.Theta-90) < 1e-9
	return Check{
		Name:    "right-angle",
		OK:      ok,
		Message: fmt.Sprintf("theta=%.6f", result.Theta),
	}
}

func CheckOrderList(lambda, d float64) Check {
	orders, err := xrd.AllOrders(lambda, d)
	if err != nil {
		return Check{Name: "order-list", OK: false, Message: err.Error()}
	}
	ok := len(orders) >= 2
	return Check{
		Name:    "order-list",
		OK:      ok,
		Message: fmt.Sprintf("orders=%d", len(orders)),
	}
}

func RunAll() []Check {
	return []Check{
		CheckEquation(1.5406, 2.087, 1),
		CheckSinThetaGate(1.5406, 1),
		CheckLambdaOrder(1.2, 1.54, 2.5),
		CheckAOrder(1.5406, 3.0, 4.0, xrd.HKL{H: 1, K: 1, L: 1}),
		CheckFCCForbidden(1.5406, 3.615),
		CheckBCCForbidden(1.5406, 2.866),
		CheckFCC111Above40(),
		CheckRightAngle(3.0, 1.5, 1),
		CheckOrderList(1.5406, 2.5),
	}
}

func AllPass(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func FormatChecks(checks []Check) string {
	out := ""
	for _, check := range checks {
		state := "PASS"
		if !check.OK {
			state = "FAIL"
		}
		out += fmt.Sprintf("%-18s %s %s\n", check.Name, state, check.Message)
	}
	return out
}

func CheckValidation(lambda, d float64, n int) Check {
	_, err := xrd.BraggAngle(lambda, d, n)
	ok := err != nil
	return Check{Name: "validation", OK: ok, Message: fmt.Sprintf("err=%v", err)}
}

func CheckSpacingFormula(a float64, hkl xrd.HKL, expected float64) Check {
	d, err := xrd.LatticeSpacing(a, hkl)
	if err != nil {
		return Check{Name: "spacing", OK: false, Message: err.Error()}
	}
	ok := math.Abs(d-expected) < 1e-9
	return Check{Name: "spacing", OK: ok, Message: fmt.Sprintf("d=%.9f expected=%.9f", d, expected)}
}

func CheckPeakCount(result xrd.PowderResult, min int) Check {
	ok := result.Count >= min
	return Check{Name: "peak-count", OK: ok, Message: fmt.Sprintf("count=%d", result.Count)}
}

func CheckHas111(result xrd.PowderResult) Check {
	ok := xrd.HasHKL(result, xrd.HKL{H: 1, K: 1, L: 1})
	return Check{Name: "has-111", OK: ok, Message: fmt.Sprintf("peaks=%d", result.Count)}
}

func CheckNoForbiddenHKL(result xrd.PowderResult, lattice string) Check {
	ok := !xrd.ContainsForbidden(result, lattice)
	return Check{Name: "no-forbidden", OK: ok, Message: fmt.Sprintf("lattice=%s", lattice)}
}

func CheckMonotonic(result xrd.PowderResult) Check {
	ok := xrd.IsMonotonicAngles(result)
	return Check{Name: "monotonic", OK: ok, Message: fmt.Sprintf("peaks=%d", result.Count)}
}

func CheckFirstAngle(result xrd.PowderResult, min float64) Check {
	ok := xrd.MinPeakAngle(result) >= min
	return Check{Name: "first-angle", OK: ok, Message: fmt.Sprintf("first=%.4f", xrd.MinPeakAngle(result))}
}

func CheckLastAngle(result xrd.PowderResult, max float64) Check {
	ok := xrd.MaxPeakAngle(result) <= max
	return Check{Name: "last-angle", OK: ok, Message: fmt.Sprintf("last=%.4f", xrd.MaxPeakAngle(result))}
}

func CheckPositiveAngles(result xrd.PowderResult) Check {
	ok := xrd.MinPeakAngle(result) > 0
	return Check{Name: "positive-angles", OK: ok, Message: fmt.Sprintf("min=%.4f", xrd.MinPeakAngle(result))}
}

func CheckPossibleFlag(result xrd.BraggResult) Check {
	ok := result.Possible
	return Check{Name: "possible", OK: ok, Message: fmt.Sprintf("possible=%v", result.Possible)}
}

func CheckImpossibleFlag(result xrd.BraggResult) Check {
	ok := !result.Possible
	return Check{Name: "impossible", OK: ok, Message: fmt.Sprintf("possible=%v", result.Possible)}
}

func CheckAngleFinite(result xrd.BraggResult) Check {
	ok := result.IsFinite()
	return Check{Name: "finite", OK: ok, Message: fmt.Sprintf("%+v", result)}
}

func CheckTwoThetaRange(result xrd.BraggResult) Check {
	ok := result.TwoTheta >= 0 && result.TwoTheta <= 180
	return Check{Name: "range", OK: ok, Message: fmt.Sprintf("2theta=%.4f", result.TwoTheta)}
}

func CheckThetaOrder(results []xrd.BraggResult) Check {
	ok := true
	for i := 1; i < len(results); i++ {
		if results[i].Theta <= results[i-1].Theta {
			ok = false
		}
	}
	return Check{Name: "theta-order", OK: ok, Message: fmt.Sprintf("orders=%d", len(results))}
}

func CheckAllowedOrder(n int) Check {
	ok := xrd.IsOrderAllowed(n)
	return Check{Name: "allowed-order", OK: ok, Message: fmt.Sprintf("n=%d", n)}
}

func CheckHKLValid(hkl xrd.HKL) Check {
	ok := xrd.IsValidMiller(hkl)
	return Check{Name: "hkl-valid", OK: ok, Message: hkl.String()}
}

func CheckHKLInvalid(hkl xrd.HKL) Check {
	ok := !xrd.IsValidMiller(hkl)
	return Check{Name: "hkl-invalid", OK: ok, Message: hkl.String()}
}

func CheckLatticeValid(lattice string) Check {
	ok := xrd.IsKnownLattice(lattice)
	return Check{Name: "lattice-valid", OK: ok, Message: lattice}
}

func CheckLatticeInvalid(lattice string) Check {
	ok := !xrd.IsKnownLattice(lattice)
	return Check{Name: "lattice-invalid", OK: ok, Message: lattice}
}

func CheckStructureFactor(lattice string, hkl xrd.HKL, expected int) Check {
	factor, err := xrd.StructureFactor(lattice, hkl)
	if err != nil {
		return Check{Name: "structure-factor", OK: false, Message: err.Error()}
	}
	ok := factor == expected
	return Check{Name: "structure-factor", OK: ok, Message: fmt.Sprintf("F=%d expected=%d", factor, expected)}
}

func CheckObservable(lattice string, hkl xrd.HKL) Check {
	ok, _ := xrd.IsObservable(lattice, hkl)
	return Check{Name: "observable", OK: ok, Message: hkl.String()}
}

func CheckForbidden(lattice string, hkl xrd.HKL) Check {
	ok, _ := xrd.IsForbidden(lattice, hkl)
	return Check{Name: "forbidden", OK: ok, Message: hkl.String()}
}

func CheckAllowed(lattice string, hkl xrd.HKL) Check {
	ok, _ := xrd.Allowed(lattice, hkl)
	return Check{Name: "allowed", OK: ok, Message: hkl.String()}
}

func CheckIntensityPositive(lattice string, hkl xrd.HKL) Check {
	intensity, err := xrd.Intensity(lattice, hkl)
	if err != nil {
		return Check{Name: "intensity", OK: false, Message: err.Error()}
	}
	ok := intensity > 0
	return Check{Name: "intensity", OK: ok, Message: fmt.Sprintf("I=%.1f", intensity)}
}

func CheckIntensityZero(lattice string, hkl xrd.HKL) Check {
	intensity, err := xrd.Intensity(lattice, hkl)
	if err != nil {
		return Check{Name: "intensity-zero", OK: false, Message: err.Error()}
	}
	ok := intensity == 0
	return Check{Name: "intensity-zero", OK: ok, Message: fmt.Sprintf("I=%.1f", intensity)}
}

func CheckFamilyLabel(hkl xrd.HKL, expected string) Check {
	ok := xrd.FamilyLabel(hkl) == expected
	return Check{Name: "family-label", OK: ok, Message: fmt.Sprintf("%s vs %s", xrd.FamilyLabel(hkl), expected)}
}

func CheckReducedHKL(hkl, expected xrd.HKL) Check {
	ok := xrd.ReducedHKL(hkl) == expected
	return Check{Name: "reduced-hkl", OK: ok, Message: fmt.Sprintf("%s vs %s", xrd.ReducedHKL(hkl), expected)}
}

func CheckMultiplicity(hkl xrd.HKL, expected int) Check {
	ok := xrd.Multiplicity(hkl) == expected
	return Check{Name: "multiplicity", OK: ok, Message: fmt.Sprintf("m=%d expected=%d", xrd.Multiplicity(hkl), expected)}
}

func CheckNormalizedHKL(hkl, expected xrd.HKL) Check {
	ok := xrd.NormalizeHKL(hkl) == expected
	return Check{Name: "normalized-hkl", OK: ok, Message: fmt.Sprintf("%s vs %s", xrd.NormalizeHKL(hkl), expected)}
}

func CheckEquivalentFamily(first, second xrd.HKL) Check {
	ok := xrd.EquivalentFamily(first, second)
	return Check{Name: "equivalent-family", OK: ok, Message: fmt.Sprintf("%s vs %s", first, second)}
}

func CheckNotEquivalent(first, second xrd.HKL) Check {
	ok := !xrd.EquivalentFamily(first, second)
	return Check{Name: "not-equivalent", OK: ok, Message: fmt.Sprintf("%s vs %s", first, second)}
}

func CheckSameFamily(first, second xrd.HKL) Check {
	ok := xrd.SameFamily(first, second)
	return Check{Name: "same-family", OK: ok, Message: fmt.Sprintf("%s vs %s", first, second)}
}

func CheckOppositeSign(first, second xrd.HKL) Check {
	ok := xrd.OppositeSign(first, second)
	return Check{Name: "opposite-sign", OK: ok, Message: fmt.Sprintf("%s vs %s", first, second)}
}

func CheckNotOpposite(first, second xrd.HKL) Check {
	ok := !xrd.OppositeSign(first, second)
	return Check{Name: "not-opposite", OK: ok, Message: fmt.Sprintf("%s vs %s", first, second)}
}

func CheckSymmetryEquivalent(first, second xrd.HKL) Check {
	ok := xrd.IsSymmetryEquivalent(first, second)
	return Check{Name: "symmetry-equivalent", OK: ok, Message: fmt.Sprintf("%s vs %s", first, second)}
}

func CheckNotSymmetryEquivalent(first, second xrd.HKL) Check {
	ok := !xrd.IsSymmetryEquivalent(first, second)
	return Check{Name: "not-symmetry-equivalent", OK: ok, Message: fmt.Sprintf("%s vs %s", first, second)}
}

func CheckMillerMagnitude(hkl xrd.HKL, expected int) Check {
	ok := xrd.MillerMagnitude(hkl) == expected
	return Check{Name: "miller-magnitude", OK: ok, Message: fmt.Sprintf("m=%d expected=%d", xrd.MillerMagnitude(hkl), expected)}
}

func CheckLowIndex(hkl xrd.HKL, max int) Check {
	ok := xrd.IsLowIndex(hkl, max)
	return Check{Name: "low-index", OK: ok, Message: hkl.String()}
}

func CheckNotLowIndex(hkl xrd.HKL, max int) Check {
	ok := !xrd.IsLowIndex(hkl, max)
	return Check{Name: "not-low-index", OK: ok, Message: hkl.String()}
}

func CheckMaxOrder(lambda, d float64, expected int) Check {
	max, err := xrd.MaxOrder(lambda, d)
	if err != nil {
		return Check{Name: "max-order", OK: false, Message: err.Error()}
	}
	ok := max == expected
	return Check{Name: "max-order", OK: ok, Message: fmt.Sprintf("max=%d expected=%d", max, expected)}
}

func CheckWavelengthForAngle(theta, d float64, n int, expected float64) Check {
	lambda, err := xrd.WavelengthForAngle(theta, d, n)
	if err != nil {
		return Check{Name: "wavelength", OK: false, Message: err.Error()}
	}
	ok := math.Abs(lambda-expected) < 1e-9
	return Check{Name: "wavelength", OK: ok, Message: fmt.Sprintf("lambda=%.9f expected=%.9f", lambda, expected)}
}

func CheckSpacingForAngle(theta, lambda float64, n int, expected float64) Check {
	d, err := xrd.SpacingForAngle(theta, lambda, n)
	if err != nil {
		return Check{Name: "spacing-angle", OK: false, Message: err.Error()}
	}
	ok := math.Abs(d-expected) < 1e-9
	return Check{Name: "spacing-angle", OK: ok, Message: fmt.Sprintf("d=%.9f expected=%.9f", d, expected)}
}

func CheckIsPossible(lambda, d float64, n int) Check {
	ok, err := xrd.IsPossible(lambda, d, n)
	if err != nil {
		return Check{Name: "is-possible", OK: false, Message: err.Error()}
	}
	return Check{Name: "is-possible", OK: ok, Message: fmt.Sprintf("possible=%v", ok)}
}

func CheckNotPossible(lambda, d float64, n int) Check {
	ok, err := xrd.IsPossible(lambda, d, n)
	if err != nil {
		return Check{Name: "not-possible", OK: false, Message: err.Error()}
	}
	return Check{Name: "not-possible", OK: !ok, Message: fmt.Sprintf("possible=%v", ok)}
}

func CheckSinThetaValue(lambda, d float64, n int, expected float64) Check {
	sin, err := xrd.SinTheta(lambda, d, n)
	if err != nil {
		return Check{Name: "sin-theta", OK: false, Message: err.Error()}
	}
	ok := math.Abs(sin-expected) < 1e-12
	return Check{Name: "sin-theta", OK: ok, Message: fmt.Sprintf("sin=%.12f expected=%.12f", sin, expected)}
}

func CheckAngleForSin(sin, expected float64) Check {
	ok := math.Abs(xrd.AngleForSin(sin)-expected) < 1e-9
	return Check{Name: "angle-for-sin", OK: ok, Message: fmt.Sprintf("%.6f vs %.6f", xrd.AngleForSin(sin), expected)}
}

func CheckTwoThetaForSin(sin, expected float64) Check {
	ok := math.Abs(xrd.TwoThetaForSin(sin)-expected) < 1e-9
	return Check{Name: "two-theta-for-sin", OK: ok, Message: fmt.Sprintf("%.6f vs %.6f", xrd.TwoThetaForSin(sin), expected)}
}

func CheckOrderCount(lambda, d float64, expected int) Check {
	count, err := xrd.OrderCount(lambda, d)
	if err != nil {
		return Check{Name: "order-count", OK: false, Message: err.Error()}
	}
	ok := count == expected
	return Check{Name: "order-count", OK: ok, Message: fmt.Sprintf("count=%d expected=%d", count, expected)}
}

func CheckHighestOrder(lambda, d float64, expected int) Check {
	highest, err := xrd.HighestOrder(lambda, d)
	if err != nil {
		return Check{Name: "highest-order", OK: false, Message: err.Error()}
	}
	ok := highest == expected
	return Check{Name: "highest-order", OK: ok, Message: fmt.Sprintf("highest=%d expected=%d", highest, expected)}
}

func CheckOrderRange(lambda, d float64, expected [2]int) Check {
	min, max, err := xrd.OrderRange(lambda, d)
	if err != nil {
		return Check{Name: "order-range", OK: false, Message: err.Error()}
	}
	ok := min == expected[0] && max == expected[1]
	return Check{Name: "order-range", OK: ok, Message: fmt.Sprintf("[%d,%d] expected %v", min, max, expected)}
}

func CheckTwoThetaList(lambda, d float64, expected int) Check {
	angles, err := xrd.TwoThetaList(lambda, d)
	if err != nil {
		return Check{Name: "two-theta-list", OK: false, Message: err.Error()}
	}
	ok := len(angles) == expected
	return Check{Name: "two-theta-list", OK: ok, Message: fmt.Sprintf("count=%d expected=%d", len(angles), expected)}
}

func CheckMaxTwoTheta(lambda, d, expected float64) Check {
	value, err := xrd.MaxTwoTheta(lambda, d)
	if err != nil {
		return Check{Name: "max-two-theta", OK: false, Message: err.Error()}
	}
	ok := math.Abs(value-expected) < 1e-9
	return Check{Name: "max-two-theta", OK: ok, Message: fmt.Sprintf("%.6f vs %.6f", value, expected)}
}

func CheckMinTwoTheta(lambda, d, expected float64) Check {
	value, err := xrd.MinTwoTheta(lambda, d)
	if err != nil {
		return Check{Name: "min-two-theta", OK: false, Message: err.Error()}
	}
	ok := math.Abs(value-expected) < 1e-9
	return Check{Name: "min-two-theta", OK: ok, Message: fmt.Sprintf("%.6f vs %.6f", value, expected)}
}

func CheckRightAngleResult(lambda, d float64, n int) Check {
	result := xrd.RightAngleResult(lambda, d, n)
	ok := result.Theta == 90 && result.TwoTheta == 180
	return Check{Name: "right-angle-result", OK: ok, Message: fmt.Sprintf("%+v", result)}
}

func CheckFormatBragg(result xrd.BraggResult) Check {
	ok := xrd.FormatBragg(result) != ""
	return Check{Name: "format-bragg", OK: ok, Message: "nonempty"}
}

func CheckDisplayPeak(peak xrd.Peak) Check {
	ok := xrd.DisplayPeak(peak) != ""
	return Check{Name: "display-peak", OK: ok, Message: "nonempty"}
}

func CheckDisplayBragg(result xrd.BraggResult) Check {
	ok := xrd.DisplayBragg(result) != ""
	return Check{Name: "display-bragg", OK: ok, Message: "nonempty"}
}

func CheckFormatPowder(result xrd.PowderResult) Check {
	ok := xrd.FormatPowder(result) != ""
	return Check{Name: "format-powder", OK: ok, Message: "nonempty"}
}

func CheckPowderSummary(result xrd.PowderResult) Check {
	ok := xrd.Summary(result) != ""
	return Check{Name: "powder-summary", OK: ok, Message: "nonempty"}
}

func CheckPeaksCSV(result xrd.PowderResult) Check {
	ok := xrd.PeaksCSV(result) != ""
	return Check{Name: "peaks-csv", OK: ok, Message: "nonempty"}
}

func CheckHasPeaks(result xrd.PowderResult) Check {
	ok := xrd.HasPeaks(result)
	return Check{Name: "has-peaks", OK: ok, Message: fmt.Sprintf("count=%d", result.Count)}
}

func CheckEmptyResult(result xrd.PowderResult) Check {
	ok := xrd.IsEmptyResult(result)
	return Check{Name: "empty-result", OK: ok, Message: fmt.Sprintf("count=%d", result.Count)}
}

func CheckPeakCountPositive(result xrd.PowderResult) Check {
	ok := result.Count > 0
	return Check{Name: "peak-count-positive", OK: ok, Message: fmt.Sprintf("count=%d", result.Count)}
}

func CheckFirstPeakExists(result xrd.PowderResult) Check {
	_, err := xrd.FirstPeak(result)
	ok := err == nil
	return Check{Name: "first-peak", OK: ok, Message: fmt.Sprintf("err=%v", err)}
}

func CheckLastPeakExists(result xrd.PowderResult) Check {
	_, err := xrd.LastPeak(result)
	ok := err == nil
	return Check{Name: "last-peak", OK: ok, Message: fmt.Sprintf("err=%v", err)}
}

func CheckClosestPeak(result xrd.PowderResult, target float64) Check {
	ok := xrd.ClosestAngle(result, target) > 0
	return Check{Name: "closest-peak", OK: ok, Message: fmt.Sprintf("angle=%.4f", xrd.ClosestAngle(result, target))}
}

func CheckMedianAngle(result xrd.PowderResult) Check {
	ok := xrd.MedianAngle(result) > 0
	return Check{Name: "median-angle", OK: ok, Message: fmt.Sprintf("median=%.4f", xrd.MedianAngle(result))}
}

func CheckMeanAngle(result xrd.PowderResult) Check {
	ok := xrd.MeanAngle(result) > 0
	return Check{Name: "mean-angle", OK: ok, Message: fmt.Sprintf("mean=%.4f", xrd.MeanAngle(result))}
}

func CheckAngleRange(result xrd.PowderResult) Check {
	ok := xrd.AngleRange(result) >= 0
	return Check{Name: "angle-range", OK: ok, Message: fmt.Sprintf("range=%.4f", xrd.AngleRange(result))}
}

func CheckAngleStddev(result xrd.PowderResult) Check {
	ok := xrd.AngleStddev(result) >= 0
	return Check{Name: "angle-stddev", OK: ok, Message: fmt.Sprintf("std=%.4f", xrd.AngleStddev(result))}
}

func CheckPeakDensity(result xrd.PowderResult) Check {
	ok := xrd.PeakDensity(result) >= 0
	return Check{Name: "peak-density", OK: ok, Message: fmt.Sprintf("density=%.4f", xrd.PeakDensity(result))}
}

func CheckRelativeIntensities(result xrd.PowderResult) Check {
	ok := len(xrd.RelativeIntensities(result)) == len(result.Peaks)
	return Check{Name: "relative-intensities", OK: ok, Message: fmt.Sprintf("n=%d", len(xrd.RelativeIntensities(result)))}
}

func CheckIntensitySum(result xrd.PowderResult) Check {
	ok := xrd.IntensitySum(result) > 0
	return Check{Name: "intensity-sum", OK: ok, Message: fmt.Sprintf("sum=%.1f", xrd.IntensitySum(result))}
}

func CheckPeakFamilies(result xrd.PowderResult) Check {
	ok := len(xrd.PeakFamilies(result)) > 0
	return Check{Name: "peak-families", OK: ok, Message: fmt.Sprintf("families=%d", len(xrd.PeakFamilies(result)))}
}

func CheckUniqueAngles(result xrd.PowderResult) Check {
	ok := xrd.CountUniqueAngles(result) > 0
	return Check{Name: "unique-angles", OK: ok, Message: fmt.Sprintf("unique=%d", xrd.CountUniqueAngles(result))}
}

func CheckOrderedPeaks(result xrd.PowderResult) Check {
	ok := xrd.IsOrdered(result)
	return Check{Name: "ordered-peaks", OK: ok, Message: fmt.Sprintf("count=%d", result.Count)}
}

func CheckHighestIntensity(result xrd.PowderResult) Check {
	ok := xrd.MaxIntensity(result) > 0
	return Check{Name: "highest-intensity", OK: ok, Message: fmt.Sprintf("max=%.1f", xrd.MaxIntensity(result))}
}

func CheckLowestAnglePeak(result xrd.PowderResult) Check {
	peak := xrd.LowestAnglePeak(result)
	ok := peak.TwoTheta > 0
	return Check{Name: "lowest-angle", OK: ok, Message: fmt.Sprintf("angle=%.4f", peak.TwoTheta)}
}

func CheckHighestAnglePeak(result xrd.PowderResult) Check {
	peak := xrd.HighestAnglePeak(result)
	ok := peak.TwoTheta > 0
	return Check{Name: "highest-angle", OK: ok, Message: fmt.Sprintf("angle=%.4f", peak.TwoTheta)}
}

func CheckAngleSpread(result xrd.PowderResult) Check {
	ok := xrd.AngleSpread(result) >= 0
	return Check{Name: "angle-spread", OK: ok, Message: fmt.Sprintf("spread=%.4f", xrd.AngleSpread(result))}
}

func CheckRelativeIntensityText(result xrd.PowderResult) Check {
	ok := xrd.RelativeIntensityText(result) != ""
	return Check{Name: "relative-text", OK: ok, Message: "nonempty"}
}

func CheckPowderText(result xrd.PowderResult) Check {
	ok := xrd.PowderText(result) != ""
	return Check{Name: "powder-text", OK: ok, Message: "nonempty"}
}

func CheckPeakAnglesText(result xrd.PowderResult) Check {
	ok := xrd.PeakAnglesText(result) != ""
	return Check{Name: "peak-angles-text", OK: ok, Message: "nonempty"}
}

func CheckHKLList(result xrd.PowderResult) Check {
	ok := xrd.DisplayHKLList(result) != ""
	return Check{Name: "hkl-list", OK: ok, Message: "nonempty"}
}

func CheckIntensityList(result xrd.PowderResult) Check {
	ok := xrd.DisplayIntensityList(result) != ""
	return Check{Name: "intensity-list", OK: ok, Message: "nonempty"}
}

func CheckOrderCountResult(result xrd.PowderResult) Check {
	ok := xrd.OrderCountPeaks(result) > 0
	return Check{Name: "order-count-result", OK: ok, Message: fmt.Sprintf("count=%d", xrd.OrderCountPeaks(result))}
}

func CheckAngleListSum(result xrd.PowderResult) Check {
	ok := xrd.AngleListSum(result) > 0
	return Check{Name: "angle-list-sum", OK: ok, Message: fmt.Sprintf("sum=%.2f", xrd.AngleListSum(result))}
}

func CheckAngleListProduct(result xrd.PowderResult) Check {
	ok := xrd.AngleListProduct(result) > 0
	return Check{Name: "angle-list-product", OK: ok, Message: fmt.Sprintf("product=%.2f", xrd.AngleListProduct(result))}
}

func CheckPeakDensityPositive(result xrd.PowderResult) Check {
	ok := xrd.PeakDensity(result) > 0
	return Check{Name: "peak-density-positive", OK: ok, Message: fmt.Sprintf("density=%.4f", xrd.PeakDensity(result))}
}

func CheckAngleStddevPositive(result xrd.PowderResult) Check {
	ok := xrd.AngleStddev(result) >= 0
	return Check{Name: "stddev-positive", OK: ok, Message: fmt.Sprintf("std=%.4f", xrd.AngleStddev(result))}
}

func CheckAllDisplayHelpers(result xrd.PowderResult) Check {
	ok := xrd.DisplaySummary(result) != "" && xrd.DisplayPeakCount(result) != "" &&
		xrd.DisplayFirstAngle(result) != "" && xrd.DisplayLastAngle(result) != ""
	return Check{Name: "display-helpers", OK: ok, Message: "nonempty"}
}

func CheckUnits() Check {
	ok := xrd.AngstromsToMeters(1) == 1e-10 && xrd.MetersToAngstroms(1e-10) == 1
	return Check{Name: "units", OK: ok, Message: "conversion"}
}

func CheckScenarioExample() Check {
	ok, err := xrd.IsFCCExampleAbove40()
	if err != nil {
		return Check{Name: "scenario-example", OK: false, Message: err.Error()}
	}
	return Check{Name: "scenario-example", OK: ok, Message: "above 40"}
}
