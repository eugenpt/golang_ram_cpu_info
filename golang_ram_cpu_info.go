// === golang_ram_cpu_info.go ===
package main

import (
    "bytes"
    "fmt"
    "time"

    "github.com/biessek/golang-ico"
    "github.com/getlantern/systray"
    "github.com/shirou/gopsutil/v3/mem"
)

// Smaller cache with LRU-like behavior for most frequently used icons
// We'll use a fixed size array and only store the most common icons
var iconCache map[[3]int][]byte
var lastIconKey [3]int // Keep track of last icon to avoid redundant updates

func main() {
    fmt.Println("Hello, World!")

    onExit := func() {
        now := time.Now()
        fmt.Println("%v", now)
    }

    systray.Run(onReady, onExit)
}

func toGB(b uint64) float32 {
    return float32(b) / (1024.0 * 1024 * 1024)
}

func mem_str_one(used uint64, total uint64) string {
    return fmt.Sprintf("%.0f%% %.1f GB Free", 100.0*float32(used)/float32(total), toGB(total-used))
}

func mem_str() (str string, ram_usage float32, phys_usage float32, cpu_usage float32) {
    ram, _ := mem.SwapMemory()
    phys_ram, _ := mem.VirtualMemory()
    cpu_usage = get_cpu_usage()

    str = fmt.Sprintf("RAM : %s\nPhys: %s\nCPU : %.0f%%\n%s", mem_str_one(ram.Used, ram.Total), mem_str_one(phys_ram.Used, phys_ram.Total), 100.0*cpu_usage, memoryString())
    ram_usage = float32(ram.Used) / float32(ram.Total)
    phys_usage = float32(phys_ram.Used) / float32(phys_ram.Total)
    return
}

func roundUsage(u float32) int {
    val := int(u*10 + 0.5)
    if val > 10 {
        val = 10
    }
    return val
}

// Generate a subset of icons instead of all possible combinations
func generateCommonIcons() {
    fmt.Println("Generating icon cache...")
    printMemoryUsage("BEFORE")
    
    // Only allocate space for the most commonly used icons
    // Typically system resources don't jump around wildly, so we can
    // generate icons for the most common usage patterns
    iconCache = make(map[[3]int][]byte, 125) // Much smaller allocation: 5×5×5 instead of 11×11×11
    
    // Generate icons for the range of 0-100% in 25% increments
    // This covers most common usage scenarios with far fewer icons
    for r := 0; r <= 10; r += 2 {
        for p := 0; p <= 10; p += 2 {
            for c := 0; c <= 10; c += 2 {
                ram := float32(r) / 10.0
                phys := float32(p) / 10.0
                cpu := float32(c) / 10.0

                img := gen_img(ram, phys, cpu)
                buf := new(bytes.Buffer)
                ico.Encode(buf, img)
                iconCache[[3]int{r, p, c}] = buf.Bytes()
            }
        }
    }
    
    fmt.Println("Finished generating icons.")
    printMemoryUsage("AFTER")
}

// Get nearest available icon or generate on-demand if needed
func getIcon(r, p, c int) []byte {
    key := [3]int{r, p, c}
    
    // Return the icon if it's in the cache
    if icon, exists := iconCache[key]; exists {
        return icon
    }
    
    // If not in cache, find nearest pre-generated icon
    // Round to nearest even number (since we generate at 2-step intervals)
    nearestR := (r + 1) / 2 * 2
    nearestP := (p + 1) / 2 * 2
    nearestC := (c + 1) / 2 * 2
    
    if nearestR > 10 {
        nearestR = 10
    }
    if nearestP > 10 {
        nearestP = 10
    }
    if nearestC > 10 {
        nearestC = 10
    }
    
    nearestKey := [3]int{nearestR, nearestP, nearestC}
    
    // Check if nearest icon exists
    if icon, exists := iconCache[nearestKey]; exists {
        return icon
    }
    
    // Generate the icon on-demand if not in cache
    ram := float32(r) / 10.0
    phys := float32(p) / 10.0
    cpu := float32(c) / 10.0

    img := gen_img(ram, phys, cpu)
    buf := new(bytes.Buffer)
    ico.Encode(buf, img)
    
    // Store it in the cache
    iconData := make([]byte, buf.Len())
    copy(iconData, buf.Bytes())
    iconCache[key] = iconData
    
    return iconData
}

func update() {
    str, ram, phys, cpu := mem_str()
    r := roundUsage(ram)
    p := roundUsage(phys)
    c := roundUsage(cpu)
    
    // Avoid redundant icon updates for the same values
    newKey := [3]int{r, p, c}
    if newKey != lastIconKey {
        systray.SetIcon(getIcon(r, p, c))
        lastIconKey = newKey
    }
    
    // Always update tooltip
    systray.SetTooltip(str)
    fmt.Println(str)
}

func onReady() {
    systray.SetTitle("Awesome App")
    systray.SetTooltip("Lantern")

    // Generate a smaller set of icons
    generateCommonIcons()

    mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
    go func() {
        <-mQuitOrig.ClickedCh
        fmt.Println("Requesting quit")
        systray.Quit()
        fmt.Println("Finished quitting")
    }()

    uptimeTicker := time.NewTicker(1 * time.Second)
    for {
        select {
        case <-uptimeTicker.C:
            update()
        }
    }
}