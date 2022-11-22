package main

import (
    "bytes"
    "github.com/biessek/golang-ico"
    "image"
    "image/color"
)

var img = image.NewRGBA(image.Rect(0, 0, 32, 32))

// HLine draws a horizontal line
func HLine(x1, y, x2 int, col color.Color) {
    for ; x1 <= x2; x1++ {
        img.Set(x1, y, col)
    }
}

// VLine draws a veritcal line
func VLine(x, y1, y2 int, col color.Color) {
    for ; y1 <= y2; y1++ {
        img.Set(x, y1, col)
    }
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(x1, y1, x2, y2 int, col color.Color) {
    HLine(x1, y1, x2, col)
    HLine(x1, y2, x2, col)
    VLine(x1, y1, y2, col)
    VLine(x2, y1, y2, col)
}

func FillRect(x1, y1, x2, y2 int, col color.Color) {
    for x := x1; x <= x2; x++ {
        for y := y1; y <= y2; y++ {
            img.Set(x, y, col)
        }
    }
}

func bar_y(usage float32) int {
    return int(0.5 + 32*(1.0-usage))
}

func gen_img(ram_usage, phys_usage, cpu_usage float32) {
    img = image.NewRGBA(image.Rect(0, 0, 32, 32))

    Rect(0, 0, 32, 32, color.RGBA{0, 0, 0, 255})
    FillRect(1, 1, 31, 31, color.RGBA{255, 255, 255, 255})
    FillRect(1, bar_y(ram_usage), 10, 32, color.RGBA{255, 0, 0, 255})
    FillRect(11, bar_y(phys_usage), 21, 32, color.RGBA{0, 255, 0, 255})
    FillRect(21, bar_y(cpu_usage), 31, 32, color.RGBA{0, 0, 255, 255})
}

func gen_image_data(ram_usage, phys_usage, cpu_usage float32) []byte {
    gen_img(ram_usage, phys_usage, cpu_usage)

    var imageBuf bytes.Buffer
    ico.Encode(&imageBuf, img)
    return imageBuf.Bytes()
}
