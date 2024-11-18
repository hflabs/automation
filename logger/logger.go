package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/rs/zerolog"
	"io"
	"os"
	"strings"
)

// Interface -.
type Interface interface {
	Trace(message interface{}, args ...interface{})
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
	SetLogLevel(level string)
	GetLogLevel() string
	Close() error
}

// Logger - структура логера
type Logger struct {
	logger    *log.Logger
	logFile   *os.File
	formatter log.Formatter
	rotating  *RotateConfig
	isSimple  bool
}

func (r *RotateConfig) GetDatePattern() string {
	return strings.ReplaceAll(r.DatePattern, "%", "%%")
}

var _ Interface = (*Logger)(nil)

// New - Создает новый экземпляр логгера с выводом в консоль и файл с заданным форматированием и настройками ро
func New(level, filename string, formatter log.Formatter, rotating *RotateConfig) *Logger {
	lev := parseLogLevel(level)
	var output log.LevelWriter
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o666)
	writer := log.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02.01.2006 15:04:05",
		NoColor:    false,
	}
	if err != nil {
		output = log.MultiLevelWriter(writer)
		filename = ""
	} else {
		var out io.Writer
		if rotating.DatePattern != "" {
			out, err = rotatelogs.New(
				fmt.Sprintf("%s.%s", filename, rotating.DatePattern),
				rotatelogs.WithLinkName(filename),
				rotatelogs.WithRotationTime(rotating.RotationTime),
				rotatelogs.WithRotationCount(rotating.RotationCount),
				rotatelogs.WithClock(rotatelogs.Local),
			)
			if err != nil {
				out = file
				rotating = nil
			}
		} else {
			out = file
		}
		writerFile := log.ConsoleWriter{
			Out:           out,
			TimeFormat:    "02.01.2006 15:04:05",
			NoColor:       true,
			FormatMessage: formatter,
		}
		output = log.MultiLevelWriter(writerFile, writer)
	}

	logger := log.New(output).Level(lev).With().Timestamp().Logger()

	msg := fmt.Sprintf("Configated logger with level - %s, filename - %s", lev.String(), filename)
	if rotating != nil {
		msg += fmt.Sprintf(", rotating - (rotationTime = %s, rotationCount = %d, datePattern = `%s`, location = %s)",
			rotating.RotationTime.String(), rotating.RotationCount, rotating.GetDatePattern(), rotating.TimeLocation.String())
	}
	if formatter != nil {
		msg += ", with formatter"
	}
	logger.Debug().Msgf(msg)
	return &Logger{
		logger:    &logger,
		logFile:   file,
		formatter: formatter,
		rotating:  rotating,
		isSimple:  false,
	}
}

// NewSimple - Создает новый экземпляр логгера с выводом только в консоль
func NewSimple(level string) *Logger {
	lev := parseLogLevel(level)
	var output log.LevelWriter
	writer := log.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02.01.2006 15:04:05",
		NoColor:    false,
	}
	output = log.MultiLevelWriter(writer)
	logger := log.New(output).Level(lev).With().Timestamp().Logger()
	return &Logger{
		logger:   &logger,
		isSimple: true,
	}
}

func (l *Logger) SetLogLevel(level string) {
	var newLogger *Logger
	if l.isSimple {
		newLogger = NewSimple(level)
	} else {
		newLogger = New(level, l.logFile.Name(), l.formatter, l.rotating)
	}
	l.logger = newLogger.logger
}

func (l *Logger) GetLogLevel() string {
	return l.logger.GetLevel().String()
}
func (l *Logger) Rotate() {

}

// Trace - Обработка результата типа TRACE
func (l *Logger) Trace(message interface{}, args ...interface{}) {
	l.msg("trace", message, args...)
}

// Debug - Обработка результата типа DEBUG
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

// Info - обработка результата типа INFO
func (l *Logger) Info(message string, args ...interface{}) {
	l.msg("info", message, args...)
}

// Warn - обработка результата типа WARN
func (l *Logger) Warn(message string, args ...interface{}) {
	l.msg("warn", message, args...)
}

// Error - обработка результата типа ERROR
func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.msg("error", message, args...)
}

// Fatal - обработка результата типа FATAL
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Logger) log(level string, message string, args ...interface{}) {
	if len(args) == 0 {
		switch level {
		case "fatal":
			l.logger.Fatal().Msg(message)
		case "error":
			l.logger.Error().Msg(message)
		case "warn":
			l.logger.Warn().Msg(message)
		case "info":
			l.logger.Info().Msg(message)
		case "debug":
			l.logger.Debug().Msg(message)
		case "trace":
			l.logger.Trace().Msg(message)
		}
	} else {
		switch level {
		case "fatal":
			l.logger.Fatal().Msgf(message, args...)
		case "error":
			l.logger.Error().Msgf(message, args...)
		case "warn":
			l.logger.Warn().Msgf(message, args...)
		case "info":
			l.logger.Info().Msgf(message, args...)
		case "debug":
			l.logger.Debug().Msgf(message, args...)
		case "trace":
			l.logger.Trace().Msgf(message, args...)
		}
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(level, msg.Error(), args...)
	case string:
		l.log(level, msg, args...)
	default:
		l.log(level, fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}

func parseLogLevel(level string) log.Level {
	var lev log.Level
	switch strings.ToLower(level) {
	case "info":
		lev = log.InfoLevel
	case "warn":
		lev = log.WarnLevel
	case "error":
		lev = log.ErrorLevel
	case "debug":
		lev = log.DebugLevel
	case "trace":
		lev = log.TraceLevel
	case "fatal":
		lev = log.FatalLevel
	case "disabled":
		lev = log.Disabled
	default:
		lev = log.InfoLevel
		fmt.Println("некорректное значения для переменной logging_level, ожидалось одно из списка: " +
			"error, warn, info, debug или trace. Используем умолчательное значение info")
	}
	return lev
}
