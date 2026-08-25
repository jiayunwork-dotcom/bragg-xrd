package patterns

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"bragg-xrd/internal/xrd"
)

type Pattern struct {
	Lambda  float64    `json:"lambda"`
	A       float64    `json:"a"`
	Lattice string     `json:"lattice"`
	Peaks   []xrd.Peak `json:"peaks"`
	Max     float64    `json:"max_intensity"`
}

func Build(lambda, a float64, lattice string, maxIndex int) (Pattern, error) {
	powder, err := xrd.Powder(lambda, a, lattice, maxIndex)
	if err != nil {
		return Pattern{}, err
	}
	return Pattern{
		Lambda: lambda, A: a, Lattice: lattice,
		Peaks: powder.Peaks, Max: xrd.MaxIntensity(powder),
	}, nil
}

func Normalize(pattern Pattern) Pattern {
	if pattern.Max == 0 {
		return pattern
	}
	for i := range pattern.Peaks {
		pattern.Peaks[i].Intensity = pattern.Peaks[i].Intensity / pattern.Max * 100
	}
	pattern.Max = 100
	return pattern
}

func Sort(pattern Pattern) Pattern {
	sort.SliceStable(pattern.Peaks, func(i, j int) bool {
		return pattern.Peaks[i].TwoTheta < pattern.Peaks[j].TwoTheta
	})
	return pattern
}

func Filter(pattern Pattern, min, max float64) Pattern {
	out := make([]xrd.Peak, 0, len(pattern.Peaks))
	for _, peak := range pattern.Peaks {
		if peak.TwoTheta >= min && peak.TwoTheta <= max {
			out = append(out, peak)
		}
	}
	pattern.Peaks = out
	return pattern
}

func FamilyPattern(pattern Pattern) Pattern {
	seen := make(map[string]bool)
	out := make([]xrd.Peak, 0)
	for _, peak := range pattern.Peaks {
		label := xrd.FamilyLabel(peak.HKL)
		if !seen[label] {
			seen[label] = true
			out = append(out, peak)
		}
	}
	pattern.Peaks = out
	return pattern
}

func Strongest(pattern Pattern, count int) Pattern {
	if count >= len(pattern.Peaks) || count <= 0 {
		return pattern
	}
	sorted := append([]xrd.Peak(nil), pattern.Peaks...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Intensity > sorted[j].Intensity
	})
	sorted = sorted[:count]
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].TwoTheta < sorted[j].TwoTheta
	})
	pattern.Peaks = sorted
	return pattern
}

func CSV(pattern Pattern) string {
	var b strings.Builder
	b.WriteString("h,k,l,n,two_theta,intensity\n")
	for _, peak := range pattern.Peaks {
		fmt.Fprintf(&b, "%d,%d,%d,%d,%.6f,%.1f\n",
			peak.HKL.H, peak.HKL.K, peak.HKL.L, peak.Order,
			peak.TwoTheta, peak.Intensity)
	}
	return b.String()
}

func Text(pattern Pattern) string {
	out := fmt.Sprintf("lattice=%s a=%.4f peaks=%d\n", pattern.Lattice, pattern.A, len(pattern.Peaks))
	for _, peak := range pattern.Peaks {
		out += fmt.Sprintf("%-8s %.4f deg I=%.1f\n",
			peak.HKL.String(), peak.TwoTheta, peak.Intensity)
	}
	return out
}

func Summary(pattern Pattern) string {
	return fmt.Sprintf("%d peaks, max 2theta %.4f, max intensity %.1f",
		len(pattern.Peaks), MaxAngle(pattern), pattern.Max)
}

func MaxAngle(pattern Pattern) float64 {
	max := 0.0
	for _, peak := range pattern.Peaks {
		if peak.TwoTheta > max {
			max = peak.TwoTheta
		}
	}
	return max
}

func MinAngle(pattern Pattern) float64 {
	if len(pattern.Peaks) == 0 {
		return 0
	}
	min := pattern.Peaks[0].TwoTheta
	for _, peak := range pattern.Peaks {
		if peak.TwoTheta < min {
			min = peak.TwoTheta
		}
	}
	return min
}

func Count(pattern Pattern) int {
	return len(pattern.Peaks)
}

