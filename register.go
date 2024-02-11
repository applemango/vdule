package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	"vdule/vtuber/schedule"
	youtube2 "vdule/vtuber/youtube"
)

func RegisterHololiveSchedule(ctx context.Context) (any, error) {
	err := schedule.RegisterHololiveSchedule()
	if err != nil {
		println(err.Error())
	}
	return nil, nil
}

func RegisterYoutubeSchedule(ctx context.Context) (any, error) {
	err := schedule.RegisterYoutuberSchedule()
	if err != nil {
		println(err.Error())
	}
	return nil, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("failed load config")
	}
	youtube2.ResetYoutube()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _ = RegisterHololiveSchedule(ctx)
	_, _ = RegisterYoutubeSchedule(ctx)

	sched := quartz.NewStdScheduler()

	every15, _ := quartz.NewCronTrigger("0 */15 * ? * *")
	every30, _ := quartz.NewCronTrigger("0 */30 * ? * *")

	_ = sched.ScheduleJob(
		quartz.NewJobDetail(
			job.NewFunctionJob(RegisterHololiveSchedule),
			quartz.NewJobKey("hololive"),
		),
		every15,
	)

	_ = sched.ScheduleJob(
		quartz.NewJobDetail(
			job.NewFunctionJob(RegisterYoutubeSchedule),
			quartz.NewJobKey("youtube"),
		),
		every30,
	)

	sched.Start(ctx)
	sched.Wait(ctx)
}
