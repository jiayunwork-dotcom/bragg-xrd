package xrd

import "fmt"

func AngstromsToMeters(value float64) float64 {
	return value * 1e-10
}

func MetersToAngstroms(value float64) float64 {
	return value * 1e10
}

func NanometersToAngstroms(value float64) float64 {
	return value * 10
}

func AngstromsToNanometers(value float64) float64 {
	return value / 10
}

func FormatAngleDegrees(value float64) string {
	return fmt.Sprintf("%.4f deg", value)
}

func FormatSpacing(value float64) string {
	return fmt.Sprintf("%.4f A", value)
}

func FormatWavelength(value float64) string {
	return fmt.Sprintf("%.4f A", value)
}

func FormatCell(value float64) string {
	return fmt.Sprintf("%.4f A", value)
}

func RadiansToDegrees(value float64) float64 {
	return value * 180 / 3.141592653589793
}

func DegreesToRadians(value float64) float64 {
	return value * 3.141592653589793 / 180
}

func TwoThetaFromTheta(theta float64) float64 {
	return 2 * theta
}

func ThetaFromTwoTheta(twoTheta float64) float64 {
	return twoTheta / 2
}

func Round4(value float64) float64 {
	return float64(int64(value*10000+0.5)) / 10000
}

func DisplayPeak(peak Peak) string {
	return fmt.Sprintf("%s %.4f deg", peak.HKL, peak.TwoTheta)
}

func DisplayBragg(result BraggResult) string {
	return fmt.Sprintf("n=%d 2theta=%.4f deg possible=%v", result.N, result.TwoTheta, result.Possible)
}

func SinThetaValueText(result BraggResult) string {
	return fmt.Sprintf("%.6f", result.SinTheta)
}

func PossibleText(result BraggResult) string {
	if result.Possible {
		return "possible"
	}
	return "impossible"
}

func IntensityPercent(value, max float64) float64 {
	if max == 0 {
		return 0
	}
	return value / max * 100
}

func FormatIntensity(value float64) string {
	return fmt.Sprintf("%.1f", value)
}

func LambdaInAngstroms(lambdaMeters float64) float64 {
	return MetersToAngstroms(lambdaMeters)
}

func DInAngstroms(dMeters float64) float64 {
	return MetersToAngstroms(dMeters)
}

func AInAngstroms(aMeters float64) float64 {
	return MetersToAngstroms(aMeters)
}

func TwoThetaListText(angles []float64) string {
	out := ""
	for _, angle := range angles {
		out += fmt.Sprintf("%.4f\n", angle)
	}
	return out
}

func OrderText(n int) string {
	return fmt.Sprintf("%d", n)
}

func FamilyText(hkl HKL) string {
	return FamilyLabel(hkl)
}

func UnitLabel() string {
	return "angstrom"
}

func ConvertScenarioToSI(scenario Scenario) Scenario {
	scenario.Lambda = AngstromsToMeters(scenario.Lambda)
	scenario.A = AngstromsToMeters(scenario.A)
	return scenario
}

func ConvertScenarioToAngstrom(scenario Scenario) Scenario {
	scenario.Lambda = MetersToAngstroms(scenario.Lambda)
	scenario.A = MetersToAngstroms(scenario.A)
	return scenario
}

func DisplayAngle(angle float64) string {
	return FormatAngleDegrees(angle)
}

func DisplayTwoTheta(result BraggResult) string {
	return FormatAngleDegrees(result.TwoTheta)
}

func DisplayTheta(result BraggResult) string {
	return FormatAngleDegrees(result.Theta)
}

func DisplayD(result BraggResult) string {
	return FormatSpacing(result.D)
}

func DisplayLambda(result BraggResult) string {
	return FormatWavelength(result.Lambda)
}

func DisplayPossible(result BraggResult) string {
	return PossibleText(result)
}

func DisplayOrder(result BraggResult) string {
	return OrderText(result.N)
}

func DisplayPeaks(result PowderResult) string {
	return FormatPowder(result)
}

func DisplaySummary(result PowderResult) string {
	return Summary(result)
}

