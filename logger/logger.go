package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/rs/zerolog"
)

const HumanTimeFormat = "02.01.2006 15:04:05"

// Interface -.
type Interface interface {
	Trace(message interface{}, args ...interface{})
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
	Lifecycle(message string, args ...interface{})
	SetLogLevel(level string)
	GetLogLevel() string
	GetLifecycleFilename() string
	EnableLifecycle(baseFilename string) error
	Close() error
}

// Logger - структура логера
type Logger struct {
	logger        *log.Logger
	lifecycleLog  *log.Logger    // Отдельный логгер для lifecycle событий
	logFile       io.WriteCloser // Используем интерфейс, чтобы закрывать не только файл, но и ротатор тоже
	logFilename   string
	lifecycleFile *os.File // Файл для lifecycle
	formatter     log.Formatter
	rotating      *RotateConfig
	isSimple      bool
}

var _ Interface = (*Logger)(nil)

// New - Создает новый экземпляр логгера с выводом в консоль и файл с заданным форматированием и настройками ро
func New(level, filename string, formatter log.Formatter, rotating *RotateConfig) Interface {
	lev := parseLogLevel(level)
	var output log.LevelWriter
	var closer io.WriteCloser

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o666)
	console := newConsoleWriter()
	if err != nil {
		output = log.MultiLevelWriter(console)
		filename = ""
	} else {
		if rotating != nil && rotating.DatePattern != "" {
			rl, rlErr := rotatelogs.New(
				fmt.Sprintf("%s.%s", filename, rotating.DatePattern),
				rotatelogs.WithLinkName(filename),
				rotatelogs.WithRotationTime(rotating.RotationTime),
				rotatelogs.WithRotationCount(rotating.RotationCount),
				rotatelogs.WithLocation(rotating.TimeLocation),
			)
			if rlErr != nil {
				output = log.MultiLevelWriter(newFormattedWriter(file, formatter), console)
				closer = file
			} else {
				// Ротация работает, поэтому больше не нужен дескриптор 'file', так как rl сам открывает и закрывает файлы.
				file.Close()
				output = log.MultiLevelWriter(newFormattedWriter(rl, formatter), console)
				closer = rl // Теперь Close() закроет именно ротатор
			}
		} else {
			output = log.MultiLevelWriter(newFormattedWriter(file, formatter), console)
			closer = file
		}
	}

	logger := log.New(output).Level(lev).With().Timestamp().Logger()

	msg := fmt.Sprintf("Configured logger with level - %s, filename - %s", lev.String(), filename)
	if rotating != nil {
		msg += fmt.Sprintf(", rotating - (rotationTime = %s, rotationCount = %d, datePattern = `%s`, location = %s)",
			rotating.RotationTime.String(), rotating.RotationCount, rotating.GetDatePattern(), rotating.TimeLocation.String())
	}
	if formatter != nil {
		msg += ", with formatter"
	}
	logger.Trace().Msg(msg)

	return &Logger{
		logger:      &logger,
		logFile:     closer,
		logFilename: filename,
		formatter:   formatter,
		rotating:    rotating,
		isSimple:    false,
	}
}

// NewSimple - Создает новый экземпляр логгера с выводом только в консоль
func NewSimple(level string) *Logger {
	lev := parseLogLevel(level)
	var output log.LevelWriter
	writer := newConsoleWriter()
	output = log.MultiLevelWriter(writer)
	logger := log.New(output).Level(lev).With().Timestamp().Logger()
	return &Logger{
		logger:   &logger,
		isSimple: true,
	}
}

func (l *Logger) EnableLifecycle(baseFilename string) error {
	lifecycleLogger, lifecycleFile, err := SetLifecycleLogFile(baseFilename)
	if err != nil {
		l.Error("Failed to create lifecycle log: %v", err)
		return err
	}
	l.lifecycleLog = lifecycleLogger
	l.lifecycleFile = lifecycleFile
	return nil
}

func (l *Logger) SetLogLevel(level string) {
	lev := parseLogLevel(level)

	var output log.LevelWriter
	writerConsole := newConsoleWriter()
	if l.isSimple {
		output = log.MultiLevelWriter(writerConsole)
	} else {
		// Для обычного логгера переиспользуем существующий l.logFile, таким образом избегаем накопления открытых дескрипторов в системе
		output = log.MultiLevelWriter(newFormattedWriter(l.logFile, l.formatter), writerConsole)
	}
	newLogger := log.New(output).Level(lev).With().Timestamp().Logger()
	l.logger = &newLogger

	l.logger.Trace().Msgf("Log level changed to %s on the fly", lev.String())
}

func (l *Logger) GetLogLevel() string {
	return l.logger.GetLevel().String()
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
	var errCommon, errLifecycle error
	if l.lifecycleFile != nil {
		// Проверяем, не является ли файл стандартным потоком перед закрытием
		if l.lifecycleFile != os.Stdout && l.lifecycleFile != os.Stderr {
			errLifecycle = l.lifecycleFile.Close()
		}
		l.lifecycleFile = nil // Защита от повторного закрытия
	}
	if l.logFile != nil {
		// Закрываем либо файл, либо ротатор
		errCommon = l.logFile.Close()
		l.logFile = nil
	}
	return errors.Join(errLifecycle, errCommon)
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

func newConsoleWriter() log.ConsoleWriter {
	return log.ConsoleWriter{Out: os.Stdout, TimeFormat: HumanTimeFormat, NoColor: false}
}

func newFormattedWriter(out io.WriteCloser, formatter log.Formatter) log.ConsoleWriter {
	return log.ConsoleWriter{
		Out:           out,
		TimeFormat:    HumanTimeFormat,
		NoColor:       true,
		FormatMessage: formatter,
	}
}
