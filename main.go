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
	"flag"
)

type GIFMetadata struct {
	Width int
	Height int
	Delays []int
	Frames int
}

func StichFrames(g *gif.GIF, gifMeta *GIFMetadata) image.Image {
	t := image.NewUniform(color.Transparent)
	imgRect := image.Rect(0, 0, gifMeta.Width * len(g.Image), gifMeta.Height)
	img := image.NewRGBA(imgRect)
	draw.Draw(img, img.Bounds(), t, image.ZP, draw.Src)
	for i, e := range g.Image {
		index := image.Point{
			X: e.Bounds().Min.X + (i * gifMeta.Width),
			Y: e.Bounds().Min.Y,
		}
		
		r := img.Bounds()
		r.Min = r.Min.Add(index)
		draw.Draw(img, r, e, e.Bounds().Min, draw.Over)
	}

	return img
}

func ExtractGIFMetadata(g *gif.GIF) *GIFMetadata {
	width := g.Image[0].Rect.Max.X
	height := g.Image[0].Rect.Max.Y
	delays := g.Delay
		
	gifMeta := &GIFMetadata {
		Width: width,
		Height: height,
		Delays: delays,
		Frames: len(g.Image),
	}
	
	return gifMeta
}

func main() {
	var inputDir string
	var outputDir string

	flag.StringVar(&inputDir, "input-dir", "input", "input-dir: A directory of the gifs you want to convert to a sprite image")
	flag.StringVar(&outputDir, "output-dir", "output", "output-dir: A directory of your processed gifs with the stiched sprite images and metadata json files")

	flag.Parse()

	if (inputDir == "") {
		panic("Missing input directory. Please supply an input directory and try again.")
	}

	gifs, err := filepath.Glob(inputDir + "/*.gif")
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

		gifMeta := ExtractGIFMetadata(g)
		img := StichFrames(g, gifMeta)
		basefilen := filepath.Base(gz)
		filen := basefilen[0: len(basefilen) - 4] // better way to do this, I'm sure
		
		err = os.Mkdir(outputDir + "/" + filen, 0755)
		if err != nil {
			switch err.(type) {
			case *os.PathError:
				// no-op
			default:
				panic(err)
			}
		}
		
		p, err := os.Create(outputDir + "/" + filen + "/" + filen + ".png")
		if err != nil {
			panic(err)
		}
		
		err = png.Encode(p, img)
		if err != nil {
			panic(err)
		}
		
		d, err := os.Create(outputDir + "/" + filen + "/" + filen + ".json")
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

	fmt.Printf("Converted %d gifs and placed them in %s", len(gifs), outputDir)  
}