func DisplayPeakCount(result PowderResult) string {
	return fmt.Sprintf("%d", result.Count)
}

func DisplayFirstAngle(result PowderResult) string {
	return FormatAngleDegrees(MinPeakAngle(result))
}

func DisplayLastAngle(result PowderResult) string {
	return FormatAngleDegrees(MaxPeakAngle(result))
}

func DisplayLattice(lattice string) string {
	return LatticeName(lattice)
}

func DisplayForbidden(lattice string) string {
	return ForbiddenText(lattice)
}

func DisplayCell(a float64) string {
	return FormatCell(a)
}

func DisplayFamily(hkl HKL) string {
	return FamilyLabel(hkl)
}

func DisplayMiller(hkl HKL) string {
	return hkl.String()
}

func DisplayMultiplicity(hkl HKL) string {
	return fmt.Sprintf("%d", Multiplicity(hkl))
}

func DisplayStructureFactor(lattice string, hkl HKL) string {
	factor, err := StructureFactor(lattice, hkl)
	if err != nil {
		return "error"
	}
	return fmt.Sprintf("%d", factor)
}

func DisplayObservable(lattice string, hkl HKL) string {
	ok, _ := IsObservable(lattice, hkl)
	if ok {
		return "observable"
	}
	return "forbidden"
}

func DisplayIntensity(lattice string, hkl HKL) string {
	intensity, _ := Intensity(lattice, hkl)
	return FormatIntensity(intensity)
}

func DisplayDSpacing(a float64, hkl HKL) string {
	d, err := LatticeSpacing(a, hkl)
	if err != nil {
		return "error"
	}
	return FormatSpacing(d)
}

func DisplayAngleRange(result PowderResult) string {
	return fmt.Sprintf("%.4f..%.4f", MinPeakAngle(result), MaxPeakAngle(result))
}

func DisplayFamilyCount(result PowderResult) string {
	return fmt.Sprintf("%d", FamilyCount(result))
}

func DisplayUniqueCount(result PowderResult) string {
	return fmt.Sprintf("%d", CountUniqueAngles(result))
}

func DisplayMaxIntensity(result PowderResult) string {
	return FormatIntensity(MaxIntensity(result))
}

func DisplayMeanAngle(result PowderResult) string {
	return FormatAngleDegrees(MeanAngle(result))
}

func DisplayMedianAngle(result PowderResult) string {
	return FormatAngleDegrees(MedianAngle(result))
}

func DisplayStddevAngle(result PowderResult) string {
	return fmt.Sprintf("%.4f deg", AngleStddev(result))
}

func DisplayPeakDensity(result PowderResult) string {
	return fmt.Sprintf("%.4f", PeakDensity(result))
}

func DisplayAngleSpread(result PowderResult) string {
	return FormatAngleDegrees(AngleSpread(result))
}

func DisplayOrderCount(result PowderResult) string {
	return fmt.Sprintf("%d", OrderCountPeaks(result))
}

func OrderCountPeaks(result PowderResult) int {
	count := 0
	for _, peak := range result.Peaks {
		count += peak.Order
	}
	return count
}

func DisplayHKLList(result PowderResult) string {
	out := ""
	for _, peak := range result.Peaks {
		out += peak.HKL.String() + "\n"
	}
	return out
}

func DisplayIntensityList(result PowderResult) string {
	out := ""
	for _, peak := range result.Peaks {
		out += fmt.Sprintf("%.1f\n", peak.Intensity)
	}
	return out
}

func DisplayRelativeIntensityList(result PowderResult) string {
	out := ""
	for _, value := range RelativeIntensities(result) {
		out += fmt.Sprintf("%.1f%%\n", value)
	}
	return out
}

func DisplayCSV(result PowderResult) string {
	return PeaksCSV(result)
}

func AttachImpossibleAngle(result *BraggResult) {
	if result == nil {
		return
	}
	if result.Possible {
		return
	}
	result.TwoTheta = TwoThetaForSin(result.SinTheta)
	result.Theta = ThetaFromTwoTheta(result.TwoTheta)
}

func ImpossibleAngleFromSin(sin float64) float64 {
	return TwoThetaForSin(sin)
}
