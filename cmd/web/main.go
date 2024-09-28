package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(baseDir string) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"formatAsPin": func(pin uint) string { return fmt.Sprintf("%04d", pin) },
	})

	r.LoadHTMLGlob(filepath.Join(baseDir, "/templates") + string(os.PathSeparator) + "*")
	r.Static("/static", filepath.Join(baseDir, "/frontend")+string(os.PathSeparator))

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
	listenPort := flag.String("port", "0.0.0.0:8080", "listen port")
	basePath := flag.String("basePath", "../../", "base path")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if help != nil && *help {
		flag.PrintDefaults()
		return
	}

	r := setupRouter(*basePath)
	r.Run(*listenPort)
}
