package xrd

import (
	"fmt"
	"math"
	"sort"
)

type PowderResult struct {
	Lambda   float64 `json:"lambda"`
	A        float64 `json:"a"`
	Lattice  string  `json:"lattice"`
	MaxIndex int     `json:"max_index"`
	Peaks    []Peak  `json:"peaks"`
	Count    int     `json:"count"`
}

func Powder(lambda, a float64, lattice string, maxIndex int) (PowderResult, error) {
	if err := ValidateWavelength(lambda); err != nil {
		return PowderResult{}, err
	}
	if err := ValidateCellConstant(a); err != nil {
		return PowderResult{}, err
	}
	if err := ValidateLattice(lattice); err != nil {
		return PowderResult{}, err
	}
	if maxIndex < 1 {
		return PowderResult{}, fmt.Errorf("max miller index must be >= 1")
	}
	peaks := make([]Peak, 0)
	for _, hkl := range MillerSet(maxIndex) {
		allowed, err := Allowed(lattice, hkl)
		if err != nil {
			return PowderResult{}, err
		}
		if !allowed {
			continue
		}
		d, err := LatticeSpacing(a, hkl)
		if err != nil {
			return PowderResult{}, err
		}
		maxOrder, err := MaxOrder(lambda, d)
		if err != nil {
			return PowderResult{}, err
		}
		for n := 1; n <= maxOrder; n++ {
			result, err := BraggAngle(lambda, d, n)
			if err != nil {
				return PowderResult{}, err
			}
			if !result.Possible {
				continue
			}
			intensity, err := Intensity(lattice, hkl)
			if err != nil {
				return PowderResult{}, err
			}
			peaks = append(peaks, Peak{
				HKL: hkl, TwoTheta: result.TwoTheta, Theta: result.Theta,
				Order: n, Intensity: intensity,
			})
		}
	}
	sort.SliceStable(peaks, func(i, j int) bool {
		if math.Abs(peaks[i].TwoTheta-peaks[j].TwoTheta) > 1e-9 {
			return peaks[i].TwoTheta < peaks[j].TwoTheta
		}
		return peaks[i].Order < peaks[j].Order
	})
	return PowderResult{
		Lambda: lambda, A: a, Lattice: lattice, MaxIndex: maxIndex,
		Peaks: peaks, Count: len(peaks),
	}, nil
}

func PowderUnique(lambda, a float64, lattice string, maxIndex int) (PowderResult, error) {
	result, err := Powder(lambda, a, lattice, maxIndex)
	if err != nil {
		return PowderResult{}, err
	}
	unique := make([]Peak, 0, len(result.Peaks))
	for _, peak := range result.Peaks {
		if len(unique) == 0 || math.Abs(unique[len(unique)-1].TwoTheta-peak.TwoTheta) > 1e-7 {
			unique = append(unique, peak)
		}
	}
	result.Peaks = unique
	result.Count = len(unique)
	return result, nil
}

func SingleCrystal(lambda, a float64, lattice string, hkl HKL, maxOrder int) ([]Peak, error) {
	if err := ValidateWavelength(lambda); err != nil {
		return nil, err
	}
	if err := ValidateCellConstant(a); err != nil {
		return nil, err
	}
	if err := ValidateLattice(lattice); err != nil {
		return nil, err
	}
	if err := ValidateHKL(hkl); err != nil {
		return nil, err
	}
	allowed, err := Allowed(lattice, hkl)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return []Peak{}, nil
	}
	d, err := LatticeSpacing(a, hkl)
	if err != nil {
		return nil, err
	}
	limit, err := MaxOrder(lambda, d)
	if err != nil {
		return nil, err
	}
	if maxOrder > 0 && maxOrder < limit {
		limit = maxOrder
	}
	peaks := make([]Peak, 0, limit)
	for n := 1; n <= limit; n++ {
		result, err := BraggAngle(lambda, d, n)
		if err != nil {
			return nil, err
		}
		if result.Possible {
			peaks = append(peaks, Peak{HKL: hkl, TwoTheta: result.TwoTheta, Theta: result.Theta, Order: n})
		}
	}
	return peaks, nil
}

func FirstPeaks(lambda, a float64, lattice string, maxIndex int, count int) (PowderResult, error) {
	result, err := Powder(lambda, a, lattice, maxIndex)
	if err != nil {
		return PowderResult{}, err
	}
	if count > 0 && count < len(result.Peaks) {
		result.Peaks = result.Peaks[:count]
		result.Count = count
	}
	return result, nil
}

