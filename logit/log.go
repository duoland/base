package logit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the log printer
var Logger *zap.SugaredLogger

// exported functions
var Debugw func(msg string, keysAndValues ...interface{})
var Debugf func(template string, args ...interface{})
var Debug func(args ...interface{})

var Infow func(msg string, keysAndValues ...interface{})
var Infof func(template string, args ...interface{})
var Info func(args ...interface{})

var Warnw func(msg string, keysAndValues ...interface{})
var Warnf func(template string, args ...interface{})
var Warn func(args ...interface{})

var Errorw func(msg string, keysAndValues ...interface{})
var Errorf func(template string, args ...interface{})
var Error func(args ...interface{})

var Fatalw func(msg string, keysAndValues ...interface{})
var Fatalf func(template string, args ...interface{})
var Fatal func(args ...interface{})

var Panicw func(msg string, keysAndValues ...interface{})
var Panicf func(template string, args ...interface{})
var Panic func(args ...interface{})

var LogLevels = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

type Config struct {
	PrintJSON     bool
	PrintStdout   bool
	PrintCaller   bool
	LogLevel      string
	LogOutputDir  string
	BuiltinFields map[string]string
}

// InitLogs init the logging system.
func InitLogs(cfg *Config) (err error) {
	zapLogLevel := zapcore.InfoLevel
	if v, ok := LogLevels[cfg.LogLevel]; ok {
		zapLogLevel = v
	}

	// detect output dir
	if cfg.LogOutputDir == "" {
		cfg.LogOutputDir = "./logs"
	}

	// detect log output levels
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapLogLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	var infoLogWriter io.Writer
	var errorLogWriter io.Writer
	var logEncoder zapcore.Encoder

	// create log writers
	infoLogWriter, err = getRotateWriter(cfg.LogOutputDir, "info.log")
	if err != nil {
		err = fmt.Errorf("create info log writer error, %s", err.Error())
		return
	}
	errorLogWriter, err = getRotateWriter(cfg.LogOutputDir, "error.log")
	if err != nil {
		err = fmt.Errorf("create error log writer error, %s", err.Error())
		return
	}

	// setup the encoder config and options
	logEncoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "ts",
		CallerKey:  "caller",
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006-01-02T15:04:05.999-0700"))
		},
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	logCoreOptions := make([]zap.Option, 0)

	// add the builtin fields
	builtinFields := make([]zap.Field, 0, len(cfg.BuiltinFields))
	for key, value := range cfg.BuiltinFields {
		builtinFields = append(builtinFields, zap.Field{Key: key, String: value, Type: zapcore.StringType})
	}
	logCoreOptions = append(logCoreOptions, zap.Fields(builtinFields...))

	// add caller
	if cfg.PrintCaller {
		logCoreOptions = append(logCoreOptions, zap.AddCaller())
	}

	// check output format
	if cfg.PrintJSON {
		// format as json
		logEncoder = zapcore.NewJSONEncoder(logEncoderConfig)
	} else {
		// format as text
		logEncoder = zapcore.NewConsoleEncoder(logEncoderConfig)
	}

	// pack up output writers
	logOutputList := []zapcore.Core{
		zapcore.NewCore(logEncoder, zapcore.AddSync(infoLogWriter), infoLevel),
		zapcore.NewCore(logEncoder, zapcore.AddSync(errorLogWriter), errorLevel)}
	// add print stdout is required
	if cfg.PrintStdout {
		logOutputList = append(logOutputList, zapcore.NewCore(logEncoder, zapcore.AddSync(os.Stdout), infoLevel))
	}
	logCore := zapcore.NewTee(logOutputList...)

	// create the logger
	Logger = zap.New(logCore, logCoreOptions...).Sugar()

	// exported functions
	Debugw = Logger.Debugw
	Debugf = Logger.Debugf
	Debug = Logger.Debug

	Infow = Logger.Infow
	Infof = Logger.Infof
	Info = Logger.Info

	Warnw = Logger.Warnw
	Warnf = Logger.Warnf
	Warn = Logger.Warn

	Errorw = Logger.Errorw
	Errorf = Logger.Errorf
	Error = Logger.Error

	Fatalw = Logger.Fatalw
	Fatalf = Logger.Fatalf
	Fatal = Logger.Fatal

	Panicw = Logger.Panicw
	Panicf = Logger.Panicf
	Panic = Logger.Panic

	return
}

// create {logDir}/{appName}/2022-05-14/infoLog.log
func getRotateWriter(logOutputDir, fileName string) (w io.Writer, err error) {
	logFilePath := filepath.Join(logOutputDir, "%Y-%m-%d", fileName)
	w, err = rotatelogs.New(logFilePath,
		rotatelogs.WithMaxAge(time.Hour*24*14),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	return
}
