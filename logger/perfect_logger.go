package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/afero"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// 日志文件的时间格式
	LogFileNameTimeFormat = "2006-01-02T15-04-05.000"

	// 日志文件大小限制
	LogMaxSize = 25

	// 日志文件的备份个数
	LogMaxBackups = 20

	// 日志文件的最大天数
	LogMaxAge = 7
)

var (
	// 缓存日志
	LumberLogger *Logger
	Log          = NewZeroLog("", "log", true, 3)
)

func CreateLogger(logName string) zerolog.Logger {
	return NewZeroLog("", logName, true, -1)
}

func GetCurrentExecDir() (dir string, err error) {
	dir, err1 := os.Getwd()
	if err1 != nil {
		fmt.Println(err1)
	}
	return strings.Replace(dir, "\\", "/", -1), err1
}

// 初始化 zero log 日志
func NewZeroLog(path, logFileNamePrefix string, stdoutFlag bool, skipCall int) zerolog.Logger {
	var logfilename string
	dataTimeStr := time.Now().Format(LogFileNameTimeFormat)
	if len(path) == 0 {
		path, _ = GetCurrentExecDir()
	}
	logpath := path + "/logs/" // + dataTimeStr

	afs := afero.NewOsFs()
	check, _ := afero.DirExists(afs, logpath)
	if !check {
		err := os.MkdirAll(logpath, 0755)
		fmt.Printf("Logger create dir error  %d ", err)
	}

	if len(logFileNamePrefix) == 0 {
		logfilename = logpath + "/pid-" + strconv.Itoa(os.Getpid()) + "-" + dataTimeStr + ".log"
	} else {
		logfilename = logpath + "/" + logFileNamePrefix + ".log"
	}

	LumberLogger = &Logger{
		Filename:   logfilename,
		MaxSize:    LogMaxSize, // megabytes
		MaxBackups: LogMaxBackups,
		MaxAge:     LogMaxAge, // days
		Compress:   false,     // 开发时不压缩
	}

	wdiode := diode.NewWriter(LumberLogger, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})

	var writers []io.Writer

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if stdoutFlag {
		writers = []io.Writer{
			wdiode,
			os.Stdout,
		}
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		writers = []io.Writer{
			wdiode,
		}
	}

	multi := io.MultiWriter(writers...)
	// 	zerolog.TimeFieldFormat = time.RFC3339Nano //  "2006-01-02/15:04:05.999999999" //15:04:05.999999999
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.999999999" // 15:04:05.999999999
	//zerolog.TimestampFieldName = "t"
	//zerolog.LevelFieldName = "l"
	//zerolog.MessageFieldName = "m"
	//zerolog.CallerFieldName = "c"
	zerolog.DurationFieldUnit = time.Nanosecond
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if skipCall == -1 {
		perfectlog := zerolog.New(multi).With().Timestamp().Stack().Caller().Logger()
		return perfectlog
	} else {
		perfectlog := zerolog.New(multi).With().Timestamp().Stack().CallerWithSkipFrameCount(skipCall).Logger()
		return perfectlog
	}
}

// Error Log ERROR level message.
func Error(format string, v ...interface{}) {
	if v == nil {
		Log.Error().Msg(format)
	} else {
		Log.Error().Msgf(format, v...)
	}
}

// Debug Log DEBUG level message.
func Debug(format string, v ...interface{}) {
	if v == nil {
		Log.Debug().Msg(format)
	} else {
		Log.Debug().Msgf(format, v...)
	}
}

// Warn Log WARN level message.
// compatibility alias for Warning()
func Fatal(format string, v ...interface{}) {
	if v == nil {
		Log.Fatal().Msg(format)
	} else {
		Log.Fatal().Msgf(format, v...)
	}
}

// Info Log INFO level message.
func Info(format string, v ...interface{}) {
	if v == nil {
		Log.Info().Msg(format)
	} else {
		Log.Info().Msgf(format, v...)
	}
}