func PeakAngles(lambda, a float64, lattice string, maxIndex int) ([]float64, error) {
	result, err := Powder(lambda, a, lattice, maxIndex)
	if err != nil {
		return nil, err
	}
	angles := make([]float64, 0, len(result.Peaks))
	for _, peak := range result.Peaks {
		angles = append(angles, peak.TwoTheta)
	}
	return angles, nil
}

func ContainsForbidden(result PowderResult, lattice string) bool {
	for _, peak := range result.Peaks {
		forbidden, _ := IsForbidden(lattice, peak.HKL)
		if forbidden {
			return true
		}
	}
	return false
}

func HasHKL(result PowderResult, hkl HKL) bool {
	for _, peak := range result.Peaks {
		if peak.HKL == hkl {
			return true
		}
	}
	return false
}

func AngleForHKL(lambda, a float64, lattice string, hkl HKL) (float64, error) {
	allowed, err := Allowed(lattice, hkl)
	if err != nil {
		return 0, err
	}
	if !allowed {
		return 0, fmt.Errorf("forbidden reflection")
	}
	d, err := LatticeSpacing(a, hkl)
	if err != nil {
		return 0, err
	}
	return TwoTheta(lambda, d, 1)
}

func FormatPowder(result PowderResult) string {
	out := fmt.Sprintf("lattice=%s a=%.4f peaks=%d\n", result.Lattice, result.A, result.Count)
	for _, peak := range result.Peaks {
		out += fmt.Sprintf("%s %.4f deg n=%d\n", peak.HKL, peak.TwoTheta, peak.Order)
	}
	return out
}

func PowderText(result PowderResult) string {
	return FormatPowder(result)
}

func PeakCount(result PowderResult) int {
	return result.Count
}

func MaxPeakAngle(result PowderResult) float64 {
	max := 0.0
	for _, peak := range result.Peaks {
		if peak.TwoTheta > max {
			max = peak.TwoTheta
		}
	}
	return max
}

func MinPeakAngle(result PowderResult) float64 {
	if len(result.Peaks) == 0 {
		return 0
	}
	return result.Peaks[0].TwoTheta
}

func PeakAnglesList(result PowderResult) []float64 {
	angles := make([]float64, 0, len(result.Peaks))
	for _, peak := range result.Peaks {
		angles = append(angles, peak.TwoTheta)
	}
	return angles
}

func SortedPeaks(result PowderResult) PowderResult {
	sort.SliceStable(result.Peaks, func(i, j int) bool {
		return result.Peaks[i].TwoTheta < result.Peaks[j].TwoTheta
	})
	return result
}

func FirstN(result PowderResult, n int) PowderResult {
	if n < 0 {
		n = 0
	}
	if n > len(result.Peaks) {
		n = len(result.Peaks)
	}
	result.Peaks = result.Peaks[:n]
	result.Count = n
	return result
}

func PeakFamilies(result PowderResult) []string {
	seen := make(map[string]bool)
	families := make([]string, 0)
	for _, peak := range result.Peaks {
		label := FamilyLabel(peak.HKL)
		if !seen[label] {
			seen[label] = true
			families = append(families, label)
		}
	}
	return families
}

func HasFamily(result PowderResult, hkl HKL) bool {
	label := FamilyLabel(hkl)
	for _, family := range PeakFamilies(result) {
		if family == label {
			return true
		}
	}
	return false
}

func CountUniqueAngles(result PowderResult) int {
	seen := make(map[string]bool)
	for _, peak := range result.Peaks {
		key := fmt.Sprintf("%.6f", peak.TwoTheta)
		seen[key] = true
	}
	return len(seen)
}

func IsMonotonicAngles(result PowderResult) bool {
	for i := 1; i < len(result.Peaks); i++ {
		if result.Peaks[i].TwoTheta < result.Peaks[i-1].TwoTheta {
			return false
		}
	}
	return true
}

func MaxIntensity(result PowderResult) float64 {
	max := 0.0
	for _, peak := range result.Peaks {
		if peak.Intensity > max {
			max = peak.Intensity
		}
	}
	return max
}

func RelativeIntensities(result PowderResult) []float64 {
	max := MaxIntensity(result)
	if max == 0 {
		return nil
	}
	out := make([]float64, len(result.Peaks))
	for i, peak := range result.Peaks {
		out[i] = peak.Intensity / max * 100
	}
	return out
}

func IntensitySum(result PowderResult) float64 {
	total := 0.0
	for _, peak := range result.Peaks {
		total += peak.Intensity
	}
	return total
}

func PeakAt(result PowderResult, index int) (Peak, error) {
	if index < 0 || index >= len(result.Peaks) {
		return Peak{}, fmt.Errorf("peak index out of range")
	}
	return result.Peaks[index], nil
}

