package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/k4zb3k/project/utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

var (
	Info      *log.Logger
	Error     *log.Logger
	Warn      *log.Logger
	Debug     *log.Logger
	MachineID *log.Logger
)

func Init() {
	fileInfo, err := os.OpenFile(utils.AppSettings.AppParams.LogInfo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	fileError, err := os.OpenFile(utils.AppSettings.AppParams.LogError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	fileWarn, err := os.OpenFile(utils.AppSettings.AppParams.LogWarning, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	fileDebug, err := os.OpenFile(utils.AppSettings.AppParams.LogDebug, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	fileMachineID, err := os.OpenFile(utils.AppSettings.AppParams.LogDebug, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		return
	}

	Info = log.New(fileInfo, "", log.Lmicroseconds)
	Error = log.New(fileError, "", log.Ldate|log.Lmicroseconds)
	Warn = log.New(fileWarn, "", log.Ldate|log.Lmicroseconds)
	Debug = log.New(fileDebug, "", log.Ldate|log.Lmicroseconds)
	MachineID = log.New(fileMachineID, "MachineHWID = ", 0)

	lumberLogInfo := &lumberjack.Logger{
		Filename:   utils.AppSettings.AppParams.LogInfo,
		MaxSize:    utils.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: utils.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     utils.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   utils.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	lumberLogError := &lumberjack.Logger{
		Filename:   utils.AppSettings.AppParams.LogError,
		MaxSize:    utils.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: utils.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     utils.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   utils.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	lumberLogWarn := &lumberjack.Logger{
		Filename:   utils.AppSettings.AppParams.LogWarning,
		MaxSize:    utils.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: utils.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     utils.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   utils.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	lumberLogDebug := &lumberjack.Logger{
		Filename:   utils.AppSettings.AppParams.LogDebug,
		MaxSize:    utils.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: utils.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     utils.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   utils.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	lumberLogMachineID := &lumberjack.Logger{
		Filename:   utils.AppSettings.AppParams.LogMachineHWID,
		MaxSize:    utils.AppSettings.AppParams.LogMaxSize, // megabytes
		MaxBackups: utils.AppSettings.AppParams.LogMaxBackups,
		MaxAge:     utils.AppSettings.AppParams.LogMaxAge,   //days
		Compress:   utils.AppSettings.AppParams.LogCompress, // disabled by default
		LocalTime:  true,
	}

	gin.DefaultWriter = io.MultiWriter(os.Stdout, lumberLogInfo)

	Info.SetOutput(gin.DefaultWriter)
	Error.SetOutput(lumberLogError)
	Warn.SetOutput(lumberLogWarn)
	Debug.SetOutput(lumberLogDebug)
	MachineID.SetOutput(lumberLogMachineID)
}
