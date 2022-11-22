package main

import (
    "fmt"
    // "io/ioutil"
    "time"

    "github.com/getlantern/systray"
    "github.com/getlantern/systray/example/icon"

    "github.com/shirou/gopsutil/v3/mem"
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

func toGB(b uint64) float32 {
    return float32(b) / (1024.0 * 1024 * 1024)
}

func mem_str_one(used uint64, total uint64) string {
    return fmt.Sprintf("%.1f/%.1fGB %.0f%% %.1f GB Free", toGB(used), toGB(total), 100.0*float32(used)/float32(total), toGB(total-used))
}

func mem_str() string {
    ram, _ := mem.SwapMemory()
    phys_ram, _ := mem.VirtualMemory()

    return fmt.Sprintf("RAM : %s\nPhys: %s", mem_str_one(ram.Used, ram.Total), mem_str_one(phys_ram.Used, phys_ram.Total))
}

func update() {
    systray.SetTooltip(mem_str())
}

func printstuff() {
    fmt.Println(mem_str())
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
