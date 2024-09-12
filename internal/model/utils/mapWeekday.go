package utils

import (
	"errors"
	"strings"

	"github.com/minacio00/easyCourt/internal/model"
)

func MapWeekDay(day string) (model.Weekday, error) {
	day = strings.ToLower((strings.TrimSpace(day)))
	day = strings.ReplaceAll(day, "-", "")
	day = strings.ReplaceAll(day, " ", "")

	switch day {
	case "domingo":
		return model.Domingo, nil
	case "segundafeira":
		return model.SegundaFeira, nil
	case "tercafeira":
		return model.TercaFeira, nil
	case "quartafeira":
		return model.QuartaFeira, nil
	case "quintafeira":
		return model.QuintaFeira, nil
	case "sextafeira":
		return model.SextaFeira, nil
	case "sabado":
		return model.Sabado, nil
	default:
		return "", errors.New("invalid week_day")
	}
}
