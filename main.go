package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
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

func DisplayHeader(header http.Header) string {
	builder := strings.Builder{}
	keys := make([]string, 0)
	for key := range header {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%v: %v\n", key, strings.Join(header[key], ", ")))
	}
	return builder.String()
}

func main() {
	var port uint
	flag.UintVar(&port, "p", 8000, "port to listen")
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Any("/*any", func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		response := fmt.Sprintf("Welcome to http-echo!\n[Client: %v]\n[Url: %v]\n[Method: %v]\n\n[Header]\n%s",
			c.Request.RemoteAddr, c.Request.URL, c.Request.Method, DisplayHeader(c.Request.Header))
		if len(body) > 0 {
			response += fmt.Sprintf("\n[Body]\n%v", string(body))
		}
		c.String(http.StatusOK, response)
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
