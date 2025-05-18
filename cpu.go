package main

import (
    "github.com/shirou/gopsutil/v3/cpu"
)

func get_cpu_usage() float32 {
    percent, err := cpu.Percent(0, false)
    if err != nil || len(percent) == 0 {
        return 0.0
    }
    return float32(percent[0]) / 100.0
}
