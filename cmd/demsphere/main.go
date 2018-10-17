package main

import (
	"fmt"
	"log"

	"github.com/fogleman/demsphere"
	"github.com/fogleman/fauxgl"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	inputFile  = kingpin.Flag("input", "Input DEM image to process.").Required().Short('i').ExistingFile()
	outputFile = kingpin.Flag("output", "Output STL file to write.").Required().Short('o').String()
)

func main() {
	kingpin.Parse()

	im, err := fauxgl.LoadImage(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	triangulator := demsphere.NewTriangulator(
		im, 6, 11, 1737400, -18257, 21563, 50, 3, 1.0/1737400)

	mesh := triangulator.Triangulate()
	fmt.Println(len(mesh.Triangles))

	inner := fauxgl.NewSphere(4)
	inner.Transform(fauxgl.Scale(fauxgl.V(0.85, 0.85, 0.85)))
	inner.ReverseWinding()
	mesh.Add(inner)

	fmt.Println(len(mesh.Triangles))

	mesh.SaveSTL(*outputFile)
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
