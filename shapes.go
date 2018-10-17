package demsphere

func NewIcosahedron() []Triangle {
	const a = 0.8506507174597755
	const b = 0.5257312591858783
	vertices := []Vector{
		{-a, -b, 0},
		{-a, b, 0},
		{-b, 0, -a},
		{-b, 0, a},
		{0, -a, -b},
		{0, -a, b},
		{0, a, -b},
		{0, a, b},
		{b, 0, -a},
		{b, 0, a},
		{a, -b, 0},
		{a, b, 0},
	}
	indices := [][3]int{
		{0, 3, 1},
		{1, 3, 7},
		{2, 0, 1},
		{2, 1, 6},
		{4, 0, 2},
		{4, 5, 0},
		{5, 3, 0},
		{6, 1, 7},
		{6, 7, 11},
		{7, 3, 9},
		{8, 2, 6},
		{8, 4, 2},
		{8, 6, 11},
		{8, 10, 4},
		{8, 11, 10},
		{9, 3, 5},
		{10, 5, 4},
		{10, 9, 5},
		{11, 7, 9},
		{11, 9, 10},
	}
	triangles := make([]Triangle, len(indices))
	for i, idx := range indices {
		p1 := vertices[idx[0]]
		p2 := vertices[idx[1]]
		p3 := vertices[idx[2]]
		triangles[i] = Triangle{p1, p2, p3}
	}
	return triangles
}
