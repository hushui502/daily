package logging

import (
	"awesomeProject2/ginlearn/pkg/file"
	"awesomeProject2/ginlearn/pkg/setting"
	"fmt"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt = "log"
	TimeFormat = "19960909"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,)
}

func openLogFile(fileName, filePath string) (*os.File, error) {
	return file.MustOpen(fileName, filePath)
}