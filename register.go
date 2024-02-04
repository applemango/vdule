package main

import "vdule/schedule"

func main() {
	err := schedule.RegisterHololiveSchedule()
	if err != nil {
		println(err.Error())
	}
}
