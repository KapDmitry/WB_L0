package validator

import (
	"github.com/KapDmitry/WB_L0/internal/logger"
	"github.com/KapDmitry/WB_L0/internal/order"
	"github.com/go-playground/validator"
)

func Validate(ord order.Order, log logger.Logger) bool {
	v := validator.New()
	err := v.Struct(ord)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			//fmt.Printf("Ошибка валидации для поля %s: %s\n", e.Field(), e.Tag())
			log.LogW("Info", "validation err for: ", map[string]interface{}{e.Field(): e.Tag()})
		}
		return false
	}
	return true
}