func HasPeak(pattern Pattern, hkl xrd.HKL) bool {
	for _, peak := range pattern.Peaks {
		if peak.HKL == hkl {
			return true
		}
	}
	return false
}

func ForbiddenCount(pattern Pattern, lattice string) int {
	count := 0
	for _, peak := range pattern.Peaks {
		forbidden, _ := xrd.IsForbidden(lattice, peak.HKL)
		if forbidden {
			count++
		}
	}
	return count
}

func IsClean(pattern Pattern, lattice string) bool {
	return ForbiddenCount(pattern, lattice) == 0
}

func IntensityAt(pattern Pattern, hkl xrd.HKL) float64 {
	for _, peak := range pattern.Peaks {
		if peak.HKL == hkl {
			return peak.Intensity
		}
	}
	return 0
}

func AngleAt(pattern Pattern, hkl xrd.HKL) float64 {
	for _, peak := range pattern.Peaks {
		if peak.HKL == hkl {
			return peak.TwoTheta
		}
	}
	return 0
}

func AddPeak(pattern Pattern, peak xrd.Peak) Pattern {
	pattern.Peaks = append(pattern.Peaks, peak)
	return pattern
}

func RemovePeak(pattern Pattern, index int) Pattern {
	if index < 0 || index >= len(pattern.Peaks) {
		return pattern
	}
	pattern.Peaks = append(pattern.Peaks[:index], pattern.Peaks[index+1:]...)
	return pattern
}

func Merge(first, second Pattern) Pattern {
	first.Peaks = append(first.Peaks, second.Peaks...)
	first.Max = math.Max(first.Max, second.Max)
	return Sort(first)
}

func Equal(first, second Pattern, tolerance float64) bool {
	if len(first.Peaks) != len(second.Peaks) {
		return false
	}
	for i := range first.Peaks {
		if math.Abs(first.Peaks[i].TwoTheta-second.Peaks[i].TwoTheta) > tolerance {
			return false
		}
	}
	return true
}

func Copy(pattern Pattern) Pattern {
	out := Pattern{Lambda: pattern.Lambda, A: pattern.A, Lattice: pattern.Lattice, Max: pattern.Max}
	out.Peaks = append([]xrd.Peak(nil), pattern.Peaks...)
	return out
}

func FirstPeaks(pattern Pattern, count int) Pattern {
	if count <= 0 || count > len(pattern.Peaks) {
		return pattern
	}
	pattern.Peaks = pattern.Peaks[:count]
	return pattern
}

func LastPeaks(pattern Pattern, count int) Pattern {
	if count <= 0 || count > len(pattern.Peaks) {
		return pattern
	}
	pattern.Peaks = pattern.Peaks[len(pattern.Peaks)-count:]
	return pattern
}

func PeakRange(pattern Pattern) float64 {
	return MaxAngle(pattern) - MinAngle(pattern)
}

func MeanAngle(pattern Pattern) float64 {
	if len(pattern.Peaks) == 0 {
		return 0
	}
	sum := 0.0
	for _, peak := range pattern.Peaks {
		sum += peak.TwoTheta
	}
	return sum / float64(len(pattern.Peaks))
}

func StddevAngle(pattern Pattern) float64 {
	if len(pattern.Peaks) < 2 {
		return 0
	}
	mean := MeanAngle(pattern)
	variance := 0.0
	for _, peak := range pattern.Peaks {
		delta := peak.TwoTheta - mean
		variance += delta * delta
	}
	return math.Sqrt(variance / float64(len(pattern.Peaks)-1))
}

func MedianAngle(pattern Pattern) float64 {
	if len(pattern.Peaks) == 0 {
		return 0
	}
	return pattern.Peaks[len(pattern.Peaks)/2].TwoTheta
}

func IntensitySum(pattern Pattern) float64 {
	sum := 0.0
	for _, peak := range pattern.Peaks {
		sum += peak.Intensity
	}
	return sum
}

func RelativeIntensity(pattern Pattern, peak xrd.Peak) float64 {
	if pattern.Max == 0 {
		return 0
	}
	return peak.Intensity / pattern.Max * 100
}

