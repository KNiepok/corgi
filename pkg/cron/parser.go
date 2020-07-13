package cron

import (
	"fmt"
	"strconv"
	"strings"
)

func parseIntervalToCron(interval string) (string, error) {
	elements := strings.Split(interval, " ")
	if len(elements) == 0 {
		return "", fmt.Errorf("cannot parse requested notification interval to anything meaningful")
	}
	switch elements[0] {
	case "daily":
		return parseDailyIntervalToCron(interval)
	case "weekly":
		return parseWeeklyIntervalToCron(interval)
	}
	return "", fmt.Errorf("cannot parse %s", interval)
}

// parseDailyIntervalToCron parses daily request to valid cron.
// Legit daily interval requests look like this:
// daily @ 17:20
// daily @ 9
// Invalid daily requests are:
// daily @ 25:12
// daily @ 12:89
// daily1730
func parseDailyIntervalToCron(interval string) (string, error) {
	elements := strings.Split(interval, " ")
	hour, minute, err := parseTime(elements[len(elements)-1])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d %d * * *", minute, hour), nil
}

// parseTime will take 17:20 and return 17,20,nil
func parseTime(timeString string) (int, int, error) {
	time := strings.Split(timeString, ":")
	hour, err := parseHour(time[0])
	if err != nil {
		return 0, 0, err
	}

	if len(time) == 1 {
		return hour, 0, nil
	}

	minute, err := parseMinute(time[1])
	if err != nil {
		return 0, 0, err
	}
	return hour, minute, nil
}

// parseWeeklyIntervalToCron parses weekly request to valid cron.
// Legit weekly interval requests look like this:
// weekly @ FRI 17:20
// weekly @ SAT 9
// Invalid weekly requests are:
// weekly @ 25:12
// weekly @ ABC 17:00
// weekly @ FRI 17:xx
func parseWeeklyIntervalToCron(interval string) (string, error) {
	elements := strings.Split(interval, " ")
	if len(elements) < 2 {
		return "", fmt.Errorf("cannot parse %s as valid weekly interval; too few elements", interval)
	}

	hour, minute, err := parseTime(elements[len(elements)-1])
	if err != nil {
		return "", err
	}

	day := strings.ToUpper(elements[len(elements)-2])
	validDays := []string{"MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}
	isDayValid := false
	for _, d := range validDays {
		if d == day {
			isDayValid = true
		}
	}
	if !isDayValid {
		return "", fmt.Errorf("cannot parse %s into a valid day", day)
	}

	return fmt.Sprintf("%d %d * * %s", minute, hour, day), nil
}

func parseHour(hour string) (int, error) {
	return parseIntBetween(hour, 0, 23)
}

func parseMinute(minute string) (int, error) {
	return parseIntBetween(minute, 0, 59)
}

func parseIntBetween(s string, lower, upper int) (int, error) {
	number, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %s as integer", s)
	}
	if number < lower || number > upper {
		return 0, fmt.Errorf("number has to be between %d and %d, got :%d", lower, upper, number)
	}
	return number, nil
}
