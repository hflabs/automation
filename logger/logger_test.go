package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger_SetLogLevel(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test_level.log")

	// Создаем конфиг ротации (пустой, чтобы просто писать в файл)
	rot := &RotateConfig{}
	tests := []struct {
		name         string
		initialLevel string
		newLevel     string
		checkMsg     string
		shouldExist  bool
	}{
		{
			name:         "1. Повышение уровня: DEBUG не пишется в INFO",
			initialLevel: "info",
			newLevel:     "info",
			checkMsg:     "this debug should not be seen",
			shouldExist:  false,
		},
		{
			name:         "2. Смена уровня на лету: DEBUG начинает писаться",
			initialLevel: "info",
			newLevel:     "debug",
			checkMsg:     "this debug MUST be seen",
			shouldExist:  true,
		},
		{
			name:         "3. Понижение уровня обратно до WARN",
			initialLevel: "debug",
			newLevel:     "warn",
			checkMsg:     "info message after downgrade",
			shouldExist:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем файл перед каждым тестом
			_ = os.Truncate(logPath, 0)
			l := New(tt.initialLevel, logPath, nil, rot)
			defer l.Close()
			// Меняем уровень
			l.SetLogLevel(tt.newLevel)
			// Пытаемся записать сообщение
			switch tt.newLevel {
			case "debug", "info": // Не должно записаться при Info
				l.Debug(tt.checkMsg)
			case "warn":
				l.Info(tt.checkMsg) // Не должно записаться при Warn
			}

			// Читаем файл
			content, _ := os.ReadFile(logPath)
			got := strings.Contains(string(content), tt.checkMsg)
			if got != tt.shouldExist {
				t.Errorf("%s: ожидалось наличие сообщения = %v, получено %v. Содержимое: %q",
					tt.name, tt.shouldExist, got, string(content))
			}
		})
	}
}

func TestLogger_Close(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test_close.log")

	tests := []struct {
		name      string
		isSimple  bool
		lifecycle bool
	}{
		{name: "1. Закрытие простого логгера (stdout)", isSimple: true, lifecycle: false},
		{name: "2. Закрытие файлового логгера", isSimple: false, lifecycle: false},
		{name: "3. Закрытие логгера с активным Lifecycle", isSimple: false, lifecycle: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l Interface
			if tt.isSimple {
				l = NewSimple("info")
			} else {
				l = New("info", logPath, nil, &RotateConfig{})
			}

			if tt.lifecycle {
				err := l.EnableLifecycle(logPath)
				if err != nil {
					t.Fatalf("не удалось включить lifecycle: %v", err)
				}
			}
			// Проверяем закрытие
			err := l.Close()
			if err != nil {
				t.Errorf("%s: ошибка при закрытии: %v", tt.name, err)
			}
			// Проверка на "панику" при повторном закрытии
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s: метод Close() вызвал панику при повторном вызове: %v", tt.name, r)
				}
			}()
			err = l.Close()
			if err != nil {
				// Повторное закрытие может вернуть ошибку (зависит от ОС), но не должно вешать систему
				t.Logf("повторное закрытие вернуло ожидаемую ошибку/nil: %v", err)
			}
		})
	}
}

func TestLogger_NoFileLeaksOnSetLogLevel(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "leak_test.log")

	l := New("info", logPath, nil, &RotateConfig{})
	defer l.Close()

	// Запоминаем текущий дескриптор файла чтобы проверить, что количество открытых файлов процессом не растет
	// Для простоты проверим логику: SetLogLevel не должен менять поле logFile на новый объект
	initialFilePtr := l.(*Logger).logFile
	for i := 0; i < 10; i++ {
		l.SetLogLevel("debug")
		l.SetLogLevel("info")
	}
	if l.(*Logger).logFile != initialFilePtr {
		t.Errorf("SetLogLevel изменил дескриптор файла! Это приведет к утечке, если старый файл не был закрыт.")
	}
}
