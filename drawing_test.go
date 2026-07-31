package spider

import (
	"math"
	"testing"
)

func TestLinspace(t *testing.T) {
	tests := []struct {
		name string
		i, j float64
		n    int
		want []float64
	}{
		{"two points are the endpoints", 0, 10, 2, []float64{0, 10}},
		{"inclusive of both ends", 0, 4, 5, []float64{0, 1, 2, 3, 4}},
		{"non zero start", 2, 6, 3, []float64{2, 4, 6}},
		{"descending range", 10, 0, 3, []float64{10, 5, 0}},
		{"single point", 3, 9, 1, []float64{3}},
		// n-1 would divide by zero for n < 2, which used to yield Inf values
		{"zero points", 0, 10, 0, nil},
		{"negative count", 0, 10, -3, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := linspace(tt.i, tt.j, tt.n)
			if len(got) != len(tt.want) {
				t.Fatalf("linspace(%v, %v, %d) = %v, want %v", tt.i, tt.j, tt.n, got, tt.want)
			}
			for k := range got {
				if math.Abs(got[k]-tt.want[k]) > 1e-9 {
					t.Fatalf("linspace(%v, %v, %d) = %v, want %v", tt.i, tt.j, tt.n, got, tt.want)
				}
				if math.IsInf(got[k], 0) || math.IsNaN(got[k]) {
					t.Fatalf("linspace produced a non-finite value: %v", got)
				}
			}
		})
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		a, b, i, want float64
	}{
		{0, 10, 0, 0},
		{0, 10, 1, 10},
		{0, 10, 0.5, 5},
		{10, 20, 0.25, 12.5},
		{-10, 10, 0.5, 0},
		{0, 10, 2, 20}, // extrapolates rather than clamping
	}
	for _, tt := range tests {
		if got := lerp(tt.a, tt.b, tt.i); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("lerp(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.i, got, tt.want)
		}
	}
}

func TestLinmap(t *testing.T) {
	tests := []struct {
		name          string
		a, b, c, d, i float64
		want          float64
	}{
		{"identity", 0, 1, 0, 1, 0.5, 0.5},
		{"scale up", 0, 10, 0, 100, 5, 50},
		{"scale down", 0, 100, 0, 1, 50, 0.5},
		{"offset range", 0, 10, 100, 200, 5, 150},
		{"inverted target", 0, 10, 10, 0, 2, 8},
		{"value maps a data point onto a radius", 0, 4000, 0, 80, 1000, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := linmap(tt.a, tt.b, tt.c, tt.d, tt.i); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("linmap(%v, %v, %v, %v, %v) = %v, want %v", tt.a, tt.b, tt.c, tt.d, tt.i, got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		v, min, max, want float64
	}{
		{5, 0, 10, 5},
		{-1, 0, 10, 0},
		{99, 0, 10, 10},
		{5, 0, math.Inf(1), 5},
		{5, 10, math.Inf(1), 10},
	}
	for _, tt := range tests {
		if got := clamp(tt.v, tt.min, tt.max); got != tt.want {
			t.Errorf("clamp(%v, %v, %v) = %v, want %v", tt.v, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestVal2String(t *testing.T) {
	tests := []struct {
		val  float64
		want string
	}{
		{0, "0.000"},
		{0.5, "0.500"},
		{0.9999, "1.000"},
		{1, "1.00"},
		{5.25, "5.25"},
		{9.999, "10.00"},
		{10, "10.0"},
		{50.25, "50.2"},
		{100, "100"},
		{500.4, "500"},
		{1000, "1k"},
		{500000, "500k"},
		{1000000, "1M"},
		{500000000, "500M"},
		{1000000000, "1G"},
		{5e11, "500G"},
		{1e12, "1T"},
		{5e14, "500T"},
		{1e15, "1P"},
		{5e17, "500P"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := val2String(tt.val); got != tt.want {
				t.Errorf("val2String(%v) = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

// Negative values are below every magnitude threshold, so they take the
// smallest-value branch regardless of size.
func TestVal2StringNegative(t *testing.T) {
	for _, val := range []float64{-1, -1000, -1e9} {
		if got := val2String(val); got == "" {
			t.Errorf("val2String(%v) returned an empty string", val)
		}
	}
}
