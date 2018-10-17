package demsphere

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

func WriteSTLFile(path string, triangles []Triangle) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	return WriteSTL(w, triangles)
}

func WriteSTL(w io.Writer, triangles []Triangle) error {
	type STLHeader struct {
		_     [80]uint8
		Count uint32
	}

	type STLTriangle struct {
		N, V1, V2, V3 [3]float32
		_             uint16
	}

	header := STLHeader{}
	header.Count = uint32(len(triangles))
	if err := binary.Write(w, binary.LittleEndian, &header); err != nil {
		return err
	}

	for _, triangle := range triangles {
		n := triangle.Normal()
		d := STLTriangle{}
		d.N[0] = float32(n.X)
		d.N[1] = float32(n.Y)
		d.N[2] = float32(n.Z)
		d.V1[0] = float32(triangle.A.X)
		d.V1[1] = float32(triangle.A.Y)
		d.V1[2] = float32(triangle.A.Z)
		d.V2[0] = float32(triangle.B.X)
		d.V2[1] = float32(triangle.B.Y)
		d.V2[2] = float32(triangle.B.Z)
		d.V3[0] = float32(triangle.C.X)
		d.V3[1] = float32(triangle.C.Y)
		d.V3[2] = float32(triangle.C.Z)
		if err := binary.Write(w, binary.LittleEndian, &d); err != nil {
			return err
		}
	}

	return nil
}
