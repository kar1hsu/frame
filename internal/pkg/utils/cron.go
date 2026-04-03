package utils

import (
	"fmt"
)

// cron 语法封装

func EverySeconds(seconds int) string {
	return fmt.Sprintf("@every %ds", seconds)
}

func EveryMinutes(minutes int) string {
	return fmt.Sprintf("@every %dm", minutes)
}

func EveryHourly() string {
	return "@hourly"
}

func EveryDaily() string {
	return "@daily"
}

func DailyAt(hour, minute int) string {
	return fmt.Sprintf("%d %d * * *", minute, hour)
}

func WeeklyAt(weekday, hour, minute int) string {
	return fmt.Sprintf("%d %d * * %d", minute, hour, weekday)
}

func EveryWeek() string {
	return "@weekly"
}

func MonthlyAt(day, hour, minute int) string {
	return fmt.Sprintf("%d %d %d * *", minute, hour, day)
}

func EveryMonth() string {
	return "@monthly"
}

func EveryYearly() string {
	return "@yearly"
}
