package moremath

import "math"

// math.Min doen't comply with the Wasm spec, so we borrow from the original
// with a change that either one of NaN results in NaN even if another is -Inf.
// https://github.com/golang/go/blob/1d20a362d0ca4898d77865e314ef6f73582daef0/src/math/dim.go#L74-L91
func WasmCompatMin(x, y float64) float64 {
	switch {
	case math.IsNaN(x) || math.IsNaN(y):
		return math.NaN()
	case math.IsInf(x, -1) || math.IsInf(y, -1):
		return math.Inf(-1)
	case x == 0 && x == y:
		if math.Signbit(x) {
			return x
		}
		return y
	}
	if x < y {
		return x
	}
	return y
}

// math.Max doen't comply with the Wasm spec, so we borrow from the original
// with a change that either one of NaN results in NaN even if another is Inf.
// https://github.com/golang/go/blob/1d20a362d0ca4898d77865e314ef6f73582daef0/src/math/dim.go#L42-L59
func WasmCompatMax(x, y float64) float64 {
	switch {
	case math.IsNaN(x) || math.IsNaN(y):
		return math.NaN()
	case math.IsInf(x, 1) || math.IsInf(y, 1):
		return math.Inf(1)

	case x == 0 && x == y:
		if math.Signbit(x) {
			return y
		}
		return x
	}
	if x > y {
		return x
	}
	return y
}

// WasmCompatNearestF32 is the Wasm spec compatible variant of math.Round, which is used for Nearest instruction.
// For example, this converts 1.9 to 2.0, and this has the semantics of LLVM's rint instrinsic: https://llvm.org/docs/LangRef.html#llvm-rint-intrinsic.
// For the difference from math.Round, math.Round(-4.5) results in -5 while this produces -4.
func WasmCompatNearestF32(f float32) float32 {
	// TODO: look at https://github.com/bytecodealliance/wasmtime/pull/2171 and reconsider this algorithm
	if f != -0 && f != 0 {
		ceil := float32(math.Ceil(float64(f)))
		floor := float32(math.Floor(float64(f)))
		distToCeil := math.Abs(float64(f - ceil))
		distToFloor := math.Abs(float64(f - floor))
		h := ceil / 2.0
		if distToCeil < distToFloor {
			f = ceil
		} else if distToCeil == distToFloor && float32(math.Floor(float64(h))) == h {
			f = ceil
		} else {
			f = floor
		}
	}
	return f
}

// WasmCompatNearestF64 is the Wasm spec compatible variant of math.Round, which is used for Nearest instruction.
// For example, this converts 1.9 to 2.0, and this has the semantics of LLVM's rint instrinsic: https://llvm.org/docs/LangRef.html#llvm-rint-intrinsic.
// For the difference from math.Round, math.Round(-4.5) results in -5 while this produces -4.
func WasmCompatNearestF64(f float64) float64 {
	// TODO: look at https://github.com/bytecodealliance/wasmtime/pull/2171 and reconsider this algorithm
	if f != -0 && f != 0 {
		ceil := math.Ceil(f)
		floor := math.Floor(f)
		distToCeil := math.Abs(f - ceil)
		distToFloor := math.Abs(f - floor)
		h := ceil / 2.0
		if distToCeil < distToFloor {
			f = ceil
		} else if distToCeil == distToFloor && math.Floor(float64(h)) == h {
			f = ceil
		} else {
			f = floor
		}
	}
	return f
}