func AngleAt(result PowderResult, index int) (float64, error) {
	peak, err := PeakAt(result, index)
	if err != nil {
		return 0, err
	}
	return peak.TwoTheta, nil
}

func HKLAt(result PowderResult, index int) (HKL, error) {
	peak, err := PeakAt(result, index)
	if err != nil {
		return HKL{}, err
	}
	return peak.HKL, nil
}

func IsFirstAngleAbove(result PowderResult, threshold float64) bool {
	if len(result.Peaks) == 0 {
		return false
	}
	return result.Peaks[0].TwoTheta > threshold
}

func IsAngleBelow(result PowderResult, index int, threshold float64) bool {
	if index < 0 || index >= len(result.Peaks) {
		return false
	}
	return result.Peaks[index].TwoTheta < threshold
}

func CountInRange(result PowderResult, min, max float64) int {
	count := 0
	for _, peak := range result.Peaks {
		if peak.TwoTheta >= min && peak.TwoTheta <= max {
			count++
		}
	}
	return count
}

func AnglesInRange(result PowderResult, min, max float64) []float64 {
	angles := make([]float64, 0)
	for _, peak := range result.Peaks {
		if peak.TwoTheta >= min && peak.TwoTheta <= max {
			angles = append(angles, peak.TwoTheta)
		}
	}
	return angles
}

func Summary(result PowderResult) string {
	return fmt.Sprintf("%s a=%.4f peaks=%d first=%.3f last=%.3f",
		result.Lattice, result.A, result.Count,
		MinPeakAngle(result), MaxPeakAngle(result))
}

func FamilyCount(result PowderResult) int {
	return len(PeakFamilies(result))
}

func IsFCCWithoutForbidden(result PowderResult) bool {
	return !ContainsForbidden(result, "fcc")
}

func IsBCCWithoutForbidden(result PowderResult) bool {
	return !ContainsForbidden(result, "bcc")
}

func PeakAnglesText(result PowderResult) string {
	out := ""
	for _, angle := range PeakAnglesList(result) {
		out += fmt.Sprintf("%.4f\n", angle)
	}
	return out
}

func OrderAt(result PowderResult, index int) (int, error) {
	peak, err := PeakAt(result, index)
	if err != nil {
		return 0, err
	}
	return peak.Order, nil
}

func IntensityAt(result PowderResult, index int) (float64, error) {
	peak, err := PeakAt(result, index)
	if err != nil {
		return 0, err
	}
	return peak.Intensity, nil
}

func AllPeaksForHKL(result PowderResult, hkl HKL) []Peak {
	out := make([]Peak, 0)
	for _, peak := range result.Peaks {
		if peak.HKL == hkl {
			out = append(out, peak)
		}
	}
	return out
}

func HasOrder(result PowderResult, hkl HKL, order int) bool {
	for _, peak := range result.Peaks {
		if peak.HKL == hkl && peak.Order == order {
			return true
		}
	}
	return false
}

func HighestOrderPeak(result PowderResult) Peak {
	best := Peak{}
	for _, peak := range result.Peaks {
		if peak.Order > best.Order {
			best = peak
		}
	}
	return best
}

func PeakCountForHKL(result PowderResult, hkl HKL) int {
	return len(AllPeaksForHKL(result, hkl))
}

func EqualAngleCounts(first, second PowderResult) bool {
	return len(first.Peaks) == len(second.Peaks)
}

func DiffCount(first, second PowderResult) int {
	return len(first.Peaks) - len(second.Peaks)
}

func AngleDifference(first, second PowderResult) float64 {
	if len(first.Peaks) == 0 || len(second.Peaks) == 0 {
		return 0
	}
	return first.Peaks[0].TwoTheta - second.Peaks[0].TwoTheta
}

func IsGreaterThan(first, second PowderResult) bool {
	return MaxPeakAngle(first) > MaxPeakAngle(second)
}

func IsLessThan(first, second PowderResult) bool {
	return MaxPeakAngle(first) < MaxPeakAngle(second)
}

func EmptyResult(lambda, a float64, lattice string) PowderResult {
	return PowderResult{Lambda: lambda, A: a, Lattice: lattice}
}

func HasPeaks(result PowderResult) bool {
	return result.Count > 0
}

func IsEmptyResult(result PowderResult) bool {
	return result.Count == 0
}

func LastPeak(result PowderResult) (Peak, error) {
	return PeakAt(result, len(result.Peaks)-1)
}

func FirstPeak(result PowderResult) (Peak, error) {
	return PeakAt(result, 0)
}