func PeakDensity(pattern Pattern) float64 {
	span := PeakRange(pattern)
	if span <= 0 {
		return 0
	}
	return float64(len(pattern.Peaks)) / span
}

func FamilyCount(pattern Pattern) int {
	seen := make(map[string]bool)
	for _, peak := range pattern.Peaks {
		seen[xrd.FamilyLabel(peak.HKL)] = true
	}
	return len(seen)
}

func UniqueAngleCount(pattern Pattern) int {
	seen := make(map[string]bool)
	for _, peak := range pattern.Peaks {
		seen[fmt.Sprintf("%.6f", peak.TwoTheta)] = true
	}
	return len(seen)
}

func IsMonotonic(pattern Pattern) bool {
	for i := 1; i < len(pattern.Peaks); i++ {
		if pattern.Peaks[i].TwoTheta < pattern.Peaks[i-1].TwoTheta {
			return false
		}
	}
	return true
}

func AngleList(pattern Pattern) []float64 {
	angles := make([]float64, 0, len(pattern.Peaks))
	for _, peak := range pattern.Peaks {
		angles = append(angles, peak.TwoTheta)
	}
	return angles
}

func IntensityList(pattern Pattern) []float64 {
	values := make([]float64, 0, len(pattern.Peaks))
	for _, peak := range pattern.Peaks {
		values = append(values, peak.Intensity)
	}
	return values
}

func HKLList(pattern Pattern) []xrd.HKL {
	hkls := make([]xrd.HKL, 0, len(pattern.Peaks))
	for _, peak := range pattern.Peaks {
		hkls = append(hkls, peak.HKL)
	}
	return hkls
}

func OrderList(pattern Pattern) []int {
	orders := make([]int, 0, len(pattern.Peaks))
	for _, peak := range pattern.Peaks {
		orders = append(orders, peak.Order)
	}
	return orders
}

func Describe(pattern Pattern) string {
	return Summary(pattern)
}

func IsNonEmpty(pattern Pattern) bool {
	return len(pattern.Peaks) > 0
}

func MaxIntensity(pattern Pattern) float64 {
	max := 0.0
	for _, peak := range pattern.Peaks {
		if peak.Intensity > max {
			max = peak.Intensity
		}
	}
	return max
}

func MinIntensity(pattern Pattern) float64 {
	if len(pattern.Peaks) == 0 {
		return 0
	}
	min := pattern.Peaks[0].Intensity
	for _, peak := range pattern.Peaks {
		if peak.Intensity < min {
			min = peak.Intensity
		}
	}
	return min
}

func AnglesInRange(pattern Pattern, min, max float64) []float64 {
	angles := make([]float64, 0)
	for _, peak := range pattern.Peaks {
		if peak.TwoTheta >= min && peak.TwoTheta <= max {
			angles = append(angles, peak.TwoTheta)
		}
	}
	return angles
}

func CountInRange(pattern Pattern, min, max float64) int {
	return len(AnglesInRange(pattern, min, max))
}

func StrongestPeak(pattern Pattern) xrd.Peak {
	best := xrd.Peak{}
	for _, peak := range pattern.Peaks {
		if peak.Intensity > best.Intensity {
			best = peak
		}
	}
	return best
}

func WeakestPeak(pattern Pattern) xrd.Peak {
	if len(pattern.Peaks) == 0 {
		return xrd.Peak{}
	}
	best := pattern.Peaks[0]
	for _, peak := range pattern.Peaks {
		if peak.Intensity < best.Intensity {
			best = peak
		}
	}
	return best
}

func ClosestPeak(pattern Pattern, target float64) xrd.Peak {
	if len(pattern.Peaks) == 0 {
		return xrd.Peak{}
	}
	best := pattern.Peaks[0]
	bestDiff := math.Abs(best.TwoTheta - target)
	for _, peak := range pattern.Peaks {
		diff := math.Abs(peak.TwoTheta - target)
		if diff < bestDiff {
			best = peak
			bestDiff = diff
		}
	}
	return best
}

func ClosestAngle(pattern Pattern, target float64) float64 {
	return ClosestPeak(pattern, target).TwoTheta
}

func HighestOrderPeak(pattern Pattern) xrd.Peak {
	best := xrd.Peak{}
	for _, peak := range pattern.Peaks {
		if peak.Order > best.Order {
			best = peak
		}
	}
	return best
}

