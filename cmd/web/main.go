package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/themulle/chronokeyaccess/internal/store"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
	"github.com/themulle/chronokeyaccess/static"
)

var (
	configFileName string
	personalPinFileName string
	accessLogFile string
	listenPort string
	basePath string

	apiUser string
	webUser string
	
	showHelp bool

	
)

func init() {
	flag.StringVar(&configFileName, "c", "config.json", "config file name")
	flag.StringVar(&personalPinFileName, "p", "personalcodes.csv", "personal pin csv file name")
	flag.StringVar(&accessLogFile, "accesslog", "accesslog.csv", "accesslog file")
	flag.StringVar(&webUser, "webuser", "webuser:password", "basic auth file")
	flag.StringVar(&apiUser, "apiuser", "api:password", "basic auth file")
	flag.StringVar(&listenPort,"port", "0.0.0.0:8080", "listen port")
	flag.StringVar(&basePath, "basePath", "../../", "base path")
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.Parse()
}

func setupRouter(baseDir string, cm codemanager.CodeManager) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	funcMap:=template.FuncMap{
		"formatAsPin": func(pin uint) string { return fmt.Sprintf("%05d", pin) },
		"formatAsUTF8Pin": func(pin uint) string {
			var emojiPin strings.Builder
			for _, digit := range fmt.Sprintf("%05d*", pin) {
				emojiPin.WriteString(string(digit) + "\u20E3")
			}
			return emojiPin.String()
		},
	}

	{
		tmpl := template.Must(template.New("").Funcs(funcMap).ParseFS(static.StaticFs, "templates/*.tmpl"))
    	r.SetHTMLTemplate(tmpl)

    	staticFS, _ := fs.Sub(static.StaticFs, "frontend")
    	r.StaticFS("/static", http.FS(staticFS))
	}
	
	//r.LoadHTMLGlob(filepath.Join(baseDir, "/templates") + string(os.PathSeparator) + "*")
	//r.Static("/static", filepath.Join(baseDir, "/frontend")+string(os.PathSeparator))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/onetimepin")
	})

	userAuth := r.Group("/")
	if parts:=strings.Split(webUser,":"); len(parts)==2 {
		userAuth = r.Group("/", gin.BasicAuth(gin.Accounts{
			parts[0]: parts[1],
		}))
	}

	apiAuth := r.Group("/")
	if parts:=strings.Split(apiUser,":"); len(parts)==2 {
		apiAuth = r.Group("/", gin.BasicAuth(gin.Accounts{
			parts[0]: parts[1],
		}))
	}

	userAuth.GET("/seriespin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "seriespin.tmpl", gin.H{
			"min": time.Now().Format("2006-01-02T15:04:05"),
			"max": time.Now().Add(time.Hour * 24 * 180).Format("2006-01-02T15:04:05"),
		})
	})

	userAuth.GET("/onetimepin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "onetimepin.tmpl", gin.H{
			"min": time.Now().Format("2006-01-02T15:04:05"),
			"max": time.Now().Add(time.Hour * 24 * 180).Format("2006-01-02T15:04:05"),
		})
	})
	userAuth.GET("/personalpin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "personalpin.tmpl", gin.H{
			"min": time.Now().Format("2006-01-02T15:04:05"),
			"max": time.Now().Add(time.Hour * 24 * 180).Format("2006-01-02T15:04:05"),
		})
	})
	userAuth.GET("/accesslog", func(c *gin.Context) {
		accessLogs, err := getAccessLogs()
		if err != nil {
			c.Error(err)
			return
		}
		slices.Reverse(accessLogs) // DESC
		c.HTML(http.StatusOK, "accesslog.tmpl", gin.H{
			"accessLogs": accessLogs,
		})
	})

	userAuth.GET("/codestemplate", func(c *gin.Context) {
		var cr CodeRequest
		if err := c.ShouldBind(&cr); err != nil {
			c.Error(err)
			return
		}
		entranceCodes, err := getCodes(cr, cm)
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

	apiAuth.GET("/codes", func(c *gin.Context) {
		var cr CodeRequest
		if err := c.ShouldBind(&cr); err != nil {
			c.Error(err)
			return
		}
		entranceCodes, err := getCodes(cr, cm)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, entranceCodes)
	})

	userAuth.GET("/reload", func(c *gin.Context) {
		if tmp, err := loadCodeManager(); err!=nil {
			c.Error(err)
			return
		} else {
			cm=tmp
		}
		c.JSON(200,"OK")
	})

	return r
}

func loadCodeManager() (codemanager.CodeManager, error) {
	codeManagerStore, err := store.LoadConfiguration(configFileName, personalPinFileName, true)
	if err!=nil {
		return nil, err
	}
	return codemanager.InitFromStore(codeManagerStore)
}



func main() {
	if showHelp {
		flag.PrintDefaults()
		return
	}

	cm, err:= loadCodeManager()
	if err!=nil {
		fmt.Println(err)
		return
	}

	r := setupRouter(basePath, cm)
	r.Run(listenPort)
}
