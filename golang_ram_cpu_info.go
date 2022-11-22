package main

import (
    "fmt"
    "time"

    "github.com/getlantern/systray"

    "github.com/shirou/gopsutil/v3/mem"
)

func main() {
    fmt.Println("Hello, World!")

    onExit := func() {
        now := time.Now()
        // ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
        fmt.Println("%v", now)
    }

    systray.Run(onReady, onExit)
}

func toGB(b uint64) float32 {
    return float32(b) / (1024.0 * 1024 * 1024)
}

func mem_str_one(used uint64, total uint64) string {
    return fmt.Sprintf("%.0f%% %.1f GB Free", 100.0*float32(used)/float32(total), toGB(total-used))
    // return fmt.Sprintf("%.1f/%.1fGB %.0f%% %.1f GB Free", toGB(used), toGB(total), 100.0*float32(used)/float32(total), toGB(total-used))
}

func mem_str() (str string, ram_usage float32, phys_usage float32, cpu_usage float32) {
    ram, _ := mem.SwapMemory()
    phys_ram, _ := mem.VirtualMemory()

    cpu_usage = get_cpu_usage()

    str = fmt.Sprintf("RAM : %s\nPhys: %s\nCPU : %.0f%%", mem_str_one(ram.Used, ram.Total), mem_str_one(phys_ram.Used, phys_ram.Total), 100.0*cpu_usage)
    ram_usage = float32(ram.Used) / float32(ram.Total)
    phys_usage = float32(phys_ram.Used) / float32(phys_ram.Total)
    return
}

func update() {
    var str, ram_usage, phys_usage, cpu_usage = mem_str()
    systray.SetTooltip(str)
    systray.SetIcon(gen_image_data(ram_usage, phys_usage, cpu_usage))
    fmt.Println(str)
}

func printstuff() {
    var str, _, _, _ = mem_str()
    fmt.Println(str)
}

func onReady() {
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
    //dateTicker := time.NewTicker(3 * time.Second)

    for {
        select {
        case <-uptimeTicker.C:
            go update()
            // case <-dateTicker.C:
            //     go printstuff()
        }
    }
}
