// === images.go ===
package main

import (
    "bytes"
    "fmt"
    "image"
    "image/color"
    "runtime"
    "sync"
    "github.com/biessek/golang-ico"
)

// Pre-allocate a single shared image buffer
var globalImg = image.NewRGBA(image.Rect(0, 0, 32, 32))

// Pre-defined colors to avoid allocations
var (
    colorBlack   = color.RGBA{0, 0, 0, 255}
    colorWhite   = color.RGBA{255, 255, 255, 255}
    colorRed     = color.RGBA{255, 0, 0, 255}
    colorGreen   = color.RGBA{0, 255, 0, 255}
    colorBlue    = color.RGBA{0, 0, 255, 255}
)

// Create a buffer pool to reuse buffers and reduce allocations
var bufferPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

func gen_img(ram_usage, phys_usage, cpu_usage float32) *image.RGBA {
    img := globalImg

    // Clear the image to white
    for x := 0; x < 32; x++ {
        for y := 0; y < 32; y++ {
            img.Set(x, y, colorWhite)
        }
    }

    // Draw bars
    Rect(img, 0, 0, 32, 32, colorBlack)
    FillRect(img, 1, bar_y(ram_usage), 10, 32, colorRed)
    FillRect(img, 11, bar_y(phys_usage), 21, 32, colorGreen)
    FillRect(img, 21, bar_y(cpu_usage), 31, 32, colorBlue)

    return img
}

func HLine(img *image.RGBA, x1, y, x2 int, col color.Color) {
    for ; x1 <= x2; x1++ {
        img.Set(x1, y, col)
    }
}

func VLine(img *image.RGBA, x, y1, y2 int, col color.Color) {
    for ; y1 <= y2; y1++ {
        img.Set(x, y1, col)
    }
}

func Rect(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
    HLine(img, x1, y1, x2, col)
    HLine(img, x1, y2, x2, col)
    VLine(img, x1, y1, y2, col)
    VLine(img, x2, y1, y2, col)
}

func FillRect(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
    for x := x1; x <= x2; x++ {
        for y := y1; y <= y2; y++ {
            img.Set(x, y, col)
        }
    }
}

func bar_y(usage float32) int {
    return int(0.5 + 32*(1.0-usage))
}

// This function is more efficient since it reuses the buffer from the pool
func gen_image_data(ram_usage, phys_usage, cpu_usage float32) []byte {
    img := gen_img(ram_usage, phys_usage, cpu_usage)

    // Get a buffer from the pool
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    
    // Write to the buffer
    ico.Encode(buf, img)
    
    // Copy the data to avoid keeping references to the buffer
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    
    // Return the buffer to the pool
    bufferPool.Put(buf)

    return result
}

func printMemoryUsage(tag string) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("[MEM %s] Alloc = %.2f MB | TotalAlloc = %.2f MB | Sys = %.2f MB | NumGC = %v\n",
        tag,
        float64(m.Alloc)/1024.0/1024.0,
        float64(m.TotalAlloc)/1024.0/1024.0,
        float64(m.Sys)/1024.0/1024.0,
        m.NumGC)
}

func memoryString() string {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    return fmt.Sprintf("Mem: %.1f MB", float64(m.Alloc)/1024.0/1024.0)
}