package main

import ( 
	"image/gif"
	"image"
	"image/draw"
	"image/color"
	"os"
	"fmt"
	"image/png"
	"encoding/json"
	"path/filepath"
)

type GIFMetadata struct {
	Width int
	Height int
	Delays []int
	Frames int
}

func main() {
	gifs, err := filepath.Glob("../../data/raw/*.gif")
	if err != nil {
		panic(err)
	}

	for _, gz := range gifs {
		f, err := os.Open(gz)
		if err != nil {
			panic(err)
		}
		
		g, err := gif.DecodeAll(f)
		if err != nil {
			panic(err)
		}
		
		width := g.Image[0].Rect.Max.X
		height := g.Image[0].Rect.Max.Y
		delays := g.Delay
		
		gifMeta := &GIFMetadata {
			Width: width,
			Height: height,
			Delays: delays,
			Frames: len(g.Image),
		}
		
		t := image.NewUniform(color.Transparent)
		imgRect := image.Rect(0, 0, width * len(g.Image), height)
		img := image.NewRGBA(imgRect)
		draw.Draw(img, img.Bounds(), t, image.ZP, draw.Src)
		for j, e := range g.Image {
			index := image.Point{
				X: e.Bounds().Min.X + (j * width),
				Y: e.Bounds().Min.Y,
			}
			
			r := img.Bounds()
			r.Min = r.Min.Add(index)
			draw.Draw(img, r, e, e.Bounds().Min, draw.Over)
		}

		basefilen := filepath.Base(gz)
		filen := basefilen[0: len(basefilen) - 4] // better way to do this, I'm sure
		
		err = os.Mkdir("../../data/processed/" + filen, 0755)
		if err != nil {
			switch err.(type) {
			case *os.PathError:
				// no-op
			default:
				panic(err)
			}
		}
		
		p, err := os.Create("../../data/processed/" + filen + "/" + filen + ".png")
		if err != nil {
			panic(err)
		}
		
		err = png.Encode(p, img)
		if err != nil {
			panic(err)
		}
		
		d, err := os.Create("../../data/processed/" + filen + "/" + filen + ".json")
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(gifMeta)
		if err != nil {
			panic(err)
		}
		
		d.Write(b)
	
		fmt.Printf("%s Frames: %d\n", filen, len(g.Image))
	}
}