package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	blanket := "off"
	cmd := "/usr/sbin/fast-gpio"
	gpio := "0"
	args1 := []string{"set-output", gpio}
	args2 := []string{}
	if err := exec.Command(cmd, args1...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Electric Blanket",
		})
	})

	r.GET("/blanket", func(c *gin.Context) {
		// Need to fix this to read the actual real stautus from the hardware
		// fast-gpio read <gpio pin>
		c.JSON(200, gin.H{
			"message": blanket,
		})
	})
	r.POST("/blanket", func(c *gin.Context) {
		state := c.PostForm("state")
		if state == "on" {
			args2 = []string{"set", gpio, "1"}
			blanket = "on"
		} else {
			args2 = []string{"set", gpio, "0"}
			blanket = "off"
		}
		if err := exec.Command(cmd, args2...).Run(); err != nil {
			c.JSON(200, gin.H{
				"message": "failure",
			})
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})
	r.Run() // listen and server on 0.0.0.0:8080
}
