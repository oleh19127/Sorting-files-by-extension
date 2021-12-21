//go:build windows

package main

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

func setCpuPriorityWindows(priority string) {
	// BELOW_NORMAL_PRIORITY_CLASS, NORMAL_PRIORITY_CLASS, ABOVE_NORMAL_PRIORITY_CLASS, HIGH_PRIORITY_CLASS
	var priorityClass uint32
	var priorityText string
	if strings.ToLower(priority) == "below normal" {
		priorityText = "Below normal"
		priorityClass = windows.BELOW_NORMAL_PRIORITY_CLASS
	} else if strings.ToLower(priority) == "above normal" {
		priorityText = "Above normal"
		priorityClass = windows.ABOVE_NORMAL_PRIORITY_CLASS
	} else if strings.ToLower(priority) == "hight" {
		priorityText = "Hight"
		priorityClass = windows.HIGH_PRIORITY_CLASS
	} else {
		priorityText = "Normal"
		priorityClass = windows.NORMAL_PRIORITY_CLASS
	}
	fmt.Println("Set Priority: ", priorityText)
	windows.SetPriorityClass(windows.Handle(windows.CurrentProcess()), priorityClass)
}
