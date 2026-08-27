package xrd

import "fmt"

type HKL struct {
	H int `json:"h"`
	K int `json:"k"`
	L int `json:"l"`
}

type BraggResult struct {
	Lambda   float64 `json:"lambda"`
	D        float64 `json:"d"`
	N        int     `json:"n"`
	Theta    float64 `json:"theta_deg"`
	TwoTheta float64 `json:"two_theta_deg"`
	SinTheta float64 `json:"sin_theta"`
	Possible bool    `json:"possible"`
}

type Peak struct {
	HKL       HKL     `json:"hkl"`
	TwoTheta  float64 `json:"two_theta_deg"`
	Theta     float64 `json:"theta_deg"`
	Order     int     `json:"order"`
	Intensity float64 `json:"intensity,omitempty"`
}

func (h HKL) String() string {
	return fmt.Sprintf("(%d%d%d)", h.H, h.K, h.L)
}

func (h HKL) Sum() int {
	return h.H + h.K + h.L
}

func (h HKL) IsZero() bool {
	return h.H == 0 && h.K == 0 && h.L == 0
}

func (h HKL) Parity() (bool, bool, bool) {
	return h.H%2 == 0, h.K%2 == 0, h.L%2 == 0
}

func (h HKL) AllEvenOrAllOdd() bool {
	he, ke, le := h.Parity()
	return (he && ke && le) || (!he && !ke && !le)
}

func (h HKL) Abs() HKL {
	return HKL{
		H: absInt(h.H), K: absInt(h.K), L: absInt(h.L),
	}
}

func (p Peak) String() string {
	return fmt.Sprintf("%s %.3f deg n=%d", p.HKL, p.TwoTheta, p.Order)
}

func (r BraggResult) IsFinite() bool {
	return !bad(r.Theta) && !bad(r.TwoTheta) && !bad(r.SinTheta)
}

func bad(value float64) bool {
	return value != value || value > 1e300 || value < -1e300
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
