package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_getLifecycleFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "1. Стандартный путь с расширением",
			filename: "log/updater.log",
			want:     "log/updater-lifecycle.log",
		},
		{
			name:     "2. Файл в корне без папок",
			filename: "app.log",
			want:     "app-lifecycle.log",
		},
		{
			name:     "3. Путь без расширения файла",
			filename: "/var/log/app",
			want:     "/var/log/app-lifecycle",
		},
		{
			name:     "4. Имя файла с множеством точек",
			filename: "my.super.app.v1.log",
			want:     "my.super.app.v1-lifecycle.log",
		},
		{
			name:     "5. Пустая строка",
			filename: "",
			want:     "-lifecycle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLifecycleFilename(tt.filename); got != tt.want {
				t.Errorf("getLifecycleFilename() = %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

func TestSetLifecycleLogFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		filename string // Функция-генератор пути, чтобы использовать tempDir
		wantErr  bool
		wantNil  bool
	}{
		{
			name:     "1. Пустое имя файла (логгер не создается)",
			filename: "",
			wantErr:  false, wantNil: true,
		},
		{
			name:     "2. Успешное создание файла и логгера",
			filename: filepath.Join(tempDir, "success.log"),
			wantErr:  false, wantNil: false,
		},
		{
			name:     "3. Ошибка создания (несуществующая директория)",
			filename: filepath.Join(tempDir, "no_dir", "fail.log"), // Пытаемся создать файл в папке, которой нет
			wantErr:  true, wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, file, err := SetLifecycleLogFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLifecycleLogFile() ошибка = %v, ожидалась ошибка %v", err, tt.wantErr)
				return
			}
			if (logger == nil) != tt.wantNil {
				t.Errorf("SetLifecycleLogFile() logger is nil = %v, ожидалось nil = %v", logger == nil, tt.wantNil)
			}

			// Дополнительная проверка: если логгер создан, проверяем запись в файл
			if logger != nil && file != nil {
				defer file.Close()

				// 1. Проверяем, что файл физически существует
				expectedPath := getLifecycleFilename(tt.filename)
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Fatalf("Файл лога %s не был создан", expectedPath)
				}
				// 2. Пробуем записать что-то в логгер
				testMsg := "тестовое сообщение жизненного цикла"
				logger.Info().Msg(testMsg)
				// 3. Читаем файл и проверяем содержимое
				// Небольшая задержка не нужна, так как запись в файл идет синхронно через ConsoleWriter, но для надежности чтения:
				contentBytes, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatalf("Не удалось прочитать созданный файл: %v", err)
				}
				content := string(contentBytes)
				// 4. Проверяем наличие сообщения
				if !strings.Contains(content, testMsg) {
					t.Errorf("Файл не содержит ожидаемое сообщение. Содержимое: %q", content)
				}
				// 5. Проверяем отсутствие уровня логирования (INFO) согласно требованию
				if strings.Contains(content, "INF") || strings.Contains(content, "INFO") {
					t.Errorf("Файл содержит уровень логирования, хотя он должен быть скрыт. Содержимое: %q", content)
				}
				// 6. Проверяем наличие текущей даты
				currentDate := time.Now().Format("02.01.2006")
				if !strings.Contains(content, currentDate) {
					t.Errorf("Файл не содержит текущую дату (%s). Содержимое: %q", currentDate, content)
				}
			}
		})
	}
}
