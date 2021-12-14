//go:build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func controllCpuPriorityWindows() {
	fmt.Println("Set Priority: Below normal,", "you can feel free in work!!!")
	// You can chang speed archiving: BELOW_NORMAL_PRIORITY_CLASS, NORMAL_PRIORITY_CLASS, ABOVE_NORMAL_PRIORITY_CLASS HIGH_PRIORITY_CLASS
	windows.SetPriorityClass(windows.Handle(windows.CurrentProcess()), windows.BELOW_NORMAL_PRIORITY_CLASS)
}
