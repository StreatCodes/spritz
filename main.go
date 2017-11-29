package main

import (
	"encoding/json"
	"image"
	"log"
	"os"

	_ "image/png"
)

type Sprite struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

type Point struct {
	X int
	Y int
}

var img image.Image
var checked [][]bool

func main() {
	//Load image file
	f, err := os.Open("sheet.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err = image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	checked = make([][]bool, bounds.Max.Y)
	for i := 0; i < len(checked); i++ {
		checked[i] = make([]bool, bounds.Max.X)
	}

	var sprites []Sprite

	//Check pixels from left to right & top to bottom for sprites
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if !checked[y][x] && a > 0 {
				//Begining of sprite found, record it
				sprite := Sprite{Top: y, Bottom: y, Right: x, Left: x}
				findBounds(x, y, &sprite)

				//Add sprite bounding box to the list
				sprites = append(sprites, sprite)
			}
		}
	}

	//Output result as json
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(sprites)
	if err != nil {
		log.Fatal(err)
	}
}

func findBounds(x, y int, sprite *Sprite) {
	//This pixel is outside the edge of the image, don't check it.
	if x < 0 || x >= img.Bounds().Max.X || y < 0 || y >= img.Bounds().Max.Y {
		return
	}
	//This pixel is blank or has already been checked, don't check it.
	if _, _, _, a := img.At(x, y).RGBA(); a < 1 || checked[y][x] { //TODO add checked
		return
	}
	checked[y][x] = true

	//Deterimine if this pixel is outside the current sprite bounding box, if so extend it
	if x < sprite.Left {
		sprite.Left = x
	} else if x > sprite.Right {
		sprite.Right = x
	}
	if y < sprite.Top {
		sprite.Top = y
	} else if sprite.Bottom < y {
		sprite.Bottom = y
	}

	//Check the 8 surrounding pixels
	findBounds(x, y-1, sprite)
	findBounds(x+1, y-1, sprite)
	findBounds(x+1, y, sprite)
	findBounds(x+1, y+1, sprite)
	findBounds(x, y+1, sprite)
	findBounds(x-1, y+1, sprite)
	findBounds(x-1, y, sprite)
	findBounds(x-1, y-1, sprite)
}
