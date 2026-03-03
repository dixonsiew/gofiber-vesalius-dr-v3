package cron

import (
	"fmt"
	"time"
	"vesaliusdr/utils"
	"vesaliusdr/ws"

	"github.com/go-co-op/gocron/v2"
)

var cronTask gocron.Scheduler

func catchPanic(funcName string) {
    if err := recover(); err != nil {
        utils.LogInfo(fmt.Sprintf("recovered from panic -%s", funcName))
    }
}

func Setup() {
    s, _ := gocron.NewScheduler()
    _, _ = s.NewJob(
        gocron.DurationJob(5 * time.Minute),
        gocron.NewTask(func() {
            defer catchPanic("GetToken-GetTokenDoc")
            ws.ProcessDoctorToReviewData()
        }),
    )
    cronTask = s
    s.Start()
}

func Shutdown() {
    cronTask.Shutdown()
}