//go:build linux

package main

import (
	"golang.org/x/sys/unix"
)

func setCpuPriorityLinux(priority int) {
	// var priorityText string
	// if priority <= 20 && priority > 10 {
	// 	priorityText = "Very low"
	// } else if priority <= 10 && priority > 0 {
	// 	priorityText = "Low"
	// } else if priority <= -10 && priority >= -20 {
	// 	priorityText = "Hight"
	// } else {
	// 	priorityText = "Normal"
	// 	priority = 0
	// }
	// fmt.Println("Set Priority:", priorityText)
	unix.Setpriority(unix.PRIO_PROCESS, 0, priority)
}
