package demsphere

import (
	"image"

	"github.com/fogleman/fauxgl"
)

type Triangulator struct {
	texture *Texture

	minDetail int
	maxDetail int
	minReal   float64
	maxReal   float64
	minRadius float64
	maxRadius float64
	tolerance float64

	triangles []*fauxgl.Triangle
}

func NewTriangulator(im image.Image, minDetail, maxDetail int, meanRadius, minElevation, maxElevation, tolerance, exaggeration, scale float64) *Triangulator {
	texture := NewTexture(im)
	minReal := meanRadius + minElevation
	maxReal := meanRadius + maxElevation
	minRadius := (meanRadius + minElevation*exaggeration) * scale
	maxRadius := (meanRadius + maxElevation*exaggeration) * scale
	return &Triangulator{texture, minDetail, maxDetail, minReal, maxReal, minRadius, maxRadius, tolerance, nil}
}

func (t *Triangulator) Triangulate() *fauxgl.Mesh {
	t.triangles = nil
	for _, triangle := range fauxgl.NewIcosahedron().Triangles {
		v1 := triangle.V1.Position
		v2 := triangle.V2.Position
		v3 := triangle.V3.Position
		t.triangulate(0, v1, v2, v3)
	}
	return fauxgl.NewTriangleMesh(t.triangles)
}

func (t *Triangulator) triangulate(detail int, v1, v2, v3 fauxgl.Vector) {
	if detail == t.maxDetail {
		t.leaf(v1, v2, v3)
		return
	}

	v12 := v1.Add(v2).DivScalar(2).Normalize()
	v13 := v1.Add(v3).DivScalar(2).Normalize()
	v23 := v2.Add(v3).DivScalar(2).Normalize()

	if detail >= t.minDetail {
		p1 := t.texture.Displace(v1, t.minReal, t.maxReal)
		p2 := t.texture.Displace(v2, t.minReal, t.maxReal)
		p3 := t.texture.Displace(v3, t.minReal, t.maxReal)
		plane := MakePlane(p1, p2, p3)
		if t.withinTolerance(3, plane, v1, v2, v3) {
			t.leaf(v1, v2, v3)
			return
		}
	}

	t.triangulate(detail+1, v1, v12, v13)
	t.triangulate(detail+1, v2, v23, v12)
	t.triangulate(detail+1, v3, v13, v23)
	t.triangulate(detail+1, v12, v23, v13)
}

func (t *Triangulator) leaf(v1, v2, v3 fauxgl.Vector) {
	p1 := t.texture.Displace(v1, t.minRadius, t.maxRadius)
	p2 := t.texture.Displace(v2, t.minRadius, t.maxRadius)
	p3 := t.texture.Displace(v3, t.minRadius, t.maxRadius)
	triangle := fauxgl.NewTriangleForPoints(p1, p2, p3)
	t.triangles = append(t.triangles, triangle)
}

func (t *Triangulator) withinTolerance(depth int, plane Plane, v1, v2, v3 fauxgl.Vector) bool {
	if depth == 0 {
		return true
	}

	v12 := v1.Add(v2).DivScalar(2).Normalize()
	p12 := t.texture.Displace(v12, t.minReal, t.maxReal)
	if plane.DistanceToPoint(p12) > t.tolerance {
		return false
	}

	v13 := v1.Add(v3).DivScalar(2).Normalize()
	p13 := t.texture.Displace(v13, t.minReal, t.maxReal)
	if plane.DistanceToPoint(p13) > t.tolerance {
		return false
	}

	v23 := v2.Add(v3).DivScalar(2).Normalize()
	p23 := t.texture.Displace(v23, t.minReal, t.maxReal)
	if plane.DistanceToPoint(p23) > t.tolerance {
		return false
	}

	return t.withinTolerance(depth-1, plane, v1, v12, v13) &&
		t.withinTolerance(depth-1, plane, v2, v23, v12) &&
		t.withinTolerance(depth-1, plane, v3, v13, v23) &&
		t.withinTolerance(depth-1, plane, v12, v23, v13)
}
