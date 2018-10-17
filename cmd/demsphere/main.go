package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fogleman/demsphere"
	"github.com/fogleman/fauxgl"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	inputFile  = kingpin.Flag("input", "Input DEM image to process.").Required().Short('i').ExistingFile()
	outputFile = kingpin.Flag("output", "Output STL file to write.").Required().Short('o').String()
)

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func main() {
	var done func()

	kingpin.Parse()

	done = timed("reading input")
	im, err := fauxgl.LoadImage(*inputFile)
	done()
	if err != nil {
		log.Fatal(err)
	}

	// mercury
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 12, 2439400, -10764, 8994, 50, 4, 1.0/2439400)

	// moon
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 11, 1737400, -18257, 21563, 50, 3, 1.0/1737400)

	// mars
	triangulator := demsphere.NewTriangulator(
		im, 0, 10, 3396190, -8201, 21241, 1000, 10, 1.0/3396190)

	// pluto
	// triangulator := demsphere.NewTriangulator(
	// 	im, 6, 12, 1188300, -4101, 6491, 50, 3, 1.0/1188300)

	done = timed("generating mesh")
	mesh := triangulator.Triangulate()
	done()
	fmt.Println(len(mesh.Triangles))

	inner := fauxgl.NewSphere(4)
	inner.Transform(fauxgl.Scale(fauxgl.V(0.85, 0.85, 0.85)))
	inner.ReverseWinding()
	mesh.Add(inner)

	fmt.Println(len(mesh.Triangles))

	done = timed("writing output")
	mesh.SaveSTL(*outputFile)
	done()
}

// 4,5120,4.7372172692
// 5,20480,2.3686086346009336
// 6,81920,1.1844992435794788
// 7,327680,0.5635352519913389
// 8,1310720,0.2833191921173619
// 9,5242880,0.14087309976888143
// 10,20971520,0.07043426041304851
// 11,83886080,0.036106462472517295
// 12,335544320,0.018169800357500515
