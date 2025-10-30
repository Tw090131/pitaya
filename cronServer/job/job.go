package job

import (
	"log"

	"github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New(cron.WithSeconds())
	log.Println("StartCron")
	c.AddFunc("*/10 * * * * *", func() {
		log.Println("每10秒执行一次任务")
	})
	c.Start()
}
