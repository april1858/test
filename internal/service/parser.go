package service

import "time"

// layoutMMYYYY определяет формат даты согласно ТЗ.
const layoutMMYYYY = "01-2006"

// parseMMYYYY конвертирует строку вида "07-2025" в time.Time (всегда первое число указанного месяца).

func parseMMYYYY(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(layoutMMYYYY, dateStr)
}

// formatMMYYYY конвертирует time.Time обратно в строку "MM-YYYY" (для тестов или кастомного вывода).
func formatMMYYYY(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(layoutMMYYYY)
}
