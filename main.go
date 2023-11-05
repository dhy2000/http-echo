package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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

func Serve(server *http.Server, certFile string, keyFile string) error {
	if len(certFile) == 0 || len(keyFile) == 0 {
		log.Printf("Go visit: http://%v", server.Addr)
		return server.ListenAndServe()
	} else {
		log.Printf("Go visit: https://%v", server.Addr)
		return server.ListenAndServeTLS(certFile, keyFile)
	}
}

func main() {
	var address, name string
	var cert, key string
	flag.StringVar(&address, "a", "0.0.0.0:8000", "address to listen")
	flag.StringVar(&name, "n", "HTTP-Echo", "name of this service")
	flag.StringVar(&cert, "c", "", "path to cert file")
	flag.StringVar(&key, "k", "", "path to key file")
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Any("/*any", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		response := fmt.Sprintf("%v\n[Client: %v]\n[Url: %v]\n[Method: %v]\n\n[Header]\n%s",
			name, c.Request.RemoteAddr, c.Request.URL, c.Request.Method, DisplayHeader(c.Request.Header))
		if len(body) > 0 {
			response += fmt.Sprintf("\n[Body]\n%v", string(body))
		}
		c.String(http.StatusOK, response)
	})
	server := &http.Server{Addr: address, Handler: router}
	GracefulExit(func() {
		log.Printf("Welcome to %v, listening to %v", name, address)
		if err := Serve(server, cert, key); err != nil && err != http.ErrServerClosed {
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