func PeakAngleListText(result PowderResult) string {
	return PeakAnglesText(result)
}

func PowderSummary(result PowderResult) string {
	return Summary(result)
}

func AngleAboveThreshold(result PowderResult, threshold float64) int {
	count := 0
	for _, peak := range result.Peaks {
		if peak.TwoTheta > threshold {
			count++
		}
	}
	return count
}

func AngleBelowThreshold(result PowderResult, threshold float64) int {
	count := 0
	for _, peak := range result.Peaks {
		if peak.TwoTheta < threshold {
			count++
		}
	}
	return count
}

func ClosestPeak(result PowderResult, target float64) Peak {
	if len(result.Peaks) == 0 {
		return Peak{}
	}
	best := result.Peaks[0]
	bestDiff := math.Abs(best.TwoTheta - target)
	for _, peak := range result.Peaks {
		diff := math.Abs(peak.TwoTheta - target)
		if diff < bestDiff {
			best = peak
			bestDiff = diff
		}
	}
	return best
}

func ClosestAngle(result PowderResult, target float64) float64 {
	return ClosestPeak(result, target).TwoTheta
}

func MedianAngle(result PowderResult) float64 {
	if len(result.Peaks) == 0 {
		return 0
	}
	return result.Peaks[len(result.Peaks)/2].TwoTheta
}

func MeanAngle(result PowderResult) float64 {
	if len(result.Peaks) == 0 {
		return 0
	}
	total := 0.0
	for _, peak := range result.Peaks {
		total += peak.TwoTheta
	}
	return total / float64(len(result.Peaks))
}

func AngleStddev(result PowderResult) float64 {
	if len(result.Peaks) < 2 {
		return 0
	}
	mean := MeanAngle(result)
	variance := 0.0
	for _, peak := range result.Peaks {
		delta := peak.TwoTheta - mean
		variance += delta * delta
	}
	return math.Sqrt(variance / float64(len(result.Peaks)-1))
}

func AngleRange(result PowderResult) float64 {
	return MaxPeakAngle(result) - MinPeakAngle(result)
}

func PeakDensity(result PowderResult) float64 {
	span := AngleRange(result)
	if span <= 0 {
		return 0
	}
	return float64(result.Count) / span
}

func AngleListSum(result PowderResult) float64 {
	total := 0.0
	for _, peak := range result.Peaks {
		total += peak.TwoTheta
	}
	return total
}

func AngleListProduct(result PowderResult) float64 {
	product := 1.0
	for _, peak := range result.Peaks {
		product *= peak.TwoTheta
	}
	return product
}

func PeakAnglesDegree(result PowderResult) []float64 {
	return PeakAnglesList(result)
}

func AngleCount(result PowderResult) int {
	return result.Count
}

func FamilyLabels(result PowderResult) []string {
	return PeakFamilies(result)
}

func IsOrdered(result PowderResult) bool {
	return IsMonotonicAngles(result)
}

func HighestIntensityPeak(result PowderResult) Peak {
	max := Peak{}
	for _, peak := range result.Peaks {
		if peak.Intensity > max.Intensity {
			max = peak
		}
	}
	return max
}

func LowestAnglePeak(result PowderResult) Peak {
	if len(result.Peaks) == 0 {
		return Peak{}
	}
	return result.Peaks[0]
}

func HighestAnglePeak(result PowderResult) Peak {
	if len(result.Peaks) == 0 {
		return Peak{}
	}
	return result.Peaks[len(result.Peaks)-1]
}

func AngleSpread(result PowderResult) float64 {
	return AngleRange(result)
}

func RelativeIntensityText(result PowderResult) string {
	out := ""
	for i, value := range RelativeIntensities(result) {
		out += fmt.Sprintf("%d %.1f%%\n", i, value)
	}
	return out
}

func PeaksCSV(result PowderResult) string {
	out := "h,k,l,n,two_theta,intensity\n"
	for _, peak := range result.Peaks {
		out += fmt.Sprintf("%d,%d,%d,%d,%.6f,%.1f\n",
			peak.HKL.H, peak.HKL.K, peak.HKL.L, peak.Order,
			peak.TwoTheta, peak.Intensity)
	}
	return out
}

func PowderTextTable(result PowderResult) string {
	return PeaksCSV(result)
}

func PatternSortKey(peak Peak) float64 {
	if peak.TwoTheta == 0 {
		return peak.Intensity
	}
	return peak.Intensity / peak.TwoTheta
}

func PatternSortLess(left, right Peak) bool {
	return PatternSortKey(left) < PatternSortKey(right)
}
