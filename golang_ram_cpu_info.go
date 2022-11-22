package main

import (
    "fmt"
    "io"
    "time"

    "github.com/getlantern/systray"
    "github.com/getlantern/systray/example/icon"

    "github.com/shirou/gopsutil/v3/mem"

    "github.com/biessek/golang-ico"
    "image"
    "image/color"
    "os"
)

//     "github.com/skratchdot/open-golang/open"

func main() {
    fmt.Println("Hello, World!")

    onExit := func() {
        now := time.Now()
        // ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
        fmt.Println("%v", now)
    }

    systray.Run(onReady, onExit)
}

var img = image.NewRGBA(image.Rect(0, 0, 32, 32))
var col color.Color

// HLine draws a horizontal line
func HLine(x1, y, x2 int) {
    for ; x1 <= x2; x1++ {
        img.Set(x1, y, col)
    }
}

// VLine draws a veritcal line
func VLine(x, y1, y2 int) {
    for ; y1 <= y2; y1++ {
        img.Set(x, y1, col)
    }
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(x1, y1, x2, y2 int) {
    HLine(x1, y1, x2)
    HLine(x1, y2, x2)
    VLine(x1, y1, y2)
    VLine(x2, y1, y2)
}

func FillRect(x1, y1, x2, y2 int, col color.Color) {
    for x := x1; x <= x2; x++ {
        for y := y1; y <= y2; y++ {
            img.Set(x, y, col)
        }
    }
}

//https://stackoverflow.com/a/70115101/2624911
func ReadFileAndReturnByteArray(extractedFilePath string) ([]byte, error) {
    file, err := os.Open(extractedFilePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    return io.ReadAll(file)
}

func readIconData() []byte {
    var data, _ = ReadFileAndReturnByteArray("draw.ico")
    return data
}

func update_image(ram_usage, phys_usage float32) {
    img = image.NewRGBA(image.Rect(0, 0, 32, 32))
    // col = color.RGBA{255, 0, 0, 255} // Red
    // HLine(10, 20, 80)
    FillRect(0, int(0.5+32*(1.0-ram_usage)), 16, 32, color.RGBA{0, 255, 0, 255})
    FillRect(16, int(0.5+32*(1.0-phys_usage)), 32, 32, color.RGBA{255, 0, 0, 255})

    f, err := os.Create("draw.ico")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    ico.Encode(f, img)
}

func toGB(b uint64) float32 {
    return float32(b) / (1024.0 * 1024 * 1024)
}

func mem_str_one(used uint64, total uint64) string {
    return fmt.Sprintf("%.1f/%.1fGB %.0f%% %.1f GB Free", toGB(used), toGB(total), 100.0*float32(used)/float32(total), toGB(total-used))
}

func mem_str() (str string, ram_usage float32, phys_usage float32) {
    ram, _ := mem.SwapMemory()
    phys_ram, _ := mem.VirtualMemory()

    str = fmt.Sprintf("RAM : %s\nPhys: %s", mem_str_one(ram.Used, ram.Total), mem_str_one(phys_ram.Used, phys_ram.Total))
    ram_usage = float32(ram.Used) / float32(ram.Total)
    phys_usage = float32(phys_ram.Used) / float32(phys_ram.Total)
    return
}

func update() {
    var str, ram_usage, phys_usage = mem_str()
    systray.SetTooltip(str)
    update_image(ram_usage, phys_usage)
    systray.SetIcon(readIconData())
}

func printstuff() {
    var str, _, _ = mem_str()
    fmt.Println(str)
}

func onReady() {
    systray.SetTemplateIcon(icon.Data, icon.Data)
    systray.SetTitle("Awesome App")
    systray.SetTooltip("Lantern")
    mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
    go func() {
        <-mQuitOrig.ClickedCh
        fmt.Println("Requesting quit")
        systray.Quit()
        fmt.Println("Finished quitting")
    }()

    uptimeTicker := time.NewTicker(1 * time.Second)
    dateTicker := time.NewTicker(3 * time.Second)

    for {
        select {
        case <-uptimeTicker.C:
            go update()
        case <-dateTicker.C:
            go printstuff()
        }
    }
}
