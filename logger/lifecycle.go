package logger

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/rs/zerolog"
)

// Lifecycle - запись событий запуска-остановки приложения.
// Пишет ТОЛЬКО в lifecycle файл, игнорируя настройки уровня основного логгера
func (l *Logger) Lifecycle(message string, args ...interface{}) {
	if l.lifecycleLog != nil {
		// Используем Info(), так как это просто запись. Уровень скрыт в FormatLevel
		l.lifecycleLog.Info().Msgf(message, args...)
	}
}

func (l *Logger) GetLifecycleFilename() string {
	if l.lifecycleFile != nil {
		return l.lifecycleFile.Name()
	}
	return ""
}

// getLifecycleFilename генерирует имя файла вида name-lifecycle.ext из переданного filename
func getLifecycleFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	return base + "-lifecycle" + ext
}

// SetLifecycleLogFile - устанавливает файл для логгера lifecycle и возвращает сам логгер и файл
func SetLifecycleLogFile(filename string) (*log.Logger, *os.File, error) {
	var lifecycleLogger *log.Logger
	var lifecycleFile *os.File
	if filename == "" {
		return lifecycleLogger, lifecycleFile, nil
	}

	lfName := getLifecycleFilename(filename)
	// Открываем файл для lifecycle (создаем или добавляем)
	lFile, err := os.OpenFile(lfName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o666)
	if err == nil {
		lifecycleFile = lFile
		// Используем ConsoleWriter для файла, чтобы получить формат "Дата Текст" без JSON
		lWriter := log.ConsoleWriter{
			Out:        lFile,
			TimeFormat: HumanTimeFormat,
			NoColor:    true, // В файле цвета не нужны
			// Убираем уровень логгера (INF, DBG) из вывода lifecycle, оставляем только время и сообщение
			FormatLevel: func(i interface{}) string { return "" },
		}
		// Создаем отдельный инстанс zerolog
		logger := log.New(lWriter).With().Timestamp().Logger()
		lifecycleLogger = &logger
	} else {
		return nil, nil, err
	}
	return lifecycleLogger, lifecycleFile, nil
}
