package demsphere

type Triangle struct {
	A, B, C Vector
}

func (t *Triangle) Normal() Vector {
	ab := t.B.Sub(t.A)
	ac := t.C.Sub(t.A)
	return ab.Cross(ac).Normalize()
}
