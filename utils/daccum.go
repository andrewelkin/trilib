package utils

import (
	"fmt"
	"math"
)

type Daccum struct {
	t, w float64
}

func (d *Daccum) Sum() float64 {
	return d.t
}

func (d *Daccum) SumWeights() float64 {
	return d.w
}

func (d *Daccum) Add(a1, b1 float64) {
	b1 = math.Abs(b1)
	d.t += a1 * b1
	d.w += b1
}

func (d *Daccum) ToString() string {
	return fmt.Sprintf("%f %f avg=%f sumw=%f", d.t, d.w, d.Avg(), d.SumWeights())
}

func (d *Daccum) Avg() float64 {
	if d.w == 0 {
		return 0
	}
	return d.t / d.w
}

func (d *Daccum) Clear() {
	d.w = 0
	d.t = 0
}

func (d *Daccum) Scale(sc float64) {
	d.t *= sc
	d.w *= sc
	if d.w < 1e-16 {
		d.Clear() //otherwise precision is lost
	}
}

func EMavWeight(n float64) float64 {
	return math.Pow(0.5, 1./math.Max(n, 1.))
}
