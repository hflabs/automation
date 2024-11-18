package logger

import "time"

type RotateConfig struct {
	DatePattern   string         // Шаблон даты-времени, который будет добавляться после названия файла лога. В формате `%Y%m%d%H%M`
	RotationTime  time.Duration  // С какой периодичностью нужно производить ротацию файлов
	RotationCount int            // Количество резервных файлов
	MaxSize       int            // Максимальный размер файла в МБ
	TimeLocation  *time.Location // Часовой пояс, по которому происходит ротация
}
