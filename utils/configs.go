package utils

import (
	"github.com/k4zb3k/project/internal/models"
)

var AppSettings models.Settings

func PutAdditionalSettings() {
	AppSettings.AppParams.LogDebug = "/home/k4zb3k/Desktop/project/logs/debug.log"
	AppSettings.AppParams.LogInfo = "/home/k4zb3k/Desktop/project/logs/info.log"
	AppSettings.AppParams.LogWarning = "/home/k4zb3k/Desktop/project/logs/warning.log"
	AppSettings.AppParams.LogError = "/home/k4zb3k/Desktop/project/logs/error.log"

	AppSettings.AppParams.LogCompress = true
	AppSettings.AppParams.LogMaxSize = 10
	AppSettings.AppParams.LogMaxAge = 100
	AppSettings.AppParams.LogMaxBackups = 100
	AppSettings.AppParams.AppVersion = "1.0"
}
