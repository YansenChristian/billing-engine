package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"loan-payment/cmd/webservice/handlers"
	"loan-payment/cmd/webservice/middlewares"
	"loan-payment/configs"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	srv *http.Server
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)

	configs.Init("webservice")

	stopFunc := StartHttpServer()
	// will wait until terminate signal or interrupt happened
	for {
		<-c
		log.Println("terminate service")
		stopFunc()
		os.Exit(0)
	}

}

func StartHttpServer() (stopFunc func()) {
	const rootPath = "billing-engine/api"
	conf := configs.Get()
	r := gin.Default()

	r.Use(middlewares.WithRequestId())

	v1 := r.Group(rootPath+"/v1", timeout.New(
		timeout.WithTimeout(10*time.Second),
	))
	v1.POST("/", handlers.Handle)

	srv = &http.Server{
		Addr:    ":" + conf.HttpPort,
		Handler: r,
	}

	go func() {
		log.Printf("running %s in web service mode...\n", conf.AppName)
		log.Println("starting web, listening on", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("failed starting web on", srv.Addr, err)
		}
	}()

	return func() {
		GracefulStop()
	}
}

func GracefulStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("shuting down web on", srv.Addr)
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalln("failed shutdown server", err)
	}
	log.Println("web gracefully stopped")
}
