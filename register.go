package main

import (
	"context"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	"vdule/schedule"
)

func hololive(ctx context.Context) (any, error) {
	err := schedule.RegisterHololiveSchedule()
	if err != nil {
		println(err.Error())
	}
	return nil, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched := quartz.NewStdScheduler()

	every15, _ := quartz.NewCronTrigger("0 */15 * ? * *")

	_ = sched.ScheduleJob(
		quartz.NewJobDetail(
			job.NewFunctionJob(hololive),
			quartz.NewJobKey("hololive"),
		),
		every15,
	)

	sched.Start(ctx)
	sched.Wait(ctx)
}
