package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Form struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}

func main() {
	r := gin.Default()
	r.StaticFS("/content", http.Dir("html/"))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/content")
	})
	r.POST("/analyze", func(c *gin.Context) {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)

		// Multipart form
		//var form Form
		//s := c.PostForm("files")
		form, _ := c.MultipartForm()
		files := form.File["files"]
		fmt.Println(form)
		//c.ShouldBind(&form)
		s := []string{"lisa4ros2-0.1a1/bin/lisa4ros2"}
		for _, file := range files {
			log.Println(file.Filename)
			//Get raw file bytes - no reader method
			// Upload to disk
			// `formFile` has io.reader method
			c.SaveUploadedFile(file, "analysis/ros-sources/"+timestamp+"/"+file.Filename)
			s = append(s, "analysis/ros-sources/"+timestamp+"/"+file.Filename)
		}
		s = append(s, "html/analysis/"+timestamp)
		fmt.Println(strings.Join(s, " "))
		cmd := exec.Command("/bin/sh", s...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		stdOut, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Error 1")
			c.JSON(http.StatusOK, gin.H{
				"error": true,
			})
			return
		}
		if err = cmd.Start(); err != nil {
			fmt.Println("Error 2")
			c.JSON(http.StatusOK, gin.H{
				"error": true,
			})
			return
		}
		bytes, err := io.ReadAll(stdOut)
		if err != nil {
			fmt.Println("Error 3")
			c.JSON(http.StatusOK, gin.H{
				"error": true,
			})
			return
		}
		if err = cmd.Wait(); err != nil {
			fmt.Println("Error 4")
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Printf("Exit code is %d\n", exitError.ExitCode())
				c.JSON(http.StatusOK, gin.H{
					"error": true,
				})
				return
			}
		}
		fmt.Println(string(bytes))
		if err != nil {
			fmt.Printf("error %s", err)
			c.JSON(http.StatusOK, gin.H{
				"url":   "analysis/" + timestamp + "/" + "report.html",
				"error": false,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"url":            "analysis/" + timestamp + "/" + "report.html",
			"analysisStatus": string(bytes),
		})
	})
	r.Run(":9888")
}
