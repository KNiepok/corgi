package corgi

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ScheduledSubscription struct {
	Subscription
	EntryID int
}

// Subscription describes to which user a notification will be sent at a given interval
type Subscription struct {
	User    User
	Details SubscriptionDetails
}

type User struct {
	ID    string
	Email string
}

type SubscriptionMode int

const (
	SubModeDaily SubscriptionMode = iota
	SubModeWeekly
)

type SubscriptionDetails struct {
	Mode         SubscriptionMode
	DayOfWeek    time.Weekday
	Hour, Minute uint8
}

type SubscriptionService interface {
	//	Subscribe adds user to notify list
	Subscribe(ctx context.Context, request Subscription) error
	// Unsubscribe removes user subscription
	Unsubscribe(ctx context.Context, user User) error
}

type StarterService interface {
	Start(ctx context.Context) error
}

// Scheduler can add or remove new subscriptions
type Scheduler interface {
	// Add registers new notification schedule
	Add(ctx context.Context, interval SubscriptionDetails, task func()) (entryID int, err error)
	// Remove deletes already existing notification schedule
	Remove(ctx context.Context, entryID int) error
}

// SubscriptionStorage can return a list of already created subscriptions
type SubscriptionStorage interface {
	// Create or update subscription entry
	Add(ctx context.Context, sub ScheduledSubscription) error
	// Get scheduled subscription by user ID
	Get(ctx context.Context, userID string) (ScheduledSubscription, error)
	// Remove scheduled subscription based on user ID
	Remove(ctx context.Context, userID string) error
	// Fetch all scheduled subscriptions
	GetAll(ctx context.Context) ([]ScheduledSubscription, error)
}

type NotificationService interface {
	Notify(ctx context.Context, sub Subscription) error
}

// NotificationSender can send notifications to given users
type NotificationSender interface {
	Send(ctx context.Context, user User, message string) error
}

// NotificationGenerator can create notification messages
type NotificationGenerator interface {
	CreateNotification(ctx context.Context, sub Subscription) (msg string, err error)
}

type UserResolver interface {
	ResolveUserEmail(ctx context.Context, userID string) (email string, err error)
}

var ErrInvalidInterval = fmt.Errorf("invalid interval")
var ErrSubscriptionNotFound = fmt.Errorf("cannot find subscription")

func NewSubscriptionDetails(input string) (SubscriptionDetails, error) {
	elements := strings.Split(input, " ")
	if len(elements) == 0 {
		return SubscriptionDetails{}, fmt.Errorf("cannot parse requested notification interval to anything meaningful")
	}
	switch elements[0] {
	case "daily":
		return parseDailyIntervalToCron(input)
	case "weekly":
		return parseWeeklyIntervalToCron(input)
	}
	return SubscriptionDetails{}, fmt.Errorf("cannot parse %s", input)
}

// parseDailyIntervalToCron parses daily request to valid cron.
// Legit daily interval requests look like this:
// daily @ 17:20
// daily @ 9
// Invalid daily requests are:
// daily @ 25:12
// daily @ 12:89
// daily1730
func parseDailyIntervalToCron(input string) (SubscriptionDetails, error) {
	elements := strings.Split(input, " ")
	hour, minute, err := parseTime(elements[len(elements)-1])
	if err != nil {
		return SubscriptionDetails{}, err
	}

	return SubscriptionDetails{
		Mode:   SubModeDaily,
		Hour:   hour,
		Minute: minute,
	}, nil
}

// parseWeeklyIntervalToCron parses weekly request to valid cron.
// Legit weekly interval requests look like this:
// weekly @ FRI 17:20
// weekly @ SAT 9
// Invalid weekly requests are:
// weekly @ 25:12
// weekly @ ABC 17:00
// weekly @ FRI 17:xx
func parseWeeklyIntervalToCron(interval string) (SubscriptionDetails, error) {
	elements := strings.Split(interval, " ")
	if len(elements) < 2 {
		return SubscriptionDetails{}, fmt.Errorf("cannot parse %s as valid weekly interval; too few elements", interval)
	}

	hour, minute, err := parseTime(elements[len(elements)-1])
	if err != nil {
		return SubscriptionDetails{}, err
	}

	dayInput := strings.ToUpper(elements[len(elements)-2])
	weekday, exists := daysOfWeek[dayInput]
	if !exists {
		return SubscriptionDetails{}, fmt.Errorf("cannot parse %s into a valid day", dayInput)
	}

	return SubscriptionDetails{
		Mode:      SubModeWeekly,
		DayOfWeek: weekday,
		Hour:      hour,
		Minute:    minute,
	}, nil
}

// parseTime will take 17:20 and return 17,20,nil
func parseTime(timeString string) (uint8, uint8, error) {
	t := strings.Split(timeString, ":")
	hour, err := parseHour(t[0])
	if err != nil {
		return 0, 0, err
	}

	if len(t) == 1 {
		return hour, 0, nil
	}

	minute, err := parseMinute(t[1])
	if err != nil {
		return 0, 0, err
	}
	return hour, minute, nil
}

var daysOfWeek = map[string]time.Weekday{
	"SUN": time.Sunday,
	"MON": time.Monday,
	"TUE": time.Tuesday,
	"WED": time.Wednesday,
	"THU": time.Thursday,
	"FRI": time.Friday,
	"SAT": time.Saturday,
}

func parseHour(hour string) (uint8, error) {
	return parseIntBetween(hour, 0, 23)
}

func parseMinute(minute string) (uint8, error) {
	return parseIntBetween(minute, 0, 59)
}

func parseIntBetween(s string, lower, upper uint8) (uint8, error) {
	number, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %s as integer", s)
	}
	n := uint8(number)
	if n < lower || n > upper {
		return 0, fmt.Errorf("number has to be between %d and %d, got :%d", lower, upper, number)
	}
	return n, nil
}