func CountByOrder(pattern Pattern, order int) int {
	count := 0
	for _, peak := range pattern.Peaks {
		if peak.Order == order {
			count++
		}
	}
	return count
}

func AngleSpread(pattern Pattern) float64 {
	return PeakRange(pattern)
}

func RelativeIntensityText(pattern Pattern) string {
	out := ""
	for _, peak := range pattern.Peaks {
		out += fmt.Sprintf("%s %.1f%%\n", peak.HKL, RelativeIntensity(pattern, peak))
	}
	return out
}

func PeakAnglesText(pattern Pattern) string {
	out := ""
	for _, angle := range AngleList(pattern) {
		out += fmt.Sprintf("%.4f\n", angle)
	}
	return out
}

func HKLText(pattern Pattern) string {
	out := ""
	for _, hkl := range HKLList(pattern) {
		out += hkl.String() + "\n"
	}
	return out
}

func IntensityText(pattern Pattern) string {
	out := ""
	for _, value := range IntensityList(pattern) {
		out += fmt.Sprintf("%.1f\n", value)
	}
	return out
}

func OrderText(pattern Pattern) string {
	out := ""
	for _, order := range OrderList(pattern) {
		out += fmt.Sprintf("%d\n", order)
	}
	return out
}

func IsEmpty(pattern Pattern) bool {
	return len(pattern.Peaks) == 0
}

func PeakCount(pattern Pattern) int {
	return len(pattern.Peaks)
}

func TotalOrder(pattern Pattern) int {
	total := 0
	for _, peak := range pattern.Peaks {
		total += peak.Order
	}
	return total
}

func AngleSum(pattern Pattern) float64 {
	sum := 0.0
	for _, peak := range pattern.Peaks {
		sum += peak.TwoTheta
	}
	return sum
}

func AngleProduct(pattern Pattern) float64 {
	product := 1.0
	for _, peak := range pattern.Peaks {
		product *= peak.TwoTheta
	}
	return product
}

func IsSorted(pattern Pattern) bool {
	return IsMonotonic(pattern)
}

func HasAngle(pattern Pattern, angle float64, tolerance float64) bool {
	for _, peak := range pattern.Peaks {
		if math.Abs(peak.TwoTheta-angle) <= tolerance {
			return true
		}
	}
	return false
}

func FamilyLabels(pattern Pattern) []string {
	seen := make(map[string]bool)
	labels := make([]string, 0)
	for _, peak := range pattern.Peaks {
		label := xrd.FamilyLabel(peak.HKL)
		if !seen[label] {
			seen[label] = true
			labels = append(labels, label)
		}
	}
	return labels
}

func PatternFromPowder(powder xrd.PowderResult) Pattern {
	return Pattern{
		Lambda: powder.Lambda, A: powder.A, Lattice: powder.Lattice,
		Peaks: powder.Peaks, Max: xrd.MaxIntensity(powder),
	}
}

func ToPowder(pattern Pattern) xrd.PowderResult {
	return xrd.PowderResult{
		Lambda: pattern.Lambda, A: pattern.A, Lattice: pattern.Lattice,
		Peaks: pattern.Peaks, Count: len(pattern.Peaks),
	}
}

func DescribeDiff(first, second Pattern) string {
	return fmt.Sprintf("first=%d peaks second=%d peaks", len(first.Peaks), len(second.Peaks))
}

func IsCleanPattern(pattern Pattern) bool {
	return ForbiddenCount(pattern, pattern.Lattice) == 0
}

func IndexOfPeak(pattern Pattern, hkl xrd.HKL) int {
	for i, peak := range pattern.Peaks {
		if peak.HKL == hkl {
			return i
		}
	}
	return -1
}

func PeakAtIndex(pattern Pattern, index int) (xrd.Peak, error) {
	if index < 0 || index >= len(pattern.Peaks) {
		return xrd.Peak{}, fmt.Errorf("peak index out of range")
	}
	return pattern.Peaks[index], nil
}

func AngleAtIndex(pattern Pattern, index int) (float64, error) {
	peak, err := PeakAtIndex(pattern, index)
	if err != nil {
		return 0, err
	}
	return peak.TwoTheta, nil
}

