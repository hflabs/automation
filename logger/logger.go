package logger

import (
	"fmt"
	log "github.com/rs/zerolog"
	"os"
	"strings"
)

// Interface -.
type Interface interface {
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
	isSimple  bool
}

var _ Interface = (*Logger)(nil)

// New - Создает новый экземпляр логгера с выводом в консоль и файл с заданным форматированием
func New(level, filename string, formatter log.Formatter) *Logger {
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
	} else {
		writerFile := log.ConsoleWriter{
			Out:        file,
			TimeFormat: "02.01.2006 15:04:05",
			NoColor:    true,
		}
		if formatter != nil {
			writerFile.FormatMessage = formatter
		}
		output = log.MultiLevelWriter(writerFile, writer)
	}

	logger := log.New(output).Level(lev).With().Timestamp().Logger()
	return &Logger{
		logger:    &logger,
		logFile:   file,
		formatter: formatter,
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
		newLogger = New(level, l.logFile.Name(), l.formatter)
	}
	l.logger = newLogger.logger
}

func (l *Logger) GetLogLevel() string {
	return l.logger.GetLevel().String()
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
	case "fatal":
		lev = log.FatalLevel
	case "disabled":
		lev = log.Disabled
	default:
		lev = log.InfoLevel
		fmt.Println("некорректное значения для переменной logging_level, ожидалось одно из списка: " +
			"error, warn, info или debug. Используем умолчательное значение info")
	}
	return lev
}
