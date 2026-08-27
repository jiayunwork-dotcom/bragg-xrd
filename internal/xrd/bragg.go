package xrd

import (
	"context"
	"fmt"
	"math"
)

func BraggAngle(lambda, d float64, n int) (BraggResult, error) {
	if err := ValidateWavelength(lambda); err != nil {
		return BraggResult{}, err
	}
	if err := ValidateSpacing(d); err != nil {
		return BraggResult{}, err
	}
	if n < 1 {
		return BraggResult{}, fmt.Errorf("order n must be >= 1, got %d", n)
	}
	sinTheta := float64(n) * lambda / (2 * d)
	result := BraggResult{
		Lambda: lambda, D: d, N: n,
		SinTheta: sinTheta, Possible: sinTheta <= 1+1e-12,
	}
	if !result.Possible {
		return result, nil
	}
	if sinTheta > 1 {
		sinTheta = 1
	}
	theta := math.Asin(sinTheta)
	result.Theta = theta * 180 / math.Pi
	result.TwoTheta = 2 * result.Theta
	return result, nil
}

func BraggAngleCtx(ctx context.Context, lambda, d float64, n int) (BraggResult, error) {
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return BraggResult{}, err
		}
		select {
		case <-ctx.Done():
			return BraggResult{}, ctx.Err()
		default:
		}
	}
	return BraggAngle(lambda, d, n)
}

func TwoTheta(lambda, d float64, n int) (float64, error) {
	result, err := BraggAngle(lambda, d, n)
	if err != nil {
		return 0, err
	}
	if !result.Possible {
		return 0, fmt.Errorf("sin(theta) > 1 for n=%d", n)
	}
	return result.TwoTheta, nil
}

func BraggEquation(lambda, d float64, thetaRad float64, n int) (float64, error) {
	if err := ValidateWavelength(lambda); err != nil {
		return 0, err
	}
	if err := ValidateSpacing(d); err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, fmt.Errorf("order must be >= 1")
	}
	return float64(n)*lambda - 2*d*math.Sin(thetaRad), nil
}

func MaxOrder(lambda, d float64) (int, error) {
	if err := ValidateWavelength(lambda); err != nil {
		return 0, err
	}
	if err := ValidateSpacing(d); err != nil {
		return 0, err
	}
	max := int(2 * d / lambda)
	if max < 1 {
		return 0, nil
	}
	return max, nil
}

func AllOrders(lambda, d float64) ([]BraggResult, error) {
	max, err := MaxOrder(lambda, d)
	if err != nil {
		return nil, err
	}
	results := make([]BraggResult, 0, max)
	for n := 1; n <= max; n++ {
		result, err := BraggAngle(lambda, d, n)
		if err != nil {
			return nil, err
		}
		if result.Possible {
			results = append(results, result)
		}
	}
	return results, nil
}

func WavelengthForAngle(thetaDeg, d float64, n int) (float64, error) {
	if err := ValidateSpacing(d); err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, fmt.Errorf("order must be >= 1")
	}
	if thetaDeg <= 0 || thetaDeg >= 90 {
		return 0, fmt.Errorf("theta must be in (0,90) degrees")
	}
	theta := thetaDeg * math.Pi / 180
	return 2 * d * math.Sin(theta) / float64(n), nil
}

func SpacingForAngle(thetaDeg, lambda float64, n int) (float64, error) {
	if err := ValidateWavelength(lambda); err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, fmt.Errorf("order must be >= 1")
	}
	if thetaDeg <= 0 || thetaDeg >= 90 {
		return 0, fmt.Errorf("theta must be in (0,90) degrees")
	}
	theta := thetaDeg * math.Pi / 180
	return float64(n) * lambda / (2 * math.Sin(theta)), nil
}

func IsPossible(lambda, d float64, n int) (bool, error) {
	result, err := BraggAngle(lambda, d, n)
	if err != nil {
		return false, err
	}
	return result.Possible, nil
}

func SinTheta(lambda, d float64, n int) (float64, error) {
	if err := ValidateWavelength(lambda); err != nil {
		return 0, err
	}
	if err := ValidateSpacing(d); err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, fmt.Errorf("order must be >= 1")
	}
	return float64(n) * lambda / (2 * d), nil
}

func AngleDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func AngleRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func FormatBragg(result BraggResult) string {
	return fmt.Sprintf(
		"n=%d lambda=%.4f d=%.4f theta=%.4f 2theta=%.4f possible=%v",
		result.N, result.Lambda, result.D, result.Theta, result.TwoTheta, result.Possible,
	)
}

func CheckEquation(result BraggResult, tolerance float64) float64 {
	theta := AngleRadians(result.Theta)
	value, err := BraggEquation(result.Lambda, result.D, theta, result.N)
	if err != nil {
		return math.NaN()
	}
	return math.Abs(value)
}

func LambdaLargerTwoThetaHigher(lambda1, lambda2, d float64) bool {
	t1, err1 := TwoTheta(lambda1, d, 1)
	t2, err2 := TwoTheta(lambda2, d, 1)
	if err1 != nil || err2 != nil {
		return false
	}
	return t2 > t1
}

func LargerDTwoThetaLower(lambda, d1, d2 float64) bool {
	t1, err1 := TwoTheta(lambda, d1, 1)
	t2, err2 := TwoTheta(lambda, d2, 1)
	if err1 != nil || err2 != nil {
		return false
	}
	return t2 < t1
}

func ExactRightAngle(lambda, d float64, n int) bool {
	sinTheta := float64(n) * lambda / (2 * d)
	return math.Abs(sinTheta-1) < 1e-12
}

func RightAngleResult(lambda, d float64, n int) BraggResult {
	return BraggResult{
		Lambda: lambda, D: d, N: n,
		SinTheta: 1, Theta: 90, TwoTheta: 180, Possible: true,
	}
}

func MaxTwoTheta(lambda, d float64) (float64, error) {
	max, err := MaxOrder(lambda, d)
	if err != nil {
		return 0, err
	}
	if max < 1 {
		return 0, nil
	}
	result, err := BraggAngle(lambda, d, max)
	if err != nil {
		return 0, err
	}
	return result.TwoTheta, nil
}

func MinTwoTheta(lambda, d float64) (float64, error) {
	result, err := BraggAngle(lambda, d, 1)
	if err != nil {
		return 0, err
	}
	return result.TwoTheta, nil
}

func OrderCount(lambda, d float64) (int, error) {
	return MaxOrder(lambda, d)
}

func TwoThetaList(lambda, d float64) ([]float64, error) {
	results, err := AllOrders(lambda, d)
	if err != nil {
		return nil, err
	}
	angles := make([]float64, 0, len(results))
	for _, result := range results {
		angles = append(angles, result.TwoTheta)
	}
	return angles, nil
}

func HighestOrder(lambda, d float64) (int, error) {
	return MaxOrder(lambda, d)
}

func IsOrderAllowed(n int) bool {
	return n >= 1
}

func OrderRange(lambda, d float64) (int, int, error) {
	max, err := MaxOrder(lambda, d)
	if err != nil {
		return 0, 0, err
	}
	return 1, max, nil
}

func SinLimit() float64 {
	return 1
}

func AngleForSin(sin float64) float64 {
	if sin > 1 {
		return 90
	}
	if sin < -1 {
		return -90
	}
	return math.Asin(sin) * 180 / math.Pi
}

func TwoThetaForSin(sin float64) float64 {
	return 2 * AngleForSin(sin)
}

func EquationResidual(result BraggResult) float64 {
	return CheckEquation(result, 0)
}

func IsPossibleAngle(result BraggResult) bool {
	return result.Possible
}

func InvalidOrderError(n int) error {
	return fmt.Errorf("order n must be >= 1, got %d", n)
}

func FormatAngle(angle float64) string {
	return fmt.Sprintf("%.4f deg", angle)
}

func Degrees(result BraggResult) (float64, float64) {
	return result.Theta, result.TwoTheta
}

func Radians(result BraggResult) (float64, float64) {
	return AngleRadians(result.Theta), AngleRadians(result.TwoTheta)
}

func SinThetaValue(result BraggResult) float64 {
	return result.SinTheta
}

func PossibleFlag(result BraggResult) bool {
	return result.Possible
}

func TwoThetaDegrees(result BraggResult) float64 {
	return result.TwoTheta
}

func ThetaDegrees(result BraggResult) float64 {
	return result.Theta
}

func IsPossibleBool(result BraggResult) bool {
	return result.Possible
}

func Lambda(result BraggResult) float64 {
	return result.Lambda
}

func DSpacing(result BraggResult) float64 {
	return result.D
}

func Order(result BraggResult) int {
	return result.N
}
