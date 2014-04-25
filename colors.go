package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"io"
	"sort"

	_ "image/jpeg"
	_ "image/png"
)

type Pair struct {
	Key   string
	Value int
}

type ColorList []Pair

func (p ColorList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ColorList) Len() int           { return len(p) }
func (p ColorList) Less(i, j int) bool { return p[i].Value > p[j].Value }

func sortColorList(m map[string]int) ColorList {
	p := make(ColorList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

type ImageColor struct {
	Colors map[string]int
}

func New(reader io.Reader) (*ImageColor, error) {
	d, _, err := image.Decode(reader)

	if err != nil {
		return nil, err
	}

	bounds := d.Bounds()

	colors := make(map[string]int, 0)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := d.At(x, y)
			colors[toHex(color)]++
		}
	}

	return &ImageColor{Colors: colors}, nil
}

func (ic *ImageColor) TopColors(n int) []string {
	colors := make([]string, 0)
	colorList := sortColorList(ic.Colors)
	for i, color := range colorList {
		colorsCount := len(colors)
		colorHex := color.Key
		if colorsCount <= n {
			if colorsCount > 0 {
				prev, _ := colorful.Hex(colorList[i-1].Key)
				current, _ := colorful.Hex(color.Key)
				distance := prev.DistanceRgb(current)
				if distance > 0.25 {
					colors = append(colors, colorHex)
				}
			} else {
				colors = append(colors, colorHex)
			}
		}
	}
	return colors
}

func toHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return string(fmt.Sprintf("#%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8)))
}