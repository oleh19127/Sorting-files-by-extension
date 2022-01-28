//go:build linux

package main

import (
	"golang.org/x/sys/unix"
)

func setCpuPriorityLinux(priority int) {
	unix.Setpriority(unix.PRIO_PROCESS, 0, priority)
}
