package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulExit(run func(), exit func()) {
	go run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1)
	<-ch
	exit()
}

func main() {
	var port uint
	flag.UintVar(&port, "p", 8000, "port to listen")
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		response := fmt.Sprintf("Hello, world!\nYour address is %v\nYour User-Agent is %v\n", ctx.Request.RemoteAddr, ctx.Request.Header["User-Agent"])
		ctx.String(http.StatusOK, response)
	})
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: router}
	GracefulExit(func() {
		log.Printf("Listening at port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	})
}
