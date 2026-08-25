package xrd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func LatticeSpacing(a float64, hkl HKL) (float64, error) {
	if err := ValidateCellConstant(a); err != nil {
		return 0, err
	}
	if err := ValidateHKL(hkl); err != nil {
		return 0, err
	}
	sum := float64(hkl.H*hkl.H + hkl.K*hkl.K + hkl.L*hkl.L)
	return a / math.Sqrt(sum), nil
}

var lastForbidden HKL
var hasLastForbidden bool

func rememberForbidden(hkl HKL) {
	lastForbidden = hkl
	hasLastForbidden = true
}

func LastForbiddenHKL() (HKL, bool) {
	if !hasLastForbidden {
		return HKL{}, false
	}
	return lastForbidden, true
}

func IsForbidden(lattice string, hkl HKL) (bool, error) {
	if err := ValidateLattice(lattice); err != nil {
		return false, err
	}
	if err := ValidateHKL(hkl); err != nil {
		return false, err
	}
	var forbidden bool
	switch lattice {
	case "primitive", "simple":
		forbidden = false
	case "bcc":
		forbidden = hkl.Sum()%2 != 0
	case "fcc":
		forbidden = !hkl.AllEvenOrAllOdd()
	default:
		return false, fmt.Errorf("unsupported lattice %q", lattice)
	}
	if forbidden {
		rememberForbidden(hkl)
	}
	return forbidden, nil
}

func Allowed(lattice string, hkl HKL) (bool, error) {
	forbidden, err := IsForbidden(lattice, hkl)
	if err != nil {
		return false, err
	}
	return !forbidden, nil
}

func ValidateLattice(lattice string) error {
	switch lattice {
	case "primitive", "simple", "bcc", "fcc":
		return nil
	default:
		return fmt.Errorf("lattice must be primitive/simple, bcc, or fcc, got %q", lattice)
	}
}

func ValidateHKL(hkl HKL) error {
	if hkl.IsZero() {
		return fmt.Errorf("hkl cannot be all zero")
	}
	return nil
}

func ValidateCellConstant(a float64) error {
	if a <= 0 {
		return fmt.Errorf("lattice constant a must be positive, got %g", a)
	}
	if math.IsNaN(a) || math.IsInf(a, 0) {
		return fmt.Errorf("lattice constant must be finite")
	}
	return nil
}

func ValidateWavelength(lambda float64) error {
	if lambda <= 0 {
		return fmt.Errorf("wavelength must be positive, got %g", lambda)
	}
	if math.IsNaN(lambda) || math.IsInf(lambda, 0) {
		return fmt.Errorf("wavelength must be finite")
	}
	return nil
}

func ValidateSpacing(d float64) error {
	if d <= 0 {
		return fmt.Errorf("spacing d must be positive, got %g", d)
	}
	if math.IsNaN(d) || math.IsInf(d, 0) {
		return fmt.Errorf("spacing must be finite")
	}
	return nil
}

func MillersFromText(text string) (HKL, error) {
	cleaned := strings.NewReplacer("(", " ", ")", " ", ",", " ").Replace(text)
	fields := strings.Fields(cleaned)
	if len(fields) == 1 && len(fields[0]) == 3 {
		h, err1 := strconv.Atoi(fields[0][0:1])
		k, err2 := strconv.Atoi(fields[0][1:2])
		l, err3 := strconv.Atoi(fields[0][2:3])
		if err1 == nil && err2 == nil && err3 == nil {
			return HKL{H: h, K: k, L: l}, nil
		}
	}
	if len(fields) != 3 {
		return HKL{}, fmt.Errorf("parse hkl: expected 3 integers, got %q", text)
	}
	h, err1 := strconv.Atoi(fields[0])
	k, err2 := strconv.Atoi(fields[1])
	l, err3 := strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return HKL{}, fmt.Errorf("parse hkl: invalid integer")
	}
	return HKL{H: h, K: k, L: l}, nil
}

func MillerSet(maxIndex int) []HKL {
	set := make([]HKL, 0)
	for h := -maxIndex; h <= maxIndex; h++ {
		for k := -maxIndex; k <= maxIndex; k++ {
			for l := -maxIndex; l <= maxIndex; l++ {
				hkl := HKL{H: h, K: k, L: l}
				if hkl.IsZero() {
					continue
				}
				if h < 0 || (h == 0 && k < 0) || (h == 0 && k == 0 && l < 0) {
					continue
				}
				set = append(set, hkl)
			}
		}
	}
	return set
}

func ReciprocalMagnitudeSquared(hkl HKL) int {
	return hkl.H*hkl.H + hkl.K*hkl.K + hkl.L*hkl.L
}

