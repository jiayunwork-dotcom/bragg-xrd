package server

import (
	"net/http"

	"bragg-xrd/internal/xrd"
)

type braggRequest struct {
	Lambda float64 `json:"lambda"`
	D      float64 `json:"d"`
	A      float64 `json:"a,omitempty"`
	N      int     `json:"n,omitempty"`
	HKL    []int   `json:"hkl,omitempty"`
}

type powderRequest struct {
	Lambda   float64 `json:"lambda"`
	A        float64 `json:"a"`
	Lattice  string  `json:"lattice"`
	MaxIndex int     `json:"max_index,omitempty"`
}

func braggHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req braggRequest
	if !readJSON(w, r, &req) {
		return
	}
	n := req.N
	if n == 0 {
		n = 1
	}
	d := req.D
	if d <= 0 && req.A > 0 {
		hkl, err := xrd.HKLFromArray(req.HKL)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		d, err = xrd.LatticeSpacing(req.A, hkl)
		if err != nil {
			badRequest(w, err.Error())
			return
		}
	}
	result, err := xrd.BraggAngle(req.Lambda, d, n)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func powderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req powderRequest
	if !readJSON(w, r, &req) {
		return
	}
	if req.MaxIndex <= 0 {
		req.MaxIndex = 4
	}
	result, err := xrd.Powder(req.Lambda, req.A, req.Lattice, req.MaxIndex)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": "bragg-xrd", "version": "1.0.0"})
}
