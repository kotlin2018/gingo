package api

var coreConfig = `package core

import (
	"fmt"
	"{{.Appname}}/global"
	_ "{{.Appname}}/packfile"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const defaultConfigFile = "config.yaml"

func init() {
	v := viper.New()
	v.SetConfigFile(defaultConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	global.GVA_VP = v
}
`

var coreLog = `package core

import (
	"fmt"
	"{{.Appname}}/config"
	"{{.Appname}}/global"
	"{{.Appname}}/utils"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	oplogging "github.com/op/go-logging"
	"io"
	"os"
	"strings"
	"time"
)

const (
	logDir      = "log"
	logSoftLink = "latest_log"
	module      = "gin-vue-admin"
)

func init() {
	c := global.GVA_CONFIG.Log
	if c.Prefix == "" {
		_ = fmt.Errorf("logger prefix not found")
	}
	logger := oplogging.MustGetLogger(module)
	var backends []oplogging.Backend
	registerStdout(c, &backends)
	if fileWriter := registerFile(c, &backends); fileWriter != nil {
		gin.DefaultWriter = io.MultiWriter(fileWriter, os.Stdout)
	}
	oplogging.SetBackend(backends...)
	global.GVA_LOG = logger
}

func registerStdout(c config.Log, backends *[]oplogging.Backend) {
	if c.Stdout != "" {
		level, err := oplogging.LogLevel(c.Stdout)
		if err != nil {
			fmt.Println(err)
		}
		*backends = append(*backends, createBackend(os.Stdout, c, level))
	}
}

func registerFile(c config.Log, backends *[]oplogging.Backend) io.Writer {
	if c.File != "" {
		if ok, _ := utils.PathExists(logDir); !ok {
			// directory not exist
			fmt.Println("create log directory")
			_ = os.Mkdir(logDir, os.ModePerm)
		}
		fileWriter, err := rotatelogs.New(
			logDir+string(os.PathSeparator)+"%Y-%m-%d-%H-%M.log",
			// generate soft link, point to latest log file
			rotatelogs.WithLinkName(logSoftLink),
			// maximum time to save log files
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// time period of log file switching
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			fmt.Println(err)
		}
		level, err := oplogging.LogLevel(c.File)
		if err != nil {
			fmt.Println(err)
		}
		*backends = append(*backends, createBackend(fileWriter, c, level))

		return fileWriter
	}
	return nil
}

func createBackend(w io.Writer, c config.Log, level oplogging.Level) oplogging.Backend {
	backend := oplogging.NewLogBackend(w, c.Prefix, 0)
	stdoutWriter := false
	if w == os.Stdout {
		stdoutWriter = true
	}
	format := getLogFormatter(c, stdoutWriter)
	backendLeveled := oplogging.AddModuleLevel(oplogging.NewBackendFormatter(backend, format))
	backendLeveled.SetLevel(level, module)
	return backendLeveled
}

func getLogFormatter(c config.Log, stdoutWriter bool) oplogging.Formatter {
	pattern := defaultFormatter
	if !stdoutWriter {
		// Color is only required for console output
		// Other writers don't need %{color} tag
		pattern = strings.Replace(pattern, "%{color:bold}", "", -1)
		pattern = strings.Replace(pattern, "%{color:reset}", "", -1)
	}
	if !c.LogFile {
		// Remove %{logfile} tag
		pattern = strings.Replace(pattern, "%{longfile}", "", -1)
	}
	return oplogging.MustStringFormatter(pattern)
}

`

var coreServe = `package core

import (
	"fmt"
	"{{.Appname}}/global"
	"{{.Appname}}/initialize"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint {
		// 初始化redis服务
		initialize.Redis()
	}
	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Debug("server run success on ", address)
	global.GVA_LOG.Error(s.ListenAndServe())
}
`

var coreServerOther = `package core

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"time"
)

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
`

var coreServerWin = `package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
`
