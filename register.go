package main

import (
	"context"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	"vdule/vtuber/schedule"
)

func hololive(ctx context.Context) (any, error) {
	err := schedule.RegisterHololiveSchedule()
	if err != nil {
		println(err.Error())
	}
	return nil, nil
}

func youtube(ctx context.Context) (any, error) {
	err := schedule.RegisterYoutuberSchedule()
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
	every30, _ := quartz.NewCronTrigger("0 */30 * ? * *")

	_ = sched.ScheduleJob(
		quartz.NewJobDetail(
			job.NewFunctionJob(hololive),
			quartz.NewJobKey("hololive"),
		),
		every15,
	)

	_ = sched.ScheduleJob(
		quartz.NewJobDetail(
			job.NewFunctionJob(youtube),
			quartz.NewJobKey("youtube"),
		),
		every30,
	)

	sched.Start(ctx)
	sched.Wait(ctx)
}
