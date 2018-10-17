package demsphere

type Plane struct {
	N Vector
	D float64
}

func MakePlane(p1, p2, p3 Vector) Plane {
	n := p2.Sub(p1).Cross(p3.Sub(p1)).Normalize()
	d := n.Dot(p1)
	return Plane{n, d}
}

func (p *Plane) DistanceToPoint(q Vector) float64 {
	x := q.Dot(p.N) - p.D
	if x < 0 {
		return -x
	}
	return x
}
