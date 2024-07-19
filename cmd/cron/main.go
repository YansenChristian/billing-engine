package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"loan-payment/cmd/cron/handlers"
	"loan-payment/configs"

	"github.com/robfig/cron/v3"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)

	configs.Init("cron")

	stopFunc := StartCronServer()
	// will wait until terminate signal or interrupt happened
	for {
		<-c
		log.Println("terminate service")
		stopFunc()
		os.Exit(0)
	}
}

func StartCronServer() (stopFunc func()) {
	client := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

	schedule, err := parser.Parse(configs.Get().CronCheckLoanStatusSchedule)
	if err != nil {
		log.Fatalf("invalid cron specs, err:%v", err)
	}

	client.Schedule(schedule, cron.FuncJob(handlers.CheckLoanStatus()))

	// start cron
	log.Printf("running %s in cron  mode...\n", configs.Get().AppName)
	go client.Start()
	return func() {
		client.Stop()
	}
}
