package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	width     = 1000
	height    = 1000
	createNum = 100
)

type ColorPallet struct {
	Result [][]int `json:"result"`
}

func userColorPallet(colorTheme string) []color.RGBA {
	bytes := []byte(colorTheme)
	var colorPallet ColorPallet
	json.Unmarshal(bytes, &colorPallet)
	rgbas := make([]color.RGBA, len(colorPallet.Result))
	for i, cs := range colorPallet.Result {
		rgbas[i] = color.RGBA{
			R: uint8(cs[0]),
			G: uint8(cs[1]),
			B: uint8(cs[2]),
			A: 140,
		}
	}
	return rgbas
}

func getColorPallet() []color.RGBA {
	body := strings.NewReader("{\"model\":\"default\"}")
	req, err := http.NewRequest("POST", "http://colormind.io/api/", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var colorPallet ColorPallet
	err = json.NewDecoder(resp.Body).Decode(&colorPallet)
	if err != nil {
		panic(err)
	}
	rgbas := make([]color.RGBA, len(colorPallet.Result))
	for i, cs := range colorPallet.Result {
		rgbas[i] = color.RGBA{
			R: uint8(cs[0]),
			G: uint8(cs[1]),
			B: uint8(cs[2]),
			A: 140,
		}
	}
	return rgbas
}

func gen(fileIndex int, colorTheme string) {
	var pallet []color.RGBA
	if colorTheme != "" {
		pallet = userColorPallet(colorTheme)
	} else {
		pallet = getColorPallet()
	}

	rand.Seed(time.Now().UnixNano())
	random1X := rand.Intn(width)
	random1Y := rand.Intn(height)

	random2X := rand.Intn(width)
	random2Y := rand.Intn(height)

	random3X := rand.Intn(width)
	random3Y := rand.Intn(height)

	random4X := rand.Intn(width)
	random4Y := rand.Intn(height)

	img0 := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img0.Set(x, y, pallet[0])
		}
	}

	img1 := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < random1X {
				if y > random1Y {
					img1.Set(x, y, pallet[1])
				}
			}
		}
	}

	img2 := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x > width-random2X {
				if y > height-random2Y {
					img2.Set(x, y, pallet[2])
				}
			}
		}
	}

	img3 := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x > width-random3X {
				if y < random3Y {
					img3.Set(x, y, pallet[3])
				}
			}
		}
	}

	img4 := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < random4X {
				if y < random4Y {
					img4.Set(x, y, pallet[4])
				}
			}
		}
	}
	draw.Draw(img0, img1.Rect, img1, image.Point{0, 0}, draw.Over)
	draw.Draw(img0, img2.Rect, img2, image.Point{0, 0}, draw.Over)
	draw.Draw(img0, img3.Rect, img3, image.Point{0, 0}, draw.Over)
	draw.Draw(img0, img4.Rect, img4, image.Point{0, 0}, draw.Over)

	f, _ := os.Create(fmt.Sprintf("./img/image%d.png", fileIndex))
	defer f.Close()
	png.Encode(f, img0)
	fmt.Println("your color theme is")

	theme := mapPalletToColorTheme(pallet)
	marshal, err := json.Marshal(theme)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

func mapPalletToColorTheme(pallet []color.RGBA) ColorPallet {
	var colorPallet ColorPallet
	v := make([][]int, len(pallet))
	for i, rgba := range pallet {
		v[i] = []int{int(rgba.R), int(rgba.G), int(rgba.B)}
	}
	colorPallet.Result = v
	return colorPallet
}

func main() {
	fmt.Print("is exit your color theme input?")
	// Scannerを使って一行読み
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	colorTheme := scanner.Text()
	for i := 0; i < createNum; i++ {
		gen(i, colorTheme)
	}
}
