package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"formatAsPin": func(pin uint) string { return fmt.Sprintf("%04d", pin) },
	})

	r.LoadHTMLGlob("../../templates/*")
	r.Static("/static", "../../frontend/")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/onetimepin")
	})

	r.GET("/seriespin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "seriespin.tmpl", gin.H{
			"min": time.Now().Format("2006-01-02T15:04:05"),
			"max": time.Now().Add(time.Hour * 24 * 180).Format("2006-01-02T15:04:05"),
		})
	})

	r.GET("/onetimepin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "onetimepin.tmpl", gin.H{
			"min": time.Now().Format("2006-01-02T15:04:05"),
			"max": time.Now().Add(time.Hour * 24 * 180).Format("2006-01-02T15:04:05"),
		})
	})
	r.GET("/personalpin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "personalpin.tmpl", gin.H{
			"min": time.Now().Format("2006-01-02T15:04:05"),
			"max": time.Now().Add(time.Hour * 24 * 180).Format("2006-01-02T15:04:05"),
		})
	})
	r.GET("/accesslog", func(c *gin.Context) {
		accessLogs, err := getAccessLogs()
		if err != nil {
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "accesslog.tmpl", gin.H{
			"accessLogs": accessLogs,
		})
	})

	r.GET("/codestemplate", func(c *gin.Context) {
		var cr CodeRequest
		if err := c.ShouldBind(&cr); err != nil {
			c.Error(err)
			return
		}
		entranceCodes, err := getCodes(cr)
		if err != nil {
			c.Error(err)
			return
		}

		entranceTimes := map[string]string{}
		for _, ec := range entranceCodes {
			entranceTimes[ec.Start.Format("15:04")] = fmt.Sprintf("%s - %s", ec.Start.Format("15:04"), ec.Stop.Format("15:04"))
		}

		c.HTML(http.StatusOK, "codelist.tmpl", gin.H{
			"entranceCodes": entranceCodes,
			"entranceTimes": entranceTimes,
		})
	})

	r.GET("/codes", func(c *gin.Context) {
		var cr CodeRequest
		if err := c.ShouldBind(&cr); err != nil {
			c.Error(err)
			return
		}
		entranceCodes, err := getCodes(cr)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, entranceCodes)
	})

	return r
}

func main() {

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run("127.0.0.1:8080")
}
