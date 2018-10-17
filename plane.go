package demsphere

import "github.com/fogleman/fauxgl"

type Plane struct {
	N fauxgl.Vector
	D float64
}

func MakePlane(p1, p2, p3 fauxgl.Vector) Plane {
	n := p2.Sub(p1).Cross(p3.Sub(p1)).Normalize()
	d := n.Dot(p1)
	return Plane{n, d}
}

func (p *Plane) DistanceToPoint(q fauxgl.Vector) float64 {
	x := q.Dot(p.N) - p.D
	if x < 0 {
		return -x
	}
	return x
}
