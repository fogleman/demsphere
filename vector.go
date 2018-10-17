package demsphere

import "math"

func bisect(v1, v2 Vector) Vector {
	// v1.Add(v2).DivScalar(2).Normalize()
	x := (v1.X + v2.X) / 2
	y := (v1.Y + v2.Y) / 2
	z := (v1.Z + v2.Z) / 2
	r := 1 / math.Sqrt(x*x+y*y+z*z)
	return Vector{x * r, y * r, z * r}
}

type Vector struct {
	X, Y, Z float64
}

func (a Vector) Add(b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Vector) Cross(b Vector) Vector {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return Vector{x, y, z}
}

func (a Vector) DivScalar(b float64) Vector {
	return Vector{a.X / b, a.Y / b, a.Z / b}
}

func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) MulScalar(b float64) Vector {
	return Vector{a.X * b, a.Y * b, a.Z * b}
}

func (a Vector) Normalize() Vector {
	r := 1 / math.Sqrt(a.X*a.X+a.Y*a.Y+a.Z*a.Z)
	return Vector{a.X * r, a.Y * r, a.Z * r}
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}
