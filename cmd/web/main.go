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
	config webconfig

	
	showHelp bool

	
)

func init() {
	config.ConfigFileName="config.json"
	config.PersonalPinFileName="personalcodes.csv"
	config.AccessLogFileName="accesslog.csv"
	config.ListenPort="0.0.0.0:8080"

	config.WebUsers=map[string]string{"webuser":"password"}
	config.ApiUsers=map[string]string{"apiuser":"password"}

	flag.StringVar(&configFileName, "config", "webconfig.json", "webconfig json")
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.Parse()
}

func setupRouter(cm codemanager.CodeManager) *gin.Engine {
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

		//r.LoadHTMLGlob("../../static/templates" + string(os.PathSeparator) + "*")
		//r.Static("/static", "../../static/templates"+string(os.PathSeparator))
	}
	
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/onetimepin")
	})

	userAuth := r.Group("/", gin.BasicAuth(gin.Accounts(config.WebUsers)))
	apiAuth := r.Group("/", gin.BasicAuth(gin.Accounts(config.ApiUsers)))

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
		accessLogs, err := getAccessLogs(config.AccessLogFileName)
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
	codeManagerStore, err := store.LoadConfiguration(config.ConfigFileName, config.PersonalPinFileName)
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

	var err error
	if config, err = store.LoadOrInitJsonConfiguration(configFileName,config); err!=nil {
		fmt.Println(err)
		return
	}

	cm, err:= loadCodeManager()
	if err!=nil {
		fmt.Println(err)
		return
	}

	r := setupRouter(cm)
	r.Run(config.ListenPort)
}

