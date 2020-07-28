package cron

import (
	"fmt"
	"github.com/kniepok/corgi"
	"strings"
)

func parseDetailsToCron(details corgi.SubscriptionDetails) (string, error) {
	switch details.Mode {
	case corgi.SubModeDaily:
		return fmt.Sprintf("%d %d * * *", details.Minute, details.Hour), nil
	case corgi.SubModeWeekly:
		day := strings.ToUpper(details.DayOfWeek.String()[0:2])
		return fmt.Sprintf("%d %d * * %s", details.Minute, details.Hour, day), nil
	}
	return "", fmt.Errorf("cannot parse subscription details to cron")
}
