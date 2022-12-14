package main

import (
    "github.com/shirou/gopsutil/v3/cpu"
)

var gCPUstat cpu.TimesStat

func getCPUstat() cpu.TimesStat {
    cpu_stats, _ := cpu.Times(false)
    return cpu_stats[len(cpu_stats)-1]
}

func get_delta_CPUstats() cpu.TimesStat {
    cpustat := getCPUstat()
    var r cpu.TimesStat

    // fmt.Printf("User:\nold : %v\nnew : %v\ndelta %v\n\n", gCPUstat.User, cpustat.User, cpustat.User-gCPUstat.User)
    // fmt.Printf("System:\nold : %v\nnew : %v\ndelta %v\n\n", gCPUstat.System, cpustat.System, cpustat.System-gCPUstat.System)
    // fmt.Printf("Idle:\nold : %v\nnew : %v\ndelta %v\n\n", gCPUstat.Idle, cpustat.Idle, cpustat.Idle-gCPUstat.Idle)

    r.User = cpustat.User - gCPUstat.User
    r.System = cpustat.System - gCPUstat.System
    r.Idle = cpustat.Idle - gCPUstat.Idle

    gCPUstat = cpustat
    return r
}

func get_cpu_usage() float32 {
    cpu_stat := get_delta_CPUstats()

    return float32((cpu_stat.System + cpu_stat.User) / (cpu_stat.Idle + cpu_stat.System + cpu_stat.User))
}