func Multiplicity(hkl HKL) int {
	absH, absK, absL := hkl.Abs().H, hkl.Abs().K, hkl.Abs().L
	nonZero := 0
	if absH > 0 {
		nonZero++
	}
	if absK > 0 {
		nonZero++
	}
	if absL > 0 {
		nonZero++
	}
	switch nonZero {
	case 1:
		return 6
	case 2:
		if absH == absK && absK == absL {
			return 8
		}
		if absH == absK || absK == absL || absH == absL {
			return 24
		}
		return 24
	default:
		if absH == absK && absK == absL {
			return 8
		}
		return 48
	}
}

func LatticeName(lattice string) string {
	switch lattice {
	case "primitive", "simple":
		return "Primitive"
	case "bcc":
		return "Body-centered cubic"
	case "fcc":
		return "Face-centered cubic"
	default:
		return lattice
	}
}

func DFormulaText() string {
	return "d = a / sqrt(h^2 + k^2 + l^2)"
}

func IsCubic(hkl HKL) bool {
	return true
}

func NormalizeHKL(hkl HKL) HKL {
	abs := hkl.Abs()
	if abs.H == 0 && abs.K == 0 {
		return HKL{H: 0, K: 0, L: 1}
	}
	gcd := GCD(GCD(abs.H, abs.K), abs.L)
	if gcd > 0 {
		return HKL{H: abs.H / gcd, K: abs.K / gcd, L: abs.L / gcd}
	}
	return abs
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func SimilarHKL(first, second HKL) bool {
	return NormalizeHKL(first) == NormalizeHKL(second)
}

func EquivalentFamily(first, second HKL) bool {
	permutations := [][]int{
		{first.H, first.K, first.L},
		{first.H, first.L, first.K},
		{first.K, first.H, first.L},
		{first.K, first.L, first.H},
		{first.L, first.H, first.K},
		{first.L, first.K, first.H},
	}
	for _, p := range permutations {
		if p[0] == second.H && p[1] == second.K && p[2] == second.L {
			return true
		}
	}
	return false
}

func FamilyLabel(hkl HKL) string {
	return fmt.Sprintf("{%d%d%d}", hkl.Abs().H, hkl.Abs().K, hkl.Abs().L)
}

func HKLFromString(text string) (HKL, error) {
	return MillersFromText(text)
}

func IsValidMiller(hkl HKL) bool {
	return ValidateHKL(hkl) == nil
}

func SpacingForFamily(a float64, hkl HKL) (float64, error) {
	return LatticeSpacing(a, hkl)
}

func AllowedFamily(lattice string, hkl HKL) bool {
	allowed, err := Allowed(lattice, hkl)
	return err == nil && allowed
}

func ForbiddenText(lattice string) string {
	switch lattice {
	case "bcc":
		return "h+k+l odd is forbidden"
	case "fcc":
		return "h,k,l must be all even or all odd"
	default:
		return "no systematic absences"
	}
}

func LatticeList() []string {
	return []string{"primitive", "simple", "bcc", "fcc"}
}

func IsKnownLattice(lattice string) bool {
	return ValidateLattice(lattice) == nil
}

func FamilySpacing(a float64, hkl HKL) (float64, error) {
	return LatticeSpacing(a, NormalizeHKL(hkl))
}

func ReducedHKL(hkl HKL) HKL {
	return NormalizeHKL(hkl)
}

func MaxMillerIndex(a, lambda float64, maxIndex int) []HKL {
	out := make([]HKL, 0)
	for _, hkl := range MillerSet(maxIndex) {
		d, err := LatticeSpacing(a, hkl)
		if err != nil {
			continue
		}
		if 2*d >= lambda {
			out = append(out, hkl)
		}
	}
	return out
}

func AllowedSet(lattice string, maxIndex int) ([]HKL, error) {
	allowed := make([]HKL, 0)
	for _, hkl := range MillerSet(maxIndex) {
		ok, err := Allowed(lattice, hkl)
		if err != nil {
			return nil, err
		}
		if ok {
			allowed = append(allowed, hkl)
		}
	}
	return allowed, nil
}

func ForbiddenSet(lattice string, maxIndex int) ([]HKL, error) {
	forbidden := make([]HKL, 0)
	for _, hkl := range MillerSet(maxIndex) {
		ok, err := IsForbidden(lattice, hkl)
		if err != nil {
			return nil, err
		}
		if ok {
			forbidden = append(forbidden, hkl)
		}
	}
	return forbidden, nil
}

func CountAllowed(lattice string, maxIndex int) (int, error) {
	set, err := AllowedSet(lattice, maxIndex)
	if err != nil {
		return 0, err
	}
	return len(set), nil
}

func CountForbidden(lattice string, maxIndex int) (int, error) {
	set, err := ForbiddenSet(lattice, maxIndex)
	if err != nil {
		return 0, err
	}
	return len(set), nil
}

func SameFamily(first, second HKL) bool {
	return EquivalentFamily(first.Abs(), second.Abs())
}

func OppositeSign(first, second HKL) bool {
	return first.H == -second.H && first.K == -second.K && first.L == -second.L
}

func IsSymmetryEquivalent(first, second HKL) bool {
	if first == second {
		return true
	}
	return EquivalentFamily(first.Abs(), second.Abs())
}

func LatticeDescription(lattice string) string {
	return LatticeName(lattice) + ": " + ForbiddenText(lattice)
}

func DSpacingFormula() string {
	return DFormulaText()
}

func CubicSpacing(a, hkl string) string {
	return fmt.Sprintf("d = %s / sqrt(%s)", a, hkl)
}

func StructureFactorText(lattice string) string {
	return LatticeDescription(lattice)
}

func IsStructureFactorAllowed(lattice string, hkl HKL) bool {
	return AllowedFamily(lattice, hkl)
}

func MillerIndicesString(hkl HKL) string {
	return hkl.String()
}

func HKLFromArray(values []int) (HKL, error) {
	if len(values) != 3 {
		return HKL{}, fmt.Errorf("hkl array must have 3 values")
	}
	return HKL{H: values[0], K: values[1], L: values[2]}, nil
}

func HKLToArray(hkl HKL) []int {
	return []int{hkl.H, hkl.K, hkl.L}
}

func SetLabel(hkl HKL) string {
	return FamilyLabel(hkl)
}

func ReducedLabel(hkl HKL) string {
	return FamilyLabel(NormalizeHKL(hkl))
}

func PlaneCount(hkl HKL) int {
	return Multiplicity(hkl)
}

func IsPlanarSpacingPositive(d float64) bool {
	return d > 0
}

func ValidateMillerRange(hkl HKL, maxIndex int) error {
	if hkl.Abs().H > maxIndex || hkl.Abs().K > maxIndex || hkl.Abs().L > maxIndex {
		return fmt.Errorf("miller index exceeds %d", maxIndex)
	}
	return ValidateHKL(hkl)
}

func AllMillerTriples(maxIndex int) [][3]int {
	triples := make([][3]int, 0)
	for h := -maxIndex; h <= maxIndex; h++ {
		for k := -maxIndex; k <= maxIndex; k++ {
			for l := -maxIndex; l <= maxIndex; l++ {
				if h == 0 && k == 0 && l == 0 {
					continue
				}
				triples = append(triples, [3]int{h, k, l})
			}
		}
	}
	return triples
}

func MillerMagnitude(hkl HKL) int {
	return ReciprocalMagnitudeSquared(hkl)
}

func IsLowIndex(hkl HKL, max int) bool {
	return hkl.Abs().H <= max && hkl.Abs().K <= max && hkl.Abs().L <= max
}

func PlaneLabel(hkl HKL) string {
	return hkl.String()
}

func StructureFactorSimple(hkl HKL) int {
	return 1
}

func StructureFactorBCC(hkl HKL) int {
	if hkl.Sum()%2 == 0 {
		return 2
	}
	return 0
}

func StructureFactorFCC(hkl HKL) int {
	if hkl.AllEvenOrAllOdd() {
		return 4
	}
	return 0
}

func StructureFactor(lattice string, hkl HKL) (int, error) {
	switch lattice {
	case "primitive", "simple":
		return StructureFactorSimple(hkl), nil
	case "bcc":
		return StructureFactorBCC(hkl), nil
	case "fcc":
		return StructureFactorFCC(hkl), nil
	default:
		return 0, fmt.Errorf("unsupported lattice")
	}
}

func IsObservable(lattice string, hkl HKL) (bool, error) {
	factor, err := StructureFactor(lattice, hkl)
	if err != nil {
		return false, err
	}
	return factor != 0, nil
}

func ObservableSet(lattice string, maxIndex int) ([]HKL, error) {
	allowed := make([]HKL, 0)
	for _, hkl := range MillerSet(maxIndex) {
		observable, err := IsObservable(lattice, hkl)
		if err != nil {
			return nil, err
		}
		if observable {
			allowed = append(allowed, hkl)
		}
	}
	return allowed, nil
}

func Intensity(lattice string, hkl HKL) (float64, error) {
	factor, err := StructureFactor(lattice, hkl)
	if err != nil {
		return 0, err
	}
	return float64(factor * factor), nil
}

func IsForbiddenFamily(lattice string, hkl HKL) bool {
	forbidden, _ := IsForbidden(lattice, hkl)
	return forbidden
}

func FamilyAllowed(lattice string, hkl HKL) bool {
	return AllowedFamily(lattice, hkl)
}
