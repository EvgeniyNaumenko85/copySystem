package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"tasks_API/configs"
	"time"
)

// SetLogger Установка Logger-а
var (
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
	Debug *log.Logger
)

func Init() {
	fileInfo, err := os.OpenFile(configs.AppSettings.AppParams.LogInfo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	fileError, err := os.OpenFile(configs.AppSettings.AppParams.LogError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	fileWarn, err := os.OpenFile(configs.AppSettings.AppParams.LogWarning, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	fileDebug, err := os.OpenFile(configs.AppSettings.AppParams.LogDebug, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		return
	}

	Info = log.New(fileInfo, "", log.Ldate|log.Lmicroseconds)
	Error = log.New(fileError, "", log.Ldate|log.Lmicroseconds)
	Warn = log.New(fileWarn, "", log.Ldate|log.Lmicroseconds)
	Debug = log.New(fileDebug, "", log.Ldate|log.Lmicroseconds)

	lumberLogInfo := &lumberjack.Logger{
		Filename:   configs.AppSettings.AppParams.LogInfo,
		MaxSize:    configs.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: configs.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     configs.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   configs.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	lumberLogError := &lumberjack.Logger{
		Filename:   configs.AppSettings.AppParams.LogError,
		MaxSize:    configs.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: configs.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     configs.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   configs.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	lumberLogWarn := &lumberjack.Logger{
		Filename:   configs.AppSettings.AppParams.LogWarning,
		MaxSize:    configs.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: configs.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     configs.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   configs.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	lumberLogDebug := &lumberjack.Logger{
		Filename:   configs.AppSettings.AppParams.LogDebug,
		MaxSize:    configs.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: configs.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     configs.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   configs.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	gin.DefaultWriter = io.MultiWriter(os.Stdout, lumberLogInfo)

	Info.SetOutput(lumberLogError)
	//Info.SetOutput(gin.DefaultWriter)
	Error.SetOutput(lumberLogError)
	Warn.SetOutput(lumberLogWarn)
	Debug.SetOutput(lumberLogDebug)
}

// FormatLogs Форматирование логов
func FormatLogs(r *gin.Engine) {
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("[GIN] %s - [%s] \"%s %s %s %d %s \"%s\" %s\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
}
