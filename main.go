package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.GET("/", func(ctx *gin.Context) {
		response := fmt.Sprintf("Hello, world!\nYour address is %v\nYour User-Agent is %v\n", ctx.Request.RemoteAddr, ctx.Request.Header["User-Agent"])
		ctx.String(http.StatusOK, response)
	})
	err := server.Run("0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
}
