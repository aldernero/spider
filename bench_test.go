package spider

import (
	"fmt"
	"path/filepath"
	"testing"
)

// benchChart builds a chart with the given number of axes and series.
func benchChart(nAxes, nSeries int) *Chart {
	c := NewChart()
	names := make([]string, nAxes)
	for i := range names {
		names[i] = fmt.Sprintf("axis%d", i)
		if err := c.AddAxis(names[i]); err != nil {
			panic(err)
		}
	}
	for i := 0; i < nSeries; i++ {
		data := make(map[string]float64, nAxes)
		for j, name := range names {
			data[name] = float64((i+1)*(j+1)) + 1
		}
		if err := c.AddSeries(fmt.Sprintf("series%d", i), data); err != nil {
			panic(err)
		}
	}
	return c
}

func BenchmarkSaveSVG(b *testing.B) {
	c := benchChart(3, 1)
	out := filepath.Join(b.TempDir(), "out.svg")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.SaveSVG(out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSavePNG(b *testing.B) {
	c := benchChart(3, 1)
	out := filepath.Join(b.TempDir(), "out.png")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.SavePNG(out); err != nil {
			b.Fatal(err)
		}
	}
}

// validate resolves every font face. It runs on each Draw, so its cost is paid
// per render and the font family cache is what keeps it cheap.
func BenchmarkValidate(b *testing.B) {
	c := benchChart(3, 1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.validate(); err != nil {
			b.Fatal(err)
		}
	}
}

// Axis maxima used to be recomputed per axis per series while drawing, which is
// quadratic in the series count.
func BenchmarkAxisMaxima(b *testing.B) {
	for _, size := range []struct{ axes, series int }{{5, 5}, {20, 20}} {
		b.Run(fmt.Sprintf("%dx%d", size.axes, size.series), func(b *testing.B) {
			c := benchChart(size.axes, size.series)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.axisMaxima()
			}
		})
	}
}

func BenchmarkSaveSVGLarge(b *testing.B) {
	c := benchChart(20, 20)
	out := filepath.Join(b.TempDir(), "out.svg")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.SaveSVG(out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkColorToCanvasColor(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Color("#3B82F6").ToCanvasColor()
		_ = Color("cornflowerblue").ToCanvasColor()
	}
}
