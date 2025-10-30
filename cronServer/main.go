package main

import (
	"log"

	"cronServer/job"
)

func main() {
	log.Println("start cron server")
	job.StartCron() // 启动 cron 任务
	select {}       // 阻塞主线程
}