func HKLAtIndex(pattern Pattern, index int) (xrd.HKL, error) {
	peak, err := PeakAtIndex(pattern, index)
	if err != nil {
		return xrd.HKL{}, err
	}
	return peak.HKL, nil
}

func IntensityAtIndex(pattern Pattern, index int) (float64, error) {
	peak, err := PeakAtIndex(pattern, index)
	if err != nil {
		return 0, err
	}
	return peak.Intensity, nil
}

func OrderAtIndex(pattern Pattern, index int) (int, error) {
	peak, err := PeakAtIndex(pattern, index)
	if err != nil {
		return 0, err
	}
	return peak.Order, nil
}

func FirstAngle(pattern Pattern) float64 {
	return MinAngle(pattern)
}

func LastAngle(pattern Pattern) float64 {
	return MaxAngle(pattern)
}

func AngleRangeText(pattern Pattern) string {
	return fmt.Sprintf("%.4f..%.4f", MinAngle(pattern), MaxAngle(pattern))
}

func CountInRangeText(pattern Pattern, min, max float64) string {
	return fmt.Sprintf("%d", CountInRange(pattern, min, max))
}

func FamilyCountText(pattern Pattern) string {
	return fmt.Sprintf("%d", FamilyCount(pattern))
}

func UniqueCountText(pattern Pattern) string {
	return fmt.Sprintf("%d", UniqueAngleCount(pattern))
}

func MaxIntensityText(pattern Pattern) string {
	return fmt.Sprintf("%.1f", MaxIntensity(pattern))
}

func MeanAngleText(pattern Pattern) string {
	return fmt.Sprintf("%.4f", MeanAngle(pattern))
}

func MedianAngleText(pattern Pattern) string {
	return fmt.Sprintf("%.4f", MedianAngle(pattern))
}

func StddevText(pattern Pattern) string {
	return fmt.Sprintf("%.4f", StddevAngle(pattern))
}

func DensityText(pattern Pattern) string {
	return fmt.Sprintf("%.4f", PeakDensity(pattern))
}

func SpreadText(pattern Pattern) string {
	return fmt.Sprintf("%.4f", AngleSpread(pattern))
}

func TotalOrderText(pattern Pattern) string {
	return fmt.Sprintf("%d", TotalOrder(pattern))
}

func AngleSumText(pattern Pattern) string {
	return fmt.Sprintf("%.2f", AngleSum(pattern))
}

func AngleProductText(pattern Pattern) string {
	return fmt.Sprintf("%.2f", AngleProduct(pattern))
}

func Empty() Pattern {
	return Pattern{}
}

func IsEmptyText(pattern Pattern) string {
	if IsEmpty(pattern) {
		return "empty"
	}
	return fmt.Sprintf("%d peaks", len(pattern.Peaks))
}

func StrongestText(pattern Pattern) string {
	return StrongestPeak(pattern).String()
}

func WeakestText(pattern Pattern) string {
	return WeakestPeak(pattern).String()
}

func ClosestText(pattern Pattern, target float64) string {
	return ClosestPeak(pattern, target).String()
}

func HighestOrderText(pattern Pattern) string {
	return HighestOrderPeak(pattern).String()
}

func CountByOrderText(pattern Pattern, order int) string {
	return fmt.Sprintf("%d", CountByOrder(pattern, order))
}

func HasAngleText(pattern Pattern, angle float64, tolerance float64) string {
	if HasAngle(pattern, angle, tolerance) {
		return "yes"
	}
	return "no"
}

func IsSortedText(pattern Pattern) string {
	if IsSorted(pattern) {
		return "sorted"
	}
	return "unsorted"
}

func IsCleanText(pattern Pattern) string {
	if IsCleanPattern(pattern) {
		return "clean"
	}
	return "has forbidden"
}

func DiffText(first, second Pattern) string {
	return DescribeDiff(first, second)
}

func IndexText(pattern Pattern, hkl xrd.HKL) string {
	return fmt.Sprintf("%d", IndexOfPeak(pattern, hkl))
}

func PeakCountText(pattern Pattern) string {
	return fmt.Sprintf("%d", PeakCount(pattern))
}
