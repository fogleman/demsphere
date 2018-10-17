package demsphere

import (
	"image"
)

type Triangulator struct {
	texture *Texture

	minDetail       int
	maxDetail       int
	minRadius       float64
	maxRadius       float64
	minOutputRadius float64
	maxOutputRadius float64
	tolerance       float64

	points    map[Vector]Vector
	temp      []Triangle
	triangles []Triangle
}

func NewTriangulator(im image.Image, minDetail, maxDetail int, meanRadius, minElevation, maxElevation, tolerance, exaggeration, scale float64) *Triangulator {
	texture := NewTexture(im)
	minRadius := meanRadius + minElevation
	maxRadius := meanRadius + maxElevation
	minOutputRadius := (meanRadius + minElevation*exaggeration) * scale
	maxOutputRadius := (meanRadius + maxElevation*exaggeration) * scale
	points := make(map[Vector]Vector)
	return &Triangulator{texture, minDetail, maxDetail, minRadius, maxRadius, minOutputRadius, maxOutputRadius, tolerance, points, nil, nil}
}

func (tri *Triangulator) Triangulate() []Triangle {
	tri.temp = nil
	tri.triangles = nil
	for _, t := range NewIcosahedron() {
		tri.triangulate(0, t.A, t.B, t.C)
	}
	for _, t := range tri.temp {
		tri.split(t.A, t.B, t.C)
	}
	return tri.triangles
}

func (tri *Triangulator) split(v1, v2, v3 Vector) {
	v12 := bisect(v1, v2)
	v23 := bisect(v2, v3)
	v31 := bisect(v3, v1)
	if _, ok := tri.points[v12]; ok {
		tri.split(v1, v12, v3)
		tri.split(v12, v2, v3)
	} else if _, ok := tri.points[v23]; ok {
		tri.split(v1, v2, v23)
		tri.split(v23, v3, v1)
	} else if _, ok := tri.points[v31]; ok {
		tri.split(v1, v2, v31)
		tri.split(v31, v2, v3)
	} else {
		p1 := tri.points[v1]
		p2 := tri.points[v2]
		p3 := tri.points[v3]
		tri.triangles = append(tri.triangles, Triangle{p1, p2, p3})
	}
}

func (tri *Triangulator) triangulate(detail int, v1, v2, v3 Vector) {
	if detail == tri.maxDetail {
		tri.leaf(v1, v2, v3)
		return
	}

	v12 := bisect(v1, v2)
	v23 := bisect(v2, v3)
	v31 := bisect(v3, v1)

	if detail >= tri.minDetail {
		p1 := tri.texture.Displace(v1, tri.minRadius, tri.maxRadius)
		p2 := tri.texture.Displace(v2, tri.minRadius, tri.maxRadius)
		p3 := tri.texture.Displace(v3, tri.minRadius, tri.maxRadius)
		plane := MakePlane(p1, p2, p3)
		depth := tri.maxDetail - detail + 1
		if depth > 6 {
			depth = 6
		}
		if tri.withinTolerance(depth, plane, v1, v2, v3) {
			tri.leaf(v1, v2, v3)
			return
		}
	}

	tri.triangulate(detail+1, v1, v12, v31)
	tri.triangulate(detail+1, v2, v23, v12)
	tri.triangulate(detail+1, v3, v31, v23)
	tri.triangulate(detail+1, v12, v23, v31)
}

func (tri *Triangulator) leaf(v1, v2, v3 Vector) {
	p1 := tri.texture.Displace(v1, tri.minOutputRadius, tri.maxOutputRadius)
	p2 := tri.texture.Displace(v2, tri.minOutputRadius, tri.maxOutputRadius)
	p3 := tri.texture.Displace(v3, tri.minOutputRadius, tri.maxOutputRadius)
	tri.points[v1] = p1
	tri.points[v2] = p2
	tri.points[v3] = p3
	tri.temp = append(tri.temp, Triangle{v1, v2, v3})
}

func (tri *Triangulator) withinTolerance(depth int, plane Plane, v1, v2, v3 Vector) bool {
	if depth == 0 {
		return true
	}

	v12 := bisect(v1, v2)
	p12 := tri.texture.Displace(v12, tri.minRadius, tri.maxRadius)
	if plane.DistanceToPoint(p12) > tri.tolerance {
		return false
	}

	v23 := bisect(v2, v3)
	p23 := tri.texture.Displace(v23, tri.minRadius, tri.maxRadius)
	if plane.DistanceToPoint(p23) > tri.tolerance {
		return false
	}

	v31 := bisect(v3, v1)
	p13 := tri.texture.Displace(v31, tri.minRadius, tri.maxRadius)
	if plane.DistanceToPoint(p13) > tri.tolerance {
		return false
	}

	if depth == 1 {
		return true
	}

	return tri.withinTolerance(depth-1, plane, v1, v12, v31) &&
		tri.withinTolerance(depth-1, plane, v2, v23, v12) &&
		tri.withinTolerance(depth-1, plane, v3, v31, v23) &&
		tri.withinTolerance(depth-1, plane, v12, v23, v31)
}
