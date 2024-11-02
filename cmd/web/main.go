package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/themulle/chronokeyaccess/internal/store"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
	"github.com/themulle/chronokeyaccess/static"
)

const (
	defaultConfigFile    = "config.json"
	defaultPinFile      = "personalcodes.csv"
	defaultAccessLogFile = "accesslog.csv"
	defaultWebConfig = 	"webconfig.json"
	defaultListenPort   = "0.0.0.0:8080"
	defaultLogLevel     = "info"
)

var (
	configFileName      string
	webConfigFileName string
	config webconfig
	showHelp bool
	accessLogFile      string
	logLevelFlag       string
	personalPinFileName string
)

func init() {
	flag.StringVar(&configFileName, "c", defaultConfigFile, "config file name")
    flag.StringVar(&personalPinFileName, "p", defaultPinFile, "personal pin csv file name")
    flag.StringVar(&accessLogFile, "accesslog", defaultAccessLogFile, "accesslog file")
    flag.StringVar(&logLevelFlag, "loglevel", defaultLogLevel, "log level (debug, info, warn, error)")
	flag.StringVar(&webConfigFileName, "webconfig", defaultWebConfig, "webconfig json")

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

func setupLogger(logLevelFlag string) error {
    logLevels := map[string]slog.Level{
        "debug": slog.LevelDebug,
        "info":  slog.LevelInfo,
        "warn":  slog.LevelWarn,
        "error": slog.LevelError,
    }

    logLevel, ok := logLevels[logLevelFlag]
    if !ok {
        return fmt.Errorf("ungültiges Log-Level: %s. Verwenden Sie 'debug', 'info', 'warn' oder 'error'", logLevelFlag)
    }

    handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
        Level: logLevel,
    })
    logger := slog.New(handler)
    slog.SetDefault(logger)

    slog.Info("Logger erfolgreich konfiguriert", 
        "level", logLevelFlag,
        "handler", "JSONHandler")

    return nil
}



func main() {
	flag.Parse()
	if showHelp {
		flag.PrintDefaults()
		return
	}

	if err := setupLogger(logLevelFlag); err != nil {
        slog.Error("Fehler beim Einrichten des Loggers", "error", err)
        os.Exit(1)
    }

	var err error
	if config, err = store.LoadOrInitJsonConfiguration(configFileName,config); err!=nil {
		fmt.Println(err)
		return
	}

	cm, err := loadCodeManager()
    if err != nil {
        slog.Error("Fehler beim Laden des CodeManagers", "error", err)
        os.Exit(1)
    }


	r := setupRouter(cm)
    slog.Info("Webserver wird gestartet", "port", defaultListenPort)
    r.Run(defaultListenPort)
}

