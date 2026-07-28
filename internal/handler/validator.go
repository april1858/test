package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomDateValidator реализует интерфейс validator.Func.
// Она проверяет, соответствует ли строка формату "MM-YYYY".
var CustomDateValidator validator.Func = func(fl validator.FieldLevel) bool {
	dateStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	// Если поле пустое, но помечено omitempty, валидатор go-playground 	пропустит его до вызова этой функции.
	// Если пустая строка попала сюда, значит валидация провалена.
	if dateStr == "" {
		return false
	}
	// Проверяем строгое соответствие шаблону "ММ-ГГГГ"
	_, err := time.Parse("01-2006", dateStr)
	return err == nil
}
