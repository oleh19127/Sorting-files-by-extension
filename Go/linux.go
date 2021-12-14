//go:build linux

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func controllCpuPriorityLinux() {
	fmt.Println("Set Priority: Low,", "you can feel free in work!!!")
	unix.Setpriority(unix.PRIO_PROCESS, 0, 5)
}